-- Rollback: drop tables in correct order (child -> parent)

DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;
