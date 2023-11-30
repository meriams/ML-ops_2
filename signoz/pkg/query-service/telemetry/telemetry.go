package telemetry

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	ph "github.com/posthog/posthog-go"
	"gopkg.in/segmentio/analytics-go.v3"

	"go.signoz.io/signoz/pkg/query-service/constants"
	"go.signoz.io/signoz/pkg/query-service/interfaces"
	"go.signoz.io/signoz/pkg/query-service/model"
	v3 "go.signoz.io/signoz/pkg/query-service/model/v3"
	"go.signoz.io/signoz/pkg/query-service/version"
)

const (
	TELEMETRY_EVENT_PATH                     = "API Call"
	TELEMETRY_EVENT_USER                     = "User"
	TELEMETRY_EVENT_INPRODUCT_FEEDBACK       = "InProduct Feedback Submitted"
	TELEMETRY_EVENT_NUMBER_OF_SERVICES       = "Number of Services"
	TELEMETRY_EVENT_NUMBER_OF_SERVICES_PH    = "Number of Services V2"
	TELEMETRY_EVENT_HEART_BEAT               = "Heart Beat"
	TELEMETRY_EVENT_ORG_SETTINGS             = "Org Settings"
	DEFAULT_SAMPLING                         = 0.1
	TELEMETRY_LICENSE_CHECK_FAILED           = "License Check Failed"
	TELEMETRY_LICENSE_UPDATED                = "License Updated"
	TELEMETRY_LICENSE_ACT_FAILED             = "License Activation Failed"
	TELEMETRY_EVENT_ENVIRONMENT              = "Environment"
	TELEMETRY_EVENT_LANGUAGE                 = "Language"
	TELEMETRY_EVENT_LOGS_FILTERS             = "Logs Filters"
	TELEMETRY_EVENT_DISTRIBUTED              = "Distributed"
	TELEMETRY_EVENT_QUERY_RANGE_V3           = "Query Range V3 Metadata"
	TELEMETRY_EVENT_ACTIVE_USER              = "Active User"
	TELEMETRY_EVENT_ACTIVE_USER_PH           = "Active User V2"
	TELEMETRY_EVENT_USER_INVITATION_SENT     = "User Invitation Sent"
	TELEMETRY_EVENT_USER_INVITATION_ACCEPTED = "User Invitation Accepted"
	DEFAULT_CLOUD_EMAIL                      = "admin@signoz.cloud"
)

var SAAS_EVENTS_LIST = map[string]struct{}{
	TELEMETRY_EVENT_NUMBER_OF_SERVICES:       {},
	TELEMETRY_EVENT_ACTIVE_USER:              {},
	TELEMETRY_EVENT_HEART_BEAT:               {},
	TELEMETRY_EVENT_LANGUAGE:                 {},
	TELEMETRY_EVENT_ENVIRONMENT:              {},
	TELEMETRY_EVENT_USER_INVITATION_SENT:     {},
	TELEMETRY_EVENT_USER_INVITATION_ACCEPTED: {},
}

const api_key = "4Gmoa4ixJAUHx2BpJxsjwA1bEfnwEeRz"
const ph_api_key = "H-htDCae7CR3RV57gUzmol6IAKtm5IMCvbcm_fwnL-w"

const IP_NOT_FOUND_PLACEHOLDER = "NA"
const DEFAULT_NUMBER_OF_SERVICES = 6

const HEART_BEAT_DURATION = 6 * time.Hour

const ACTIVE_USER_DURATION = 30 * time.Minute

// const HEART_BEAT_DURATION = 30 * time.Second
// const ACTIVE_USER_DURATION = 30 * time.Second

const RATE_LIMIT_CHECK_DURATION = 1 * time.Minute
const RATE_LIMIT_VALUE = 1

// const RATE_LIMIT_CHECK_DURATION = 20 * time.Second
// const RATE_LIMIT_VALUE = 5

var telemetry *Telemetry
var once sync.Once

func (a *Telemetry) IsSampled() bool {

	random_number := a.minRandInt + rand.Intn(a.maxRandInt-a.minRandInt) + 1

	if (random_number % a.maxRandInt) == 0 {
		return true
	} else {
		return false
	}

}

