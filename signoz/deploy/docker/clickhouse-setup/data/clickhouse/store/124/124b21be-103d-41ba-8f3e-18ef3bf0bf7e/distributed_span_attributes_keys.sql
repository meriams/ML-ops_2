ATTACH TABLE _ UUID 'c5301f64-8e5e-466d-9459-e77eb91bc39d'
(
    `tagKey` LowCardinality(String) CODEC(ZSTD(1)),
    `tagType` Enum8('tag' = 1, 'resource' = 2) CODEC(ZSTD(1)),
    `dataType` Enum8('string' = 1, 'bool' = 2, 'float64' = 3) CODEC(ZSTD(1)),
    `isColumn` Bool CODEC(ZSTD(1))
)
ENGINE = Distributed('cluster', 'signoz_traces', 'span_attributes_keys', cityHash64(rand()))
