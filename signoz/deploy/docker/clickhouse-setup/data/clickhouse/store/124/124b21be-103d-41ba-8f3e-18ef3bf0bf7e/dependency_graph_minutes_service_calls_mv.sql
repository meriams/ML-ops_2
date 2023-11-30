ATTACH MATERIALIZED VIEW _ UUID 'c755324d-8922-4912-9e37-e1f1e34094aa' TO signoz_traces.dependency_graph_minutes
(
    `src` LowCardinality(String),
    `dest` LowCardinality(String),
    `duration_quantiles_state` AggregateFunction(quantiles(0.5, 0.75, 0.9, 0.95, 0.99), Float64),
    `error_count` UInt64,
    `total_count` UInt64,
    `timestamp` DateTime
) AS
SELECT
    A.serviceName AS src,
    B.serviceName AS dest,
    quantilesState(0.5, 0.75, 0.9, 0.95, 0.99)(toFloat64(B.durationNano)) AS duration_quantiles_state,
    countIf(B.statusCode = 2) AS error_count,
    count(*) AS total_count,
    toStartOfMinute(B.timestamp) AS timestamp
FROM signoz_traces.signoz_index_v2 AS A, signoz_traces.signoz_index_v2 AS B
WHERE (A.serviceName != B.serviceName) AND (A.spanID = B.parentSpanID)
GROUP BY
    timestamp,
    src,
    dest
