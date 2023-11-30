ATTACH TABLE _ UUID '660e2dad-4baa-4d0c-a509-de1a0a35ebda'
(
    `env` LowCardinality(String) DEFAULT 'default',
    `temporality` LowCardinality(String) DEFAULT 'Unspecified',
    `metric_name` LowCardinality(String),
    `fingerprint` UInt64 CODEC(Delta(8), ZSTD(1)),
    `timestamp_ms` Int64 CODEC(Delta(8), ZSTD(1)),
    `labels` String CODEC(ZSTD(5))
)
ENGINE = Distributed('cluster', 'signoz_metrics', 'time_series_v3', cityHash64(env, temporality, metric_name, fingerprint))
