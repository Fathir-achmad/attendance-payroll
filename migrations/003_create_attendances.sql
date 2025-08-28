CREATE TABLE IF NOT EXISTS attendances (
        id SERIAL PRIMARY KEY,
        employee_id INT REFERENCES employees(id) ON DELETE CASCADE,
        date DATE NOT NULL,
        check_in_at TIMESTAMP NULL,
        check_out_at TIMESTAMP NULL
    );
