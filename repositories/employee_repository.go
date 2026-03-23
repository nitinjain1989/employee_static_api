package repositories

import (
	"context"
	"errors"
	"fmt"
	"static-api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmployeeRepository struct {
	PgDB *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{
		PgDB: pool,
	}
}

type EmployeeResponse struct {
	Data         []models.Employee
	ContentRange string
}

func (r *EmployeeRepository) GetEmployeeByID(
	ctx context.Context,
	id string,
) (*models.Employee, error) {

	rows, err := r.PgDB.Query(ctx, `
		SELECT 
			e.id, e.name, e.email, e.designation, e.department, e.city,
			e.country, e.img_url, e.is_active, e.joining_date,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE e.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employee *models.Employee

	for rows.Next() {
		var (
			e models.Employee
			m models.Mobile

			mobileID   *string
			mobileNum  *string
			mobileType *string

			joiningDate *time.Time
		)

		err := rows.Scan(
			&e.ID, &e.Name, &e.Email,
			&e.Designation, &e.Department, &e.City,
			&e.Country, &e.ImgURL, &e.IsActive,
			&joiningDate,
			&mobileID, &mobileNum, &mobileType,
		)
		if err != nil {
			return nil, err
		}

		// ✅ First row → initialize employee
		if employee == nil {
			e.Mobiles = []models.Mobile{} // 🔥 ensures empty array

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			}

			employee = &e
		}

		// ✅ Only add mobile if exists
		if mobileID != nil {
			m.ID = *mobileID

			if mobileNum != nil {
				m.Number = *mobileNum
			}
			if mobileType != nil {
				m.Type = *mobileType
			}

			employee.Mobiles = append(employee.Mobiles, m)
		}
	}

	// ❌ Employee not found
	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	// ✅ If no mobiles → already []
	return employee, nil
}

func (r *EmployeeRepository) FetchEmployees(
	ctx context.Context,
	filter models.EmployeeFilter,
) ([]models.Employee, models.Meta, error) {

	// =========================
	// 🔍 PAGINATION CHECK
	// =========================
	var (
		limit    int
		offset   int
		paginate bool
	)

	if filter.Limit != nil && filter.Offset != nil {
		limit = *filter.Limit
		offset = *filter.Offset
		paginate = true
	}

	// ❗ validation
	if (filter.Limit == nil && filter.Offset != nil) ||
		(filter.Limit != nil && filter.Offset == nil) {
		return nil, models.Meta{}, errors.New("limit and offset must be provided together")
	}

	// =========================
	// 🔢 COUNT QUERY
	// =========================
	countQuery := `
		SELECT COUNT(DISTINCT e.id)
		FROM employees e
		WHERE 1=1
	`

	var countArgs []interface{}
	countIdx := 1

	if filter.Search != "" {
		countQuery += fmt.Sprintf(" AND e.name ILIKE $%d", countIdx)
		countArgs = append(countArgs, "%"+filter.Search+"%")
		countIdx++
	}

	if len(filter.Designation) > 0 {
		countQuery += fmt.Sprintf(" AND e.designation = ANY($%d)", countIdx)
		countArgs = append(countArgs, filter.Designation)
		countIdx++
	}

	if len(filter.Department) > 0 {
		countQuery += fmt.Sprintf(" AND e.department = ANY($%d)", countIdx)
		countArgs = append(countArgs, filter.Department)
		countIdx++
	}

	if filter.Status == "active" {
		countQuery += " AND e.is_active = true"
	} else if filter.Status == "inactive" {
		countQuery += " AND e.is_active = false"
	}

	var totalCount int
	if err := r.PgDB.QueryRow(ctx, countQuery, countArgs...).Scan(&totalCount); err != nil {
		return nil, models.Meta{}, err
	}

	// =========================
	// 📦 DATA QUERY
	// =========================
	query := `
		SELECT 
			e.id, e.name, e.email, e.designation, e.department, e.city,
			e.country, e.img_url, e.is_active, e.joining_date,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE 1=1
	`

	var args []interface{}
	idx := 1

	// 🔍 Search
	if filter.Search != "" {
		query += fmt.Sprintf(" AND e.name ILIKE $%d", idx)
		args = append(args, "%"+filter.Search+"%")
		idx++
	}

	// 🎯 designation
	if len(filter.Designation) > 0 {
		query += fmt.Sprintf(" AND e.designation = ANY($%d)", idx)
		args = append(args, filter.Designation)
		idx++
	}

	// 🎯 department
	if len(filter.Department) > 0 {
		query += fmt.Sprintf(" AND e.department = ANY($%d)", idx)
		args = append(args, filter.Department)
		idx++
	}

	// 🔄 status
	if filter.Status == "active" {
		query += " AND e.is_active = true"
	} else if filter.Status == "inactive" {
		query += " AND e.is_active = false"
	}

	// 📊 Sorting
	if filter.Search != "" {
		query += " ORDER BY e.name ASC, e.created_at DESC, e.id DESC"
	} else {
		query += " ORDER BY e.created_at DESC, e.id DESC"
	}

	// 📦 Pagination (optional)
	if paginate {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
		args = append(args, limit, offset)
	}

	rows, err := r.PgDB.Query(ctx, query, args...)
	if err != nil {
		return nil, models.Meta{}, err
	}
	defer rows.Close()

	// =========================
	// 🧠 MAP DATA
	// =========================
	employeeMap := make(map[string]*models.Employee)
	var orderedIDs []string
	for rows.Next() {
		var (
			e models.Employee
			m models.Mobile

			mobileID    *string
			mobileNum   *string
			mobileType  *string
			joiningDate *time.Time
		)

		err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.Email,
			&e.Designation,
			&e.Department,
			&e.City,
			&e.Country,
			&e.ImgURL,
			&e.IsActive,
			&joiningDate,
			&mobileID,
			&mobileNum,
			&mobileType,
		)
		if err != nil {
			return nil, models.Meta{}, err
		}

		if _, exists := employeeMap[e.ID]; !exists {
			e.Mobiles = []models.Mobile{}

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			}

			employeeMap[e.ID] = &e
			orderedIDs = append(orderedIDs, e.ID)
		}

		// ✅ Only append if mobile exists
		if mobileID != nil {
			m.ID = *mobileID
			m.Number = *mobileNum
			m.Type = *mobileType

			employeeMap[e.ID].Mobiles = append(employeeMap[e.ID].Mobiles, m)
		}
	}

	var employees []models.Employee
	for _, id := range orderedIDs {
		employees = append(employees, *employeeMap[id])
	}

	// =========================
	// 📊 META
	// =========================
	var meta models.Meta

	if paginate {
		page := (offset / limit) + 1
		hasNextPage := offset+limit < totalCount

		meta = models.Meta{
			TotalCount:  totalCount,
			Page:        page,
			PageSize:    limit,
			HasNextPage: hasNextPage,
		}
	} else {
		meta = models.Meta{
			TotalCount:  totalCount,
			Page:        1,
			PageSize:    totalCount,
			HasNextPage: false,
		}
	}

	return employees, meta, nil
}

func (r *EmployeeRepository) FetchDistinctField(
	ctx context.Context,
	field string,
) ([]string, error) {

	// ⚠️ Prevent SQL injection (VERY IMPORTANT)
	allowedFields := map[string]bool{
		"designation": true,
		"department":  true,
	}

	if !allowedFields[field] {
		return nil, fmt.Errorf("invalid field")
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT %s
		FROM employees
		WHERE %s IS NOT NULL AND %s != ''
		ORDER BY %s ASC
	`, field, field, field, field)

	rows, err := r.PgDB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string

	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err != nil {
			return nil, err
		}
		values = append(values, val)
	}

	return values, nil
}

