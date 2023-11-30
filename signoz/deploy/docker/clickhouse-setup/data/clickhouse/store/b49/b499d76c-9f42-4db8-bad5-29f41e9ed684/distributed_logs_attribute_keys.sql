ATTACH TABLE _ UUID '88799789-60a0-4ee9-9a27-09557d8c27b5'
(
    `name` String,
    `datatype` String
)
ENGINE = Distributed('cluster', 'signoz_logs', 'logs_attribute_keys', cityHash64(datatype))
