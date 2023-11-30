ATTACH TABLE _ UUID 'd59dd3a9-60b9-47ec-bba7-b71d7dd6afc7'
(
    `name` LowCardinality(String) CODEC(ZSTD(1)),
    `serviceName` LowCardinality(String) CODEC(ZSTD(1))
)
ENGINE = Distributed('cluster', 'signoz_traces', 'top_level_operations', cityHash64(rand()))
