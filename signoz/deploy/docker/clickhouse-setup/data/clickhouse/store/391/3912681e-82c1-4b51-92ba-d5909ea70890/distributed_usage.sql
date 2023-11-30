ATTACH TABLE _ UUID '4d5c6889-fd98-4f35-af55-11428c320145'
(
    `tenant` String,
    `collector_id` String,
    `exporter_id` String,
    `timestamp` DateTime,
    `data` String
)
ENGINE = Distributed('cluster', 'signoz_metrics', 'usage', cityHash64(rand()))
