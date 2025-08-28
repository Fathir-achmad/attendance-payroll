CREATE TABLE IF NOT EXISTS payrolls (
        id SERIAL PRIMARY KEY,
        employee_id INT REFERENCES employees(id) ON DELETE CASCADE,
        amount NUMERIC(12,2) NOT NULL,
        period_start DATE NOT NULL,
        period_end DATE NOT NULL
    );
