ATTACH MATERIALIZED VIEW _ UUID 'a41e8790-84c5-47f6-9336-03dcba1ce4c4' TO signoz_logs.logs_attribute_keys
(
    `name` String,
    `datatype` String
) AS
SELECT DISTINCT
    arrayJoin(attributes_float64_key) AS name,
    'Float64' AS datatype
FROM signoz_logs.logs
ORDER BY name ASC
