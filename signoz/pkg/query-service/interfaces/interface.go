package interfaces

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/storage"
	"github.com/prometheus/prometheus/util/stats"
	am "go.signoz.io/signoz/pkg/query-service/integrations/alertManager"
	"go.signoz.io/signoz/pkg/query-service/model"
	v3 "go.signoz.io/signoz/pkg/query-service/model/v3"
)

type Reader interface {
	GetChannel(id string) (*model.ChannelItem, *model.ApiError)
	GetChannels() (*[]model.ChannelItem, *model.ApiError)
	DeleteChannel(id string) *model.ApiError
	CreateChannel(receiver *am.Receiver) (*am.Receiver, *model.ApiError)
	EditChannel(receiver *am.Receiver, id string) (*am.Receiver, *model.ApiError)

	GetInstantQueryMetricsResult(ctx context.Context, query *model.InstantQueryMetricsParams) (*promql.Result, *stats.QueryStats, *model.ApiError)
	GetQueryRangeResult(ctx context.Context, query *model.QueryRangeParams) (*promql.Result, *stats.QueryStats, *model.ApiError)
	GetServiceOverview(ctx context.Context, query *model.GetServiceOverviewParams, skipConfig *model.SkipConfig) (*[]model.ServiceOverviewItem, *model.ApiError)
	GetTopLevelOperations(ctx context.Context, skipConfig *model.SkipConfig) (*map[string][]string, *model.ApiError)
	GetServices(ctx context.Context, query *model.GetServicesParams, skipConfig *model.SkipConfig) (*[]model.ServiceItem, *model.ApiError)
	GetTopOperations(ctx context.Context, query *model.GetTopOperationsParams) (*[]model.TopOperationsItem, *model.ApiError)
	GetUsage(ctx context.Context, query *model.GetUsageParams) (*[]model.UsageItem, error)
	GetServicesList(ctx context.Context) (*[]string, error)
	GetDependencyGraph(ctx context.Context, query *model.GetServicesParams) (*[]model.ServiceMapDependencyResponseItem, error)

	GetTTL(ctx context.Context, ttlParams *model.GetTTLParams) (*model.GetTTLResponseItem, *model.ApiError)

	// GetDisks returns a list of disks configured in the underlying DB. It is supported by
	// clickhouse only.
	GetDisks(ctx context.Context) (*[]model.DiskItem, *model.ApiError)
	GetSpanFilters(ctx context.Context, query *model.SpanFilterParams) (*model.SpanFiltersResponse, *model.ApiError)
	GetTraceAggregateAttributes(ctx context.Context, req *v3.AggregateAttributeRequest) (*v3.AggregateAttributeResponse, error)
	GetTraceAttributeKeys(ctx context.Context, req *v3.FilterAttributeKeyRequest) (*v3.FilterAttributeKeyResponse, error)
	GetTraceAttributeValues(ctx context.Context, req *v3.FilterAttributeValueRequest) (*v3.FilterAttributeValueResponse, error)
	GetSpanAttributeKeys(ctx context.Context) (map[string]v3.AttributeKey, error)
	GetTagFilters(ctx context.Context, query *model.TagFilterParams) (*model.TagFilters, *model.ApiError)
	GetTagValues(ctx context.Context, query *model.TagFilterParams) (*model.TagValues, *model.ApiError)
	GetFilteredSpans(ctx context.Context, query *model.GetFilteredSpansParams) (*model.GetFilterSpansResponse, *model.ApiError)
	GetFilteredSpansAggregates(ctx context.Context, query *model.GetFilteredSpanAggregatesParams) (*model.GetFilteredSpansAggregatesResponse, *model.ApiError)

	ListErrors(ctx context.Context, params *model.ListErrorsParams) (*[]model.Error, *model.ApiError)
	CountErrors(ctx context.Context, params *model.CountErrorsParams) (uint64, *model.ApiError)
	GetErrorFromErrorID(ctx context.Context, params *model.GetErrorParams) (*model.ErrorWithSpan, *model.ApiError)
	GetErrorFromGroupID(ctx context.Context, params *model.GetErrorParams) (*model.ErrorWithSpan, *model.ApiError)
	GetNextPrevErrorIDs(ctx context.Context, params *model.GetErrorParams) (*model.NextPrevErrorIDs, *model.ApiError)

	// Search Interfaces
	SearchTraces(ctx context.Context, traceID string, spanId string, levelUp int, levelDown int, spanLimit int, smartTraceAlgorithm func(payload []model.SearchSpanResponseItem, targetSpanId string, levelUp int, levelDown int, spanLimit int) ([]model.SearchSpansResult, error)) (*[]model.SearchSpansResult, error)

	// Setter Interfaces
	SetTTL(ctx context.Context, ttlParams *model.TTLParams) (*model.SetTTLResponseItem, *model.ApiError)

	FetchTemporality(ctx context.Context, metricNames []string) (map[string]map[v3.Temporality]bool, error)
	GetMetricAutocompleteMetricNames(ctx context.Context, matchText string, limit int) (*[]string, *model.ApiError)
	GetMetricAutocompleteTagKey(ctx context.Context, params *model.MetricAutocompleteTagParams) (*[]string, *model.ApiError)
	GetMetricAutocompleteTagValue(ctx context.Context, params *model.MetricAutocompleteTagParams) (*[]string, *model.ApiError)
	GetMetricResult(ctx context.Context, query string) ([]*model.Series, error)
	GetMetricResultEE(ctx context.Context, query string) ([]*model.Series, string, error)
	GetMetricAggregateAttributes(ctx context.Context, req *v3.AggregateAttributeRequest) (*v3.AggregateAttributeResponse, error)
	GetMetricAttributeKeys(ctx context.Context, req *v3.FilterAttributeKeyRequest) (*v3.FilterAttributeKeyResponse, error)
	GetMetricAttributeValues(ctx context.Context, req *v3.FilterAttributeValueRequest) (*v3.FilterAttributeValueResponse, error)

	// QB V3 metrics/traces/logs
	GetTimeSeriesResultV3(ctx context.Context, query string) ([]*v3.Series, error)
	GetListResultV3(ctx context.Context, query string) ([]*v3.Row, error)
	LiveTailLogsV3(ctx context.Context, query string, timestampStart uint64, idStart string, client *v3.LogsLiveTailClient)

	GetTotalSpans(ctx context.Context) (uint64, error)
	GetSpansInLastHeartBeatInterval(ctx context.Context) (uint64, error)
	GetTimeSeriesInfo(ctx context.Context) (map[string]interface{}, error)
	GetSamplesInfoInLastHeartBeatInterval(ctx context.Context) (uint64, error)
	GetLogsInfoInLastHeartBeatInterval(ctx context.Context) (uint64, error)
	GetTagsInfoInLastHeartBeatInterval(ctx context.Context) (*model.TagsInfo, error)
	GetDistributedInfoInLastHeartBeatInterval(ctx context.Context) (map[string]interface{}, error)
	// Logs
	GetLogFields(ctx context.Context) (*model.GetFieldsResponse, *model.ApiError)
	UpdateLogField(ctx context.Context, field *model.UpdateField) *model.ApiError
	GetLogs(ctx context.Context, params *model.LogsFilterParams) (*[]model.SignozLog, *model.ApiError)
	TailLogs(ctx context.Context, client *model.LogsTailClient)
	AggregateLogs(ctx context.Context, params *model.LogsAggregateParams) (*model.GetLogsAggregatesResponse, *model.ApiError)
	GetLogAttributeKeys(ctx context.Context, req *v3.FilterAttributeKeyRequest) (*v3.FilterAttributeKeyResponse, error)
	GetLogAttributeValues(ctx context.Context, req *v3.FilterAttributeValueRequest) (*v3.FilterAttributeValueResponse, error)
	GetLogAggregateAttributes(ctx context.Context, req *v3.AggregateAttributeRequest) (*v3.AggregateAttributeResponse, error)

	// Connection needed for rules, not ideal but required
	GetConn() clickhouse.Conn
	GetQueryEngine() *promql.Engine
	GetFanoutStorage() *storage.Storage

	QueryDashboardVars(ctx context.Context, query string) (*model.DashboardVar, error)
	CheckClickHouse(ctx context.Context) error

	GetLatencyMetricMetadata(context.Context, string, bool) (*v3.LatencyMetricMetadataResponse, error)
}

type Querier interface {
	QueryRange(context.Context, *v3.QueryRangeParamsV3, map[string]v3.AttributeKey) ([]*v3.Result, error, map[string]string)

	// test helpers
	QueriesExecuted() []string
}
