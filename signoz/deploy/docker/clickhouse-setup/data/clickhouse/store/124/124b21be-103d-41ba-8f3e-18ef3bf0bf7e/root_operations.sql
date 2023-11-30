ATTACH MATERIALIZED VIEW _ UUID 'c1fa0ad8-5585-47d0-9a29-f040fbb23840' TO signoz_traces.top_level_operations
(
    `name` LowCardinality(String),
    `serviceName` LowCardinality(String)
) AS
SELECT DISTINCT
    name,
    serviceName
FROM signoz_traces.signoz_index_v2
WHERE parentSpanID = ''
