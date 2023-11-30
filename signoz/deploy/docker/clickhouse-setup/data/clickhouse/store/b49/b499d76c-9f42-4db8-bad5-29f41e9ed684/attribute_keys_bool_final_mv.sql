ATTACH MATERIALIZED VIEW _ UUID '716e1b20-b0e2-42e9-a421-188d58ef33f8' TO signoz_logs.logs_attribute_keys
(
    `name` String,
    `datatype` String
) AS
SELECT DISTINCT
    arrayJoin(attributes_bool_key) AS name,
    'Bool' AS datatype
FROM signoz_logs.logs
ORDER BY name ASC