func (telemetry *Telemetry) CheckSigNozSignals(postData *v3.QueryRangeParamsV3) (bool, bool) {
	signozLogsUsed := false
	signozMetricsUsed := false

	if postData.CompositeQuery.QueryType == v3.QueryTypeBuilder {
		for _, query := range postData.CompositeQuery.BuilderQueries {
			if query.DataSource == v3.DataSourceLogs && len(query.Filters.Items) > 0 {
				signozLogsUsed = true
			} else if query.DataSource == v3.DataSourceMetrics &&
				!strings.Contains(query.AggregateAttribute.Key, "signoz_") &&
				len(query.AggregateAttribute.Key) > 0 {
				signozMetricsUsed = true
			}
		}
	} else if postData.CompositeQuery.QueryType == v3.QueryTypePromQL {
		for _, query := range postData.CompositeQuery.PromQueries {
			if !strings.Contains(query.Query, "signoz_") && len(query.Query) > 0 {
				signozMetricsUsed = true
			}
		}
	} else if postData.CompositeQuery.QueryType == v3.QueryTypeClickHouseSQL {
		for _, query := range postData.CompositeQuery.ClickHouseQueries {
			if strings.Contains(query.Query, "signoz_metrics") && len(query.Query) > 0 {
				signozMetricsUsed = true
			}
		}
	}
	return signozLogsUsed, signozMetricsUsed
}

func (telemetry *Telemetry) AddActiveTracesUser() {
	telemetry.mutex.Lock()
	telemetry.activeUser["traces"] = 1
	telemetry.mutex.Unlock()
}
func (telemetry *Telemetry) AddActiveMetricsUser() {
	telemetry.mutex.Lock()
	telemetry.activeUser["metrics"] = 1
	telemetry.mutex.Unlock()
}
func (telemetry *Telemetry) AddActiveLogsUser() {
	telemetry.mutex.Lock()
	telemetry.activeUser["logs"] = 1
	telemetry.mutex.Unlock()
}

type Telemetry struct {
	operator      analytics.Client
	saasOperator  analytics.Client
	phOperator    ph.Client
	ipAddress     string
	userEmail     string
	isEnabled     bool
	isAnonymous   bool
	distinctId    string
	reader        interfaces.Reader
	companyDomain string
	minRandInt    int
	maxRandInt    int
	rateLimits    map[string]int8
	activeUser    map[string]int8
	countUsers    int8
	mutex         sync.RWMutex
}

func createTelemetry() {
	// Do not do anything in CI (not even resolving the outbound IP address)
	if testing.Testing() {
		telemetry = &Telemetry{
			isEnabled: false,
		}
		return
	}

	telemetry = &Telemetry{
		operator:   analytics.New(api_key),
		phOperator: ph.New(ph_api_key),
		ipAddress:  getOutboundIP(),
		rateLimits: make(map[string]int8),
		activeUser: make(map[string]int8),
	}
	telemetry.minRandInt = 0
	telemetry.maxRandInt = int(1 / DEFAULT_SAMPLING)

	rand.Seed(time.Now().UnixNano())

	data := map[string]interface{}{}

	telemetry.SetTelemetryEnabled(constants.IsTelemetryEnabled())
	telemetry.SendEvent(TELEMETRY_EVENT_HEART_BEAT, data, "")

	ticker := time.NewTicker(HEART_BEAT_DURATION)
	activeUserTicker := time.NewTicker(ACTIVE_USER_DURATION)

	rateLimitTicker := time.NewTicker(RATE_LIMIT_CHECK_DURATION)

	go func() {
		for {
			select {
			case <-rateLimitTicker.C:
				telemetry.rateLimits = make(map[string]int8)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-activeUserTicker.C:
				if (telemetry.activeUser["traces"] != 0) || (telemetry.activeUser["metrics"] != 0) || (telemetry.activeUser["logs"] != 0) {
					telemetry.activeUser["any"] = 1
				}
				telemetry.SendEvent(TELEMETRY_EVENT_ACTIVE_USER, map[string]interface{}{"traces": telemetry.activeUser["traces"], "metrics": telemetry.activeUser["metrics"], "logs": telemetry.activeUser["logs"], "any": telemetry.activeUser["any"]}, "")
				telemetry.activeUser = map[string]int8{"traces": 0, "metrics": 0, "logs": 0, "any": 0}

			case <-ticker.C:

				tagsInfo, _ := telemetry.reader.GetTagsInfoInLastHeartBeatInterval(context.Background())

				if len(tagsInfo.Env) != 0 {
					telemetry.SendEvent(TELEMETRY_EVENT_ENVIRONMENT, map[string]interface{}{"value": tagsInfo.Env}, "")
				}

				for language, _ := range tagsInfo.Languages {
					telemetry.SendEvent(TELEMETRY_EVENT_LANGUAGE, map[string]interface{}{"language": language}, "")
				}

				totalSpans, _ := telemetry.reader.GetTotalSpans(context.Background())
				spansInLastHeartBeatInterval, _ := telemetry.reader.GetSpansInLastHeartBeatInterval(context.Background())
				getSamplesInfoInLastHeartBeatInterval, _ := telemetry.reader.GetSamplesInfoInLastHeartBeatInterval(context.Background())
				tsInfo, _ := telemetry.reader.GetTimeSeriesInfo(context.Background())

				getLogsInfoInLastHeartBeatInterval, _ := telemetry.reader.GetLogsInfoInLastHeartBeatInterval(context.Background())

				traceTTL, _ := telemetry.reader.GetTTL(context.Background(), &model.GetTTLParams{Type: constants.TraceTTL})
				metricsTTL, _ := telemetry.reader.GetTTL(context.Background(), &model.GetTTLParams{Type: constants.MetricsTTL})
				logsTTL, _ := telemetry.reader.GetTTL(context.Background(), &model.GetTTLParams{Type: constants.LogsTTL})

				data := map[string]interface{}{
					"totalSpans":                            totalSpans,
					"spansInLastHeartBeatInterval":          spansInLastHeartBeatInterval,
					"getSamplesInfoInLastHeartBeatInterval": getSamplesInfoInLastHeartBeatInterval,
					"getLogsInfoInLastHeartBeatInterval":    getLogsInfoInLastHeartBeatInterval,
					"countUsers":                            telemetry.countUsers,
					"metricsTTLStatus":                      metricsTTL.Status,
					"tracesTTLStatus":                       traceTTL.Status,
					"logsTTLStatus":                         logsTTL.Status,
				}
				for key, value := range tsInfo {
					data[key] = value
				}
				telemetry.SendEvent(TELEMETRY_EVENT_HEART_BEAT, data, "")

				getDistributedInfoInLastHeartBeatInterval, _ := telemetry.reader.GetDistributedInfoInLastHeartBeatInterval(context.Background())
				telemetry.SendEvent(TELEMETRY_EVENT_DISTRIBUTED, getDistributedInfoInLastHeartBeatInterval, "")

			}
		}
	}()

}

