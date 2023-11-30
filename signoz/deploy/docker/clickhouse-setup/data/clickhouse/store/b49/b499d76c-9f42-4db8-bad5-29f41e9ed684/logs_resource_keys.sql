ATTACH TABLE _ UUID 'bd42cf9c-0dc4-4056-8390-7d4a3214e7c7'
(
    `name` String,
    `datatype` String
)
ENGINE = ReplacingMergeTree
ORDER BY (name, datatype)
SETTINGS index_granularity = 8192
