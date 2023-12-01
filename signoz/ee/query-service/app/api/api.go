package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.signoz.io/signoz/ee/query-service/dao"
	"go.signoz.io/signoz/ee/query-service/interfaces"
	"go.signoz.io/signoz/ee/query-service/license"
	"go.signoz.io/signoz/ee/query-service/usage"
	baseapp "go.signoz.io/signoz/pkg/query-service/app"
	"go.signoz.io/signoz/pkg/query-service/app/logparsingpipeline"
	"go.signoz.io/signoz/pkg/query-service/cache"
	baseint "go.signoz.io/signoz/pkg/query-service/interfaces"
	basemodel "go.signoz.io/signoz/pkg/query-service/model"
	rules "go.signoz.io/signoz/pkg/query-service/rules"
	"go.signoz.io/signoz/pkg/query-service/version"
)

type APIHandlerOptions struct {
	DataConnector                 interfaces.DataConnector
	SkipConfig                    *basemodel.SkipConfig
	PreferDelta                   bool
	PreferSpanMetrics             bool
	MaxIdleConns                  int
	MaxOpenConns                  int
	DialTimeout                   time.Duration
	AppDao                        dao.ModelDao
	RulesManager                  *rules.Manager
	UsageManager                  *usage.Manager
	FeatureFlags                  baseint.FeatureLookup
	LicenseManager                *license.Manager
	LogsParsingPipelineController *logparsingpipeline.LogParsingPipelineController
	Cache                         cache.Cache
	// Querier Influx Interval
	FluxInterval time.Duration
}

type APIHandler struct {
	opts APIHandlerOptions
	baseapp.APIHandler
}

// NewAPIHandler returns an APIHandler
func NewAPIHandler(opts APIHandlerOptions) (*APIHandler, error) {

	baseHandler, err := baseapp.NewAPIHandler(baseapp.APIHandlerOpts{
		Reader:                        opts.DataConnector,
		SkipConfig:                    opts.SkipConfig,
		PerferDelta:                   opts.PreferDelta,
		PreferSpanMetrics:             opts.PreferSpanMetrics,
		MaxIdleConns:                  opts.MaxIdleConns,
		MaxOpenConns:                  opts.MaxOpenConns,
		DialTimeout:                   opts.DialTimeout,
		AppDao:                        opts.AppDao,
		RuleManager:                   opts.RulesManager,
		FeatureFlags:                  opts.FeatureFlags,
		LogsParsingPipelineController: opts.LogsParsingPipelineController,
		Cache:                         opts.Cache,
		FluxInterval:                  opts.FluxInterval,
	})

	if err != nil {
		return nil, err
	}

	ah := &APIHandler{
		opts:       opts,
		APIHandler: *baseHandler,
	}
	return ah, nil
}

func (ah *APIHandler) FF() baseint.FeatureLookup {
	return ah.opts.FeatureFlags
}

func (ah *APIHandler) RM() *rules.Manager {
	return ah.opts.RulesManager
}

func (ah *APIHandler) LM() *license.Manager {
	return ah.opts.LicenseManager
}

func (ah *APIHandler) UM() *usage.Manager {
	return ah.opts.UsageManager
}

func (ah *APIHandler) AppDao() dao.ModelDao {
	return ah.opts.AppDao
}

func (ah *APIHandler) CheckFeature(f string) bool {
	err := ah.FF().CheckFeature(f)
	return err == nil
}

// RegisterRoutes registers routes for this handler on the given router
func (ah *APIHandler) RegisterRoutes(router *mux.Router, am *baseapp.AuthMiddleware) {
	// note: add ee override methods first

	// routes available only in ee version
	router.HandleFunc("/api/v1/licenses",
		am.AdminAccess(ah.listLicenses)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/licenses",
		am.AdminAccess(ah.applyLicense)).
		Methods(http.MethodPost)

	router.HandleFunc("/api/v1/featureFlags",
		am.OpenAccess(ah.getFeatureFlags)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/loginPrecheck",
		am.OpenAccess(ah.precheckLogin)).
		Methods(http.MethodGet)

	// paid plans specific routes
	router.HandleFunc("/api/v1/complete/saml",
		am.OpenAccess(ah.receiveSAML)).
		Methods(http.MethodPost)

	router.HandleFunc("/api/v1/complete/google",
		am.OpenAccess(ah.receiveGoogleAuth)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/orgs/{orgId}/domains",
		am.AdminAccess(ah.listDomainsByOrg)).
		Methods(http.MethodGet)

	router.HandleFunc("/api/v1/domains",
		am.AdminAccess(ah.postDomain)).
		Methods(http.MethodPost)

	router.HandleFunc("/api/v1/domains/{id}",
		am.AdminAccess(ah.putDomain)).
		Methods(http.MethodPut)

	router.HandleFunc("/api/v1/domains/{id}",
		am.AdminAccess(ah.deleteDomain)).
		Methods(http.MethodDelete)

	// base overrides
	router.HandleFunc("/api/v1/version", am.OpenAccess(ah.getVersion)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/invite/{token}", am.OpenAccess(ah.getInvite)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/register", am.OpenAccess(ah.registerUser)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/login", am.OpenAccess(ah.loginUser)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/traces/{traceId}", am.ViewAccess(ah.searchTraces)).Methods(http.MethodGet)
	router.HandleFunc("/api/v2/metrics/query_range", am.ViewAccess(ah.queryRangeMetricsV2)).Methods(http.MethodPost)

	// PAT APIs
	router.HandleFunc("/api/v1/pat", am.OpenAccess(ah.createPAT)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/pat", am.OpenAccess(ah.getPATs)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/pat/{id}", am.OpenAccess(ah.deletePAT)).Methods(http.MethodDelete)

	router.HandleFunc("/api/v1/checkout", am.AdminAccess(ah.checkout)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/billing", am.AdminAccess(ah.getBilling)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/portal", am.AdminAccess(ah.portalSession)).Methods(http.MethodPost)

	router.HandleFunc("/api/v1/dashboards/{uuid}/lock", am.EditAccess(ah.lockDashboard)).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/dashboards/{uuid}/unlock", am.EditAccess(ah.unlockDashboard)).Methods(http.MethodPut)

	router.HandleFunc("/api/v2/licenses",
		am.ViewAccess(ah.listLicensesV2)).
		Methods(http.MethodGet)

	ah.APIHandler.RegisterRoutes(router, am)

}

func (ah *APIHandler) getVersion(w http.ResponseWriter, r *http.Request) {
	version := version.GetVersion()
	versionResponse := basemodel.GetVersionResponse{
		Version:        version,
		EE:             "Y",
		SetupCompleted: ah.SetupCompleted,
	}

	ah.WriteJSON(w, r, versionResponse)
}