// Get preferred outbound ip of this machine
func getOutboundIP() string {

	ip := []byte(IP_NOT_FOUND_PLACEHOLDER)
	resp, err := http.Get("https://api.ipify.org?format=text")

	if err != nil {
		return string(ip)
	}

	defer resp.Body.Close()
	if err == nil {
		ipBody, err := io.ReadAll(resp.Body)
		if err == nil {
			ip = ipBody
		}
	}

	return string(ip)
}

func (a *Telemetry) IdentifyUser(user *model.User) {
	if user.Email == DEFAULT_CLOUD_EMAIL {
		return
	}
	a.SetCompanyDomain(user.Email)
	a.SetUserEmail(user.Email)
	if !a.isTelemetryEnabled() || a.isTelemetryAnonymous() {
		return
	}
	if a.saasOperator != nil {
		a.saasOperator.Enqueue(analytics.Identify{
			UserId: a.userEmail,
			Traits: analytics.NewTraits().SetName(user.Name).SetEmail(user.Email),
		})
		a.saasOperator.Enqueue(analytics.Group{
			UserId:  a.userEmail,
			GroupId: a.getCompanyDomain(),
			Traits:  analytics.NewTraits().Set("company_domain", a.getCompanyDomain()),
		})
	}

	a.operator.Enqueue(analytics.Identify{
		UserId: a.ipAddress,
		Traits: analytics.NewTraits().SetName(user.Name).SetEmail(user.Email).Set("ip", a.ipAddress),
	})
	// Updating a groups properties
	a.phOperator.Enqueue(ph.GroupIdentify{
		Type: "companyDomain",
		Key:  a.getCompanyDomain(),
		Properties: ph.NewProperties().
			Set("companyDomain", a.getCompanyDomain()),
	})

}

func (a *Telemetry) SetCountUsers(countUsers int8) {
	a.countUsers = countUsers
}

func (a *Telemetry) SetUserEmail(email string) {
	a.userEmail = email
}

func (a *Telemetry) GetUserEmail() string {
	return a.userEmail
}

func (a *Telemetry) SetSaasOperator(saasOperatorKey string) {
	if saasOperatorKey == "" {
		return
	}
	a.saasOperator = analytics.New(saasOperatorKey)
}

