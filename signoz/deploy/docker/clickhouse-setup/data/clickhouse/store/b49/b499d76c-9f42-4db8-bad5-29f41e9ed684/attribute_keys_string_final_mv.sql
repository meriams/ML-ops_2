ATTACH MATERIALIZED VIEW _ UUID '30bd251f-5dd5-4963-8640-8596d3155be2' TO signoz_logs.logs_attribute_keys
(
    `name` String,
    `datatype` String
) AS
SELECT DISTINCT
    arrayJoin(attributes_string_key) AS name,
    'String' AS datatype
FROM signoz_logs.logs
ORDER BY name ASC
