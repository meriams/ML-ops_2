ATTACH TABLE _ UUID '32d126be-f59d-4b05-ab66-2e1691d89d20'
(
    `timestamp` UInt64 CODEC(DoubleDelta, LZ4),
    `observed_timestamp` UInt64 CODEC(DoubleDelta, LZ4),
    `id` String CODEC(ZSTD(1)),
    `trace_id` String CODEC(ZSTD(1)),
    `span_id` String CODEC(ZSTD(1)),
    `trace_flags` UInt32,
    `severity_text` LowCardinality(String) CODEC(ZSTD(1)),
    `severity_number` UInt8,
    `body` String CODEC(ZSTD(2)),
    `resources_string_key` Array(String) CODEC(ZSTD(1)),
    `resources_string_value` Array(String) CODEC(ZSTD(1)),
    `attributes_string_key` Array(String) CODEC(ZSTD(1)),
    `attributes_string_value` Array(String) CODEC(ZSTD(1)),
    `attributes_int64_key` Array(String) CODEC(ZSTD(1)),
    `attributes_int64_value` Array(Int64) CODEC(ZSTD(1)),
    `attributes_float64_key` Array(String) CODEC(ZSTD(1)),
    `attributes_float64_value` Array(Float64) CODEC(ZSTD(1)),
    `attributes_bool_key` Array(String) CODEC(ZSTD(1)),
    `attributes_bool_value` Array(Bool) CODEC(ZSTD(1)),
    INDEX body_idx body TYPE tokenbf_v1(10240, 3, 0) GRANULARITY 4,
    INDEX id_minmax id TYPE minmax GRANULARITY 1,
    INDEX severity_number_idx severity_number TYPE set(25) GRANULARITY 4,
    INDEX severity_text_idx severity_text TYPE set(25) GRANULARITY 4,
    INDEX trace_flags_idx trace_flags TYPE bloom_filter GRANULARITY 4
)
ENGINE = MergeTree
PARTITION BY toDate(timestamp / 1000000000)
ORDER BY (timestamp, id)
TTL toDateTime(timestamp / 1000000000) + toIntervalSecond(1296000)
SETTINGS index_granularity = 8192, ttl_only_drop_parts = 1
