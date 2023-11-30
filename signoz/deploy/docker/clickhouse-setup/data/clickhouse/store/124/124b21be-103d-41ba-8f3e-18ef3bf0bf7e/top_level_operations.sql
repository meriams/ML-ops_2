ATTACH TABLE _ UUID 'f622e4f2-db6a-4b0a-bc0e-f1cf9b4b1853'
(
    `name` LowCardinality(String) CODEC(ZSTD(1)),
    `serviceName` LowCardinality(String) CODEC(ZSTD(1))
)
ENGINE = ReplacingMergeTree
ORDER BY (serviceName, name)
SETTINGS index_granularity = 8192
