ATTACH TABLE _ UUID '7279d2b6-237b-4d93-bf16-06919af19113'
(
    `version` Int64,
    `dirty` UInt8,
    `sequence` UInt64
)
ENGINE = MergeTree
ORDER BY sequence
SETTINGS index_granularity = 8192
