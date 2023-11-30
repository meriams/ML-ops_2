ATTACH TABLE _ UUID 'e6dd0cab-a0d1-4eaf-b6dc-b37a793de7b2'
(
    `timestamp` DateTime CODEC(ZSTD(1)),
    `tagKey` LowCardinality(String) CODEC(ZSTD(1)),
    `tagType` Enum8('tag' = 1, 'resource' = 2) CODEC(ZSTD(1)),
    `tagDataType` Enum8('string' = 1, 'bool' = 2, 'int64' = 3, 'float64' = 4) CODEC(ZSTD(1)),
    `stringTagValue` String CODEC(ZSTD(1)),
    `int64TagValue` Nullable(Int64) CODEC(ZSTD(1)),
    `float64TagValue` Nullable(Float64) CODEC(ZSTD(1))
)
ENGINE = Distributed('cluster', 'signoz_logs', 'tag_attributes', rand())
