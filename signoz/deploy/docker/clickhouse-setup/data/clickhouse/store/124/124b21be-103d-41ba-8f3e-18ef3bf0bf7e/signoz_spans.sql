ATTACH TABLE _ UUID '29083384-8eb0-4a13-a8b3-b0b90f662af6'
(
    `timestamp` DateTime64(9) CODEC(DoubleDelta, LZ4),
    `traceID` FixedString(32) CODEC(ZSTD(1)),
    `model` String CODEC(ZSTD(9))
)
ENGINE = MergeTree
PARTITION BY toDate(timestamp)
ORDER BY traceID
TTL toDateTime(timestamp) + toIntervalSecond(1296000)
SETTINGS index_granularity = 1024, ttl_only_drop_parts = 1
