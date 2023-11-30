ATTACH TABLE _ UUID '51b45e48-a063-4851-8a30-dde9233fda3f'
(
    `metric_name` LowCardinality(String),
    `fingerprint` UInt64 CODEC(DoubleDelta, LZ4),
    `timestamp_ms` Int64 CODEC(DoubleDelta, LZ4),
    `value` Float64 CODEC(Gorilla, LZ4)
)
ENGINE = Distributed('cluster', 'signoz_metrics', 'samples_v2', cityHash64(metric_name, fingerprint))
