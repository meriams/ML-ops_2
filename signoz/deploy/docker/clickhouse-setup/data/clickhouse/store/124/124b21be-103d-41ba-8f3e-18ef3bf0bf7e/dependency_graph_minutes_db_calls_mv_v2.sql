ATTACH MATERIALIZED VIEW _ UUID 'c0e149e3-421a-4906-884d-d5a4cd0c427b' TO signoz_traces.dependency_graph_minutes_v2
(
    `src` LowCardinality(String),
    `dest` String,
    `duration_quantiles_state` AggregateFunction(quantiles(0.5, 0.75, 0.9, 0.95, 0.99), Float64),
    `error_count` UInt64,
    `total_count` UInt64,
    `timestamp` DateTime,
    `deployment_environment` String,
    `k8s_cluster_name` String,
    `k8s_namespace_name` String
) AS
SELECT
    serviceName AS src,
    tagMap['db.system'] AS dest,
    quantilesState(0.5, 0.75, 0.9, 0.95, 0.99)(toFloat64(durationNano)) AS duration_quantiles_state,
    countIf(statusCode = 2) AS error_count,
    count(*) AS total_count,
    toStartOfMinute(timestamp) AS timestamp,
    resourceTagsMap['deployment.environment'] AS deployment_environment,
    resourceTagsMap['k8s.cluster.name'] AS k8s_cluster_name,
    resourceTagsMap['k8s.namespace.name'] AS k8s_namespace_name
FROM signoz_traces.signoz_index_v2
WHERE (dest != '') AND (kind != 2)
GROUP BY
    timestamp,
    src,
    dest,
    deployment_environment,
    k8s_cluster_name,
    k8s_namespace_name
