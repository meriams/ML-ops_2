ATTACH MATERIALIZED VIEW _ UUID '46f22870-380f-452f-be72-eaf0015c4195' TO signoz_traces.dependency_graph_minutes_v2
(
    `src` LowCardinality(String),
    `dest` LowCardinality(String),
    `duration_quantiles_state` AggregateFunction(quantiles(0.5, 0.75, 0.9, 0.95, 0.99), Float64),
    `error_count` UInt64,
    `total_count` UInt64,
    `timestamp` DateTime,
    `deployment_environment` String,
    `k8s_cluster_name` String,
    `k8s_namespace_name` String
) AS
SELECT
    A.serviceName AS src,
    B.serviceName AS dest,
    quantilesState(0.5, 0.75, 0.9, 0.95, 0.99)(toFloat64(B.durationNano)) AS duration_quantiles_state,
    countIf(B.statusCode = 2) AS error_count,
    count(*) AS total_count,
    toStartOfMinute(B.timestamp) AS timestamp,
    B.resourceTagsMap['deployment.environment'] AS deployment_environment,
    B.resourceTagsMap['k8s.cluster.name'] AS k8s_cluster_name,
    B.resourceTagsMap['k8s.namespace.name'] AS k8s_namespace_name
FROM signoz_traces.signoz_index_v2 AS A, signoz_traces.signoz_index_v2 AS B
WHERE (A.serviceName != B.serviceName) AND (A.spanID = B.parentSpanID)
GROUP BY
    timestamp,
    src,
    dest,
    deployment_environment,
    k8s_cluster_name,
    k8s_namespace_name
