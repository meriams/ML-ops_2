ATTACH TABLE _ UUID '45da05a0-dab8-4406-a079-8e0f3d382435'
(
    `metric_name` LowCardinality(String),
    `fingerprint` UInt64 CODEC(DoubleDelta, LZ4),
    `timestamp_ms` Int64 CODEC(DoubleDelta, LZ4),
    `labels` String CODEC(ZSTD(5)),
    `temporality` LowCardinality(String) DEFAULT 'Unspecified' CODEC(ZSTD(5))
)
ENGINE = Distributed('cluster', 'signoz_metrics', 'time_series_v2', cityHash64(metric_name, fingerprint))
