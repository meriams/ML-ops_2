ATTACH TABLE _ UUID '1c2e8887-b4c3-4b83-abc0-b6a5df78f02c'
(
    `name` String,
    `datatype` String
)
ENGINE = Distributed('cluster', 'signoz_logs', 'logs_resource_keys', cityHash64(datatype))
