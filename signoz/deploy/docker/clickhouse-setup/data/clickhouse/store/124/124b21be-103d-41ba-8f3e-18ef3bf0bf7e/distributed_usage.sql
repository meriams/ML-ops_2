ATTACH TABLE _ UUID '60659ab3-0368-4883-8528-1d7460fd3056'
(
    `tenant` String,
    `collector_id` String,
    `exporter_id` String,
    `timestamp` DateTime,
    `data` String
)
ENGINE = Distributed('cluster', 'signoz_traces', 'usage', cityHash64(rand()))