func (a *Telemetry) SetCompanyDomain(email string) {

	email_split := strings.Split(email, "@")
	if len(email_split) != 2 {
		a.companyDomain = email
	}
	a.companyDomain = email_split[1]

}

func (a *Telemetry) getCompanyDomain() string {
	return a.companyDomain
}

func (a *Telemetry) checkEvents(event string) bool {
	sendEvent := true
	if event == TELEMETRY_EVENT_USER && a.isTelemetryAnonymous() {
		sendEvent = false
	}
	return sendEvent
}

func (a *Telemetry) SendEvent(event string, data map[string]interface{}, userEmail string, opts ...bool) {

	// ignore telemetry for default user
	if userEmail == DEFAULT_CLOUD_EMAIL {
		return
	}

	if userEmail != "" {
		a.SetUserEmail(userEmail)
	}
	rateLimitFlag := true
	if len(opts) > 0 {
		rateLimitFlag = opts[0]
	}

	if !a.isTelemetryEnabled() {
		return
	}

	ok := a.checkEvents(event)
	if !ok {
		return
	}

	// drop events with properties matching
	if ignoreEvents(event, data) {
		return
	}

	if rateLimitFlag {
		telemetry.mutex.Lock()
		limit := a.rateLimits[event]
		if limit < RATE_LIMIT_VALUE {
			a.rateLimits[event] += 1
			telemetry.mutex.Unlock()
		} else {
			telemetry.mutex.Unlock()
			return
		}
	}

	// zap.S().Info(data)
	properties := analytics.NewProperties()
	properties.Set("version", version.GetVersion())
	properties.Set("deploymentType", getDeploymentType())
	properties.Set("companyDomain", a.getCompanyDomain())

	for k, v := range data {
		properties.Set(k, v)
	}

	userId := a.ipAddress
	if a.isTelemetryAnonymous() || userId == IP_NOT_FOUND_PLACEHOLDER {
		userId = a.GetDistinctId()
	}

	// check if event is part of SAAS_EVENTS_LIST
	_, isSaaSEvent := SAAS_EVENTS_LIST[event]

	if a.saasOperator != nil && a.GetUserEmail() != "" && isSaaSEvent {
		a.saasOperator.Enqueue(analytics.Track{
			Event:      event,
			UserId:     a.GetUserEmail(),
			Properties: properties,
			Context: &analytics.Context{
				Extra: map[string]interface{}{
					"groupId": a.getCompanyDomain(),
				},
			},
		})
	}

	a.operator.Enqueue(analytics.Track{
		Event:      event,
		UserId:     userId,
		Properties: properties,
	})

	if event == TELEMETRY_EVENT_NUMBER_OF_SERVICES {

		a.phOperator.Enqueue(ph.Capture{
			DistinctId: userId,
			Event:      TELEMETRY_EVENT_NUMBER_OF_SERVICES_PH,
			Properties: ph.Properties(properties),
			Groups: ph.NewGroups().
				Set("companyDomain", a.getCompanyDomain()),
		})

	}
	if event == TELEMETRY_EVENT_ACTIVE_USER {

		a.phOperator.Enqueue(ph.Capture{
			DistinctId: userId,
			Event:      TELEMETRY_EVENT_ACTIVE_USER_PH,
			Properties: ph.Properties(properties),
			Groups: ph.NewGroups().
				Set("companyDomain", a.getCompanyDomain()),
		})

	}
}

func (a *Telemetry) GetDistinctId() string {
	return a.distinctId
}
func (a *Telemetry) SetDistinctId(distinctId string) {
	a.distinctId = distinctId
}

func (a *Telemetry) isTelemetryAnonymous() bool {
	return a.isAnonymous
}

func (a *Telemetry) SetTelemetryAnonymous(value bool) {
	a.isAnonymous = value
}

func (a *Telemetry) isTelemetryEnabled() bool {
	return a.isEnabled
}

func (a *Telemetry) SetTelemetryEnabled(value bool) {
	a.isEnabled = value
}

func (a *Telemetry) SetReader(reader interfaces.Reader) {
	a.reader = reader
}

func GetInstance() *Telemetry {

	once.Do(func() {
		createTelemetry()
	})

	return telemetry
}

func getDeploymentType() string {
	deploymentType := os.Getenv("DEPLOYMENT_TYPE")
	if deploymentType == "" {
		return "unknown"
	}
	return deploymentType
}
