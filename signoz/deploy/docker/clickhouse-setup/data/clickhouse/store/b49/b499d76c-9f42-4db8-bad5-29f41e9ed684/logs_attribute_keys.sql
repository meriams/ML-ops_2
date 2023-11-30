ATTACH TABLE _ UUID 'd9a7b513-ced9-4ca2-abe3-0017447270bb'
(
    `name` String,
    `datatype` String
)
ENGINE = ReplacingMergeTree
ORDER BY (name, datatype)
SETTINGS index_granularity = 8192
