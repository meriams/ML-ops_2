ATTACH TABLE _ UUID '1fb6055e-3645-4445-a912-9cf4996ac072'
(
    `version` Int64,
    `dirty` UInt8,
    `sequence` UInt64
)
ENGINE = MergeTree
ORDER BY sequence
SETTINGS index_granularity = 8192
