SELECT
    parent.service_name,
    child.service_name,
    COUNT(*)
FROM spans AS child
         JOIN spans AS parent ON child.parent_id = parent.span_id
WHERE child.service_name != parent.service_name
GROUP BY parent.service_name, child.service_name