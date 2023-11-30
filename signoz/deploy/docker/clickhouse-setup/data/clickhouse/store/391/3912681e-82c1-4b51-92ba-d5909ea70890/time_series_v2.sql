ATTACH TABLE _ UUID '0624ba10-aefb-4b87-8fd8-fd089aa9dcf6'
(
    `metric_name` LowCardinality(String),
    `fingerprint` UInt64 CODEC(DoubleDelta, LZ4),
    `timestamp_ms` Int64 CODEC(DoubleDelta, LZ4),
    `labels` String CODEC(ZSTD(5)),
    `temporality` LowCardinality(String) DEFAULT 'Unspecified' CODEC(ZSTD(5)),
    INDEX temporality_index temporality TYPE SET(3) GRANULARITY 1
)
ENGINE = ReplacingMergeTree
PARTITION BY toDate(timestamp_ms / 1000)
ORDER BY (metric_name, fingerprint)
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1
