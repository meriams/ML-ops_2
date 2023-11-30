ATTACH TABLE _ UUID 'ca4e7ccf-8783-41eb-ba32-152b03f5e97c'
(
    `tenant` String,
    `collector_id` String,
    `exporter_id` String,
    `timestamp` DateTime,
    `data` String
)
ENGINE = Distributed('cluster', 'signoz_logs', 'usage', cityHash64(rand()))
