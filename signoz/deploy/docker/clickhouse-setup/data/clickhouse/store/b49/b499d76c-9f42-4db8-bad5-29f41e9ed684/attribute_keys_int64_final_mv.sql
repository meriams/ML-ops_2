ATTACH MATERIALIZED VIEW _ UUID '5ad8f8b9-9794-4242-a265-ff46e0e43120' TO signoz_logs.logs_attribute_keys
(
    `name` String,
    `datatype` String
) AS
SELECT DISTINCT
    arrayJoin(attributes_int64_key) AS name,
    'Int64' AS datatype
FROM signoz_logs.logs
ORDER BY name ASC
