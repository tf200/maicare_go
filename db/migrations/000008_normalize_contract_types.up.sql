WITH ranked_types AS (
    SELECT
        id,
        FIRST_VALUE(id) OVER (
            PARTITION BY LOWER(BTRIM(name))
            ORDER BY id
        ) AS canonical_id
    FROM contract_type
), duplicate_types AS (
    SELECT id, canonical_id
    FROM ranked_types
    WHERE id <> canonical_id
)
UPDATE contract
SET type_id = duplicate_types.canonical_id
FROM duplicate_types
WHERE contract.type_id = duplicate_types.id;

WITH ranked_types AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY LOWER(BTRIM(name))
            ORDER BY id
        ) AS row_number
    FROM contract_type
)
DELETE FROM contract_type
USING ranked_types
WHERE contract_type.id = ranked_types.id
  AND ranked_types.row_number > 1;

UPDATE contract_type
SET name = BTRIM(name);

CREATE UNIQUE INDEX contract_type_normalized_name_idx
ON contract_type (LOWER(name));