// Create Emoplyee
func (r *EmployeeRepository) CreateEmployeeWithMobiles(
	ctx context.Context,
	emp models.Employee,
) (string, error) {

	var empID string

	err := r.withTx(ctx, func(tx pgx.Tx) error {

		// ✅ Convert JoiningDate (string → time.Time)
		var joiningDate *time.Time
		if emp.JoiningDate != "" {
			t, err := time.Parse("2006-01-02", emp.JoiningDate)
			if err != nil {
				return err
			}
			joiningDate = &t
		}

		// ✅ Insert full employee
		err := tx.QueryRow(ctx, `
			INSERT INTO employees (
				name, email, designation, department,
				city, country, img_url, is_active, joining_date
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			RETURNING id
		`,
			emp.Name,
			emp.Email,
			emp.Designation,
			emp.Department,
			emp.City,
			emp.Country,
			emp.ImgURL,
			emp.IsActive,
			joiningDate, // ✅ proper type
		).Scan(&empID)

		if err != nil {
			return err
		}

		// ✅ Insert mobiles
		for _, m := range emp.Mobiles {
			_, err := tx.Exec(ctx, `
				INSERT INTO mobiles (employee_id, number, type)
				VALUES ($1, $2, $3)
			`,
				empID,
				m.Number,
				m.Type,
			)

			if err != nil {
				return err
			}
		}

		return nil
	})

	return empID, err
}

// Update Emoplyee
func (r *EmployeeRepository) UpdateEmployeeWithMobiles(
	ctx context.Context,
	id string,
	emp models.Employee,
) error {

	tx, err := r.PgDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // 🔥 rollback if anything fails

	// ✅ 1. Update employee
	_, err = tx.Exec(ctx, `
		UPDATE employees SET
			name = $1,
			designation = $2,
			department = $3,
			is_active = $4,
			img_url = $5,
			email = $6,
			city = $7,
			country = $8,
			joining_date = $9
		WHERE id = $10
	`,
		emp.Name,
		emp.Designation,
		emp.Department,
		emp.IsActive,
		emp.ImgURL,
		emp.Email,
		emp.City,
		emp.Country,
		emp.JoiningDate,
		id,
	)
	if err != nil {
		return err
	}

	// ✅ 2. Delete old mobiles
	_, err = tx.Exec(ctx,
		`DELETE FROM mobiles WHERE employee_id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	// ✅ 3. Insert new mobiles
	for _, m := range emp.Mobiles {
		_, err := tx.Exec(ctx, `
			INSERT INTO mobiles (employee_id, number, type)
			VALUES ($1, $2, $3)
		`,
			id, m.Number, m.Type,
		)

		if err != nil {
			return err // 🔥 rollback triggered
		}
	}

	// ✅ 4. Commit only if everything succeeds
	return tx.Commit(ctx)
}

// Delete Emoplyee
func (r *EmployeeRepository) DeleteEmployee(ctx context.Context, id string) error {

	tx, err := r.PgDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := tx.Exec(ctx,
		`DELETE FROM employees WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("employee not found")
	}

	return tx.Commit(ctx)
}

func (r *EmployeeRepository) withTx(
	ctx context.Context,
	fn func(tx pgx.Tx) error,
) error {
	tx, err := r.PgDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
