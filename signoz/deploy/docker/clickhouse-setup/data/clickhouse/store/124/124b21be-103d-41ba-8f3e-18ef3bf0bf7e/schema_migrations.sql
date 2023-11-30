ATTACH TABLE _ UUID 'd8827582-555a-4073-a6b6-1e6602c6b014'
(
    `version` Int64,
    `dirty` UInt8,
    `sequence` UInt64
)
ENGINE = MergeTree
ORDER BY sequence
SETTINGS index_granularity = 8192
