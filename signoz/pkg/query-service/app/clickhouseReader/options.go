package clickhouseReader

import (
	"context"
	"net/url"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"go.uber.org/zap"
)

type Encoding string

const (
	// EncodingJSON is used for spans encoded as JSON.
	EncodingJSON Encoding = "json"
	// EncodingProto is used for spans encoded as Protobuf.
	EncodingProto Encoding = "protobuf"
)

const (
	defaultDatasource              string        = "tcp://localhost:9000"
	defaultTraceDB                 string        = "signoz_traces"
	defaultOperationsTable         string        = "distributed_signoz_operations"
	defaultIndexTable              string        = "distributed_signoz_index_v2"
	defaultErrorTable              string        = "distributed_signoz_error_index_v2"
	defaultDurationTable           string        = "distributed_durationSort"
	defaultUsageExplorerTable      string        = "distributed_usage_explorer"
	defaultSpansTable              string        = "distributed_signoz_spans"
	defaultDependencyGraphTable    string        = "distributed_dependency_graph_minutes_v2"
	defaultTopLevelOperationsTable string        = "distributed_top_level_operations"
	defaultSpanAttributeTable      string        = "distributed_span_attributes"
	defaultSpanAttributeKeysTable  string        = "distributed_span_attributes_keys"
	defaultLogsDB                  string        = "signoz_logs"
	defaultLogsTable               string        = "distributed_logs"
	defaultLogsLocalTable          string        = "logs"
	defaultLogAttributeKeysTable   string        = "distributed_logs_attribute_keys"
	defaultLogResourceKeysTable    string        = "distributed_logs_resource_keys"
	defaultLogTagAttributeTable    string        = "distributed_tag_attributes"
	defaultLiveTailRefreshSeconds  int           = 5
	defaultWriteBatchDelay         time.Duration = 5 * time.Second
	defaultWriteBatchSize          int           = 10000
	defaultEncoding                Encoding      = EncodingJSON
)

const (
	suffixEnabled         = ".enabled"
	suffixDatasource      = ".datasource"
	suffixOperationsTable = ".operations-table"
	suffixIndexTable      = ".index-table"
	suffixSpansTable      = ".spans-table"
	suffixWriteBatchDelay = ".write-batch-delay"
	suffixWriteBatchSize  = ".write-batch-size"
	suffixEncoding        = ".encoding"
)

// NamespaceConfig is Clickhouse's internal configuration data
type namespaceConfig struct {
	namespace               string
	Enabled                 bool
	Datasource              string
	MaxIdleConns            int
	MaxOpenConns            int
	DialTimeout             time.Duration
	TraceDB                 string
	OperationsTable         string
	IndexTable              string
	DurationTable           string
	UsageExplorerTable      string
	SpansTable              string
	ErrorTable              string
	SpanAttributeTable      string
	SpanAttributeKeysTable  string
	DependencyGraphTable    string
	TopLevelOperationsTable string
	LogsDB                  string
	LogsTable               string
	LogsLocalTable          string
	LogsAttributeKeysTable  string
	LogsResourceKeysTable   string
	LogsTagAttributeTable   string
	LiveTailRefreshSeconds  int
	WriteBatchDelay         time.Duration
	WriteBatchSize          int
	Encoding                Encoding
	Connector               Connector
}

// Connecto defines how to connect to the database
type Connector func(cfg *namespaceConfig) (clickhouse.Conn, error)

func defaultConnector(cfg *namespaceConfig) (clickhouse.Conn, error) {
	ctx := context.Background()
	dsnURL, err := url.Parse(cfg.Datasource)
	if err != nil {
		return nil, err
	}
	options := &clickhouse.Options{
		Addr:         []string{dsnURL.Host},
		MaxOpenConns: cfg.MaxOpenConns,
		MaxIdleConns: cfg.MaxIdleConns,
		DialTimeout:  cfg.DialTimeout,
	}
	if dsnURL.Query().Get("username") != "" {
		auth := clickhouse.Auth{
			Username: dsnURL.Query().Get("username"),
			Password: dsnURL.Query().Get("password"),
		}
		options.Auth = auth
	}
	zap.S().Infof("Connecting to Clickhouse at %s, MaxIdleConns: %d, MaxOpenConns: %d, DialTimeout: %s", dsnURL.Host, options.MaxIdleConns, options.MaxOpenConns, options.DialTimeout)
	db, err := clickhouse.Open(options)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

// Options store storage plugin related configs
type Options struct {
	primary *namespaceConfig

	others map[string]*namespaceConfig
}

// NewOptions creates a new Options struct.
func NewOptions(
	datasource string,
	maxIdleConns int,
	maxOpenConns int,
	dialTimeout time.Duration,
	primaryNamespace string,
	otherNamespaces ...string,
) *Options {

	if datasource == "" {
		datasource = defaultDatasource
	}

	options := &Options{
		primary: &namespaceConfig{
			namespace:               primaryNamespace,
			Enabled:                 true,
			Datasource:              datasource,
			MaxIdleConns:            maxIdleConns,
			MaxOpenConns:            maxOpenConns,
			DialTimeout:             dialTimeout,
			TraceDB:                 defaultTraceDB,
			OperationsTable:         defaultOperationsTable,
			IndexTable:              defaultIndexTable,
			ErrorTable:              defaultErrorTable,
			DurationTable:           defaultDurationTable,
			UsageExplorerTable:      defaultUsageExplorerTable,
			SpansTable:              defaultSpansTable,
			SpanAttributeTable:      defaultSpanAttributeTable,
			SpanAttributeKeysTable:  defaultSpanAttributeKeysTable,
			DependencyGraphTable:    defaultDependencyGraphTable,
			TopLevelOperationsTable: defaultTopLevelOperationsTable,
			LogsDB:                  defaultLogsDB,
			LogsTable:               defaultLogsTable,
			LogsLocalTable:          defaultLogsLocalTable,
			LogsAttributeKeysTable:  defaultLogAttributeKeysTable,
			LogsResourceKeysTable:   defaultLogResourceKeysTable,
			LogsTagAttributeTable:   defaultLogTagAttributeTable,
			LiveTailRefreshSeconds:  defaultLiveTailRefreshSeconds,
			WriteBatchDelay:         defaultWriteBatchDelay,
			WriteBatchSize:          defaultWriteBatchSize,
			Encoding:                defaultEncoding,
			Connector:               defaultConnector,
		},
		others: make(map[string]*namespaceConfig, len(otherNamespaces)),
	}

	for _, namespace := range otherNamespaces {
		if namespace == archiveNamespace {
			options.others[namespace] = &namespaceConfig{
				namespace:              namespace,
				Datasource:             datasource,
				TraceDB:                "",
				OperationsTable:        "",
				IndexTable:             "",
				ErrorTable:             "",
				LogsDB:                 "",
				LogsTable:              "",
				LogsLocalTable:         "",
				LogsAttributeKeysTable: "",
				LogsResourceKeysTable:  "",
				LiveTailRefreshSeconds: defaultLiveTailRefreshSeconds,
				WriteBatchDelay:        defaultWriteBatchDelay,
				WriteBatchSize:         defaultWriteBatchSize,
				Encoding:               defaultEncoding,
				Connector:              defaultConnector,
			}
		} else {
			options.others[namespace] = &namespaceConfig{namespace: namespace}
		}
	}

	return options
}

// GetPrimary returns the primary namespace configuration
func (opt *Options) getPrimary() *namespaceConfig {
	return opt.primary
}
