ATTACH TABLE _ UUID '3f420807-c970-450a-b547-70b3a7c1c525'
(
    `env` LowCardinality(String) DEFAULT 'default',
    `temporality` LowCardinality(String) DEFAULT 'Unspecified',
    `metric_name` LowCardinality(String),
    `fingerprint` UInt64 CODEC(Delta(8), ZSTD(1)),
    `timestamp_ms` Int64 CODEC(Delta(8), ZSTD(1)),
    `labels` String CODEC(ZSTD(5))
)
ENGINE = ReplacingMergeTree
PARTITION BY toDate(timestamp_ms / 1000)
ORDER BY (env, temporality, metric_name, fingerprint)
SETTINGS index_granularity = 8192
