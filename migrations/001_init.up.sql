-- =========================
-- Departments
-- =========================
CREATE TABLE IF NOT EXISTS departments (
  id          BIGSERIAL PRIMARY KEY,
  name        TEXT NOT NULL UNIQUE,
  created_at  TIMESTAMP NOT NULL DEFAULT now(),
  updated_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- =========================
-- Employees
-- =========================
CREATE TABLE IF NOT EXISTS employees (
  id             BIGSERIAL PRIMARY KEY,
  name           TEXT NOT NULL,
  department_id  BIGINT NOT NULL,
  email          TEXT NOT NULL UNIQUE,
  age            INT,
  position       TEXT,
  salary         NUMERIC(12,2),
  created_at     TIMESTAMP NOT NULL DEFAULT now(),
  updated_at     TIMESTAMP NOT NULL DEFAULT now(),

  CONSTRAINT fk_department
    FOREIGN KEY (department_id)
    REFERENCES departments(id)
    ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_employees_department_id
ON employees(department_id);
