ATTACH TABLE _ UUID '5094749f-2029-4b67-9d34-5a1509057434'
(
    `src` LowCardinality(String) CODEC(ZSTD(1)),
    `dest` LowCardinality(String) CODEC(ZSTD(1)),
    `duration_quantiles_state` AggregateFunction(quantiles(0.5, 0.75, 0.9, 0.95, 0.99), Float64) CODEC(Default),
    `error_count` SimpleAggregateFunction(sum, UInt64) CODEC(T64, ZSTD(1)),
    `total_count` SimpleAggregateFunction(sum, UInt64) CODEC(T64, ZSTD(1)),
    `timestamp` DateTime CODEC(DoubleDelta, LZ4),
    `deployment_environment` LowCardinality(String) CODEC(ZSTD(1)),
    `k8s_cluster_name` LowCardinality(String) CODEC(ZSTD(1)),
    `k8s_namespace_name` LowCardinality(String) CODEC(ZSTD(1))
)
ENGINE = Distributed('cluster', 'signoz_traces', 'dependency_graph_minutes_v2', cityHash64(rand()))
