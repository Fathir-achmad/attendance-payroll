package models

import "time"

type Department struct {
    ID        uint       `gorm:"primaryKey" json:"id"`
    Name      string     `gorm:"size:100;not null" json:"name"`
    Employees []Employee `json:"employees,omitempty"`
}

type Employee struct {
    ID           uint       `gorm:"primaryKey" json:"id"`
    Name         string     `gorm:"size:100;not null" json:"name"`
    Username     string     `gorm:"unique;not null" json:"username"`
    PasswordHash string     `gorm:"not null" json:"-"`
    DepartmentID uint       `json:"department_id"`
    Department   Department `json:"department"`
    Attendances  []Attendance `json:"attendances,omitempty"`
    Payrolls     []Payroll    `json:"payrolls,omitempty"`
}

type Attendance struct {
    ID         uint       `gorm:"primaryKey" json:"id"`
    EmployeeID uint       `gorm:"not null" json:"employee_id"`
    Employee   Employee   `json:"employee"`
    Date       time.Time  `gorm:"not null" json:"date"`
    CheckInAt  *time.Time `json:"check_in_at"`
    CheckOutAt *time.Time `json:"check_out_at"`
}

type Payroll struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    EmployeeID  uint      `gorm:"not null" json:"employee_id"`
    Employee    Employee  `json:"employee"`
    Amount      float64   `gorm:"not null" json:"amount"`
    PeriodStart time.Time `json:"period_start"`
    PeriodEnd   time.Time `json:"period_end"`
}
