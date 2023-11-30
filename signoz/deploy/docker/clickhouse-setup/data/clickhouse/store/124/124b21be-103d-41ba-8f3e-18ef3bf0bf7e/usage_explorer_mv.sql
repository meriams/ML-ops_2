ATTACH MATERIALIZED VIEW _ UUID 'd312660b-35c9-4d1c-a0cd-d86330c778c3' TO signoz_traces.usage_explorer
(
    `timestamp` DateTime,
    `service_name` LowCardinality(String),
    `count` UInt64
) AS
SELECT
    toStartOfHour(timestamp) AS timestamp,
    serviceName AS service_name,
    count() AS count
FROM signoz_traces.signoz_index_v2
GROUP BY
    timestamp,
    serviceName
