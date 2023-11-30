ATTACH MATERIALIZED VIEW _ UUID '356b8051-9acc-4fc9-b67f-a5f8217abdd7' TO signoz_logs.logs_resource_keys
(
    `name` String,
    `datatype` String
) AS
SELECT DISTINCT
    arrayJoin(resources_string_key) AS name,
    'String' AS datatype
FROM signoz_logs.logs
ORDER BY name ASC
