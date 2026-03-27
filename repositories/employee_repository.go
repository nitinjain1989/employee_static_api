package repositories

import (
	"context"
	"errors"
	"fmt"
	"static-api/dto"
	"static-api/models"
	"static-api/utils"
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

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type EmployeeResponse struct {
	Data         []models.Employee
	ContentRange string
}

/*func (r *EmployeeRepository) GetEmployeeByID(
	ctx context.Context,
	id string,
) (*dto.EmployeeResponse, error) {

	rows, err := r.PgDB.Query(ctx, `
		SELECT
			e.id, e.name, e.email, e.designation, e.department, e.city,
			e.country, e.img_url, e.is_active, e.joining_date,
			e.version, e.updated_at, e.deleted_at,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE e.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employee *dto.EmployeeResponse

	for rows.Next() {
		var (
			e dto.EmployeeResponse
			m dto.MobileResponse

			mobileID   *string
			mobileNum  *string
			mobileType *string

			joiningDate *time.Time
			updatedAt   time.Time
			deletedAt   *time.Time
		)

		err := rows.Scan(
			&e.ID, &e.Name, &e.Email,
			&e.Designation, &e.Department, &e.City,
			&e.Country, &e.ImgURL, &e.IsActive,
			&joiningDate,
			&e.Version,
			&updatedAt,
			&deletedAt,
			&mobileID, &mobileNum, &mobileType,
		)
		if err != nil {
			return nil, err
		}

		// ✅ First row → initialize employee
		if employee == nil {
			e.Mobiles = []dto.MobileResponse{} // 🔥 ensures empty array

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			} else {
				e.JoiningDate = ""
			}

			// ✅ updated_at → ISO string
			e.UpdatedAt = updatedAt.Format(time.RFC3339)

			// ✅ deleted_at → "" if NULL
			if deletedAt != nil {
				e.DeletedAt = deletedAt.Format(time.RFC3339)
			} else {
				e.DeletedAt = ""
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
}*/

func (r *EmployeeRepository) FetchEmployees(
	ctx context.Context,
	filter dto.EmployeeFilterRequest,
) ([]dto.EmployeeResponse, dto.Meta, error) {

	limit, offset, paginate, err := validatePagination(filter)
	if err != nil {
		return nil, dto.Meta{}, err
	}

	countQuery, countArgs := buildCountQuery(filter)
	totalCount, err := r.executeCountQuery(ctx, countQuery, countArgs)
	if err != nil {
		return nil, dto.Meta{}, err
	}

	dataQuery, dataArgs := buildDataQuery(filter, paginate, limit, offset)
	rows, err := r.executeDataQuery(ctx, dataQuery, dataArgs)
	if err != nil {
		return nil, dto.Meta{}, err
	}

	defer rows.Close()

	employees, err := mapEmployees(rows)
	if err != nil {
		return nil, dto.Meta{}, err
	}

	meta := buildMeta(totalCount, paginate, limit, offset)
	return employees, meta, nil
}

func buildMeta(totalCount int, paginate bool, limit, offset int) dto.Meta {
	if paginate {
		page := (offset / limit) + 1
		return dto.Meta{
			TotalCount:  totalCount,
			Page:        page,
			PageSize:    limit,
			HasNextPage: offset+limit < totalCount,
		}
	}

	return dto.Meta{
		TotalCount:  totalCount,
		Page:        1,
		PageSize:    totalCount,
		HasNextPage: false,
	}
}

func mapEmployees(rows pgx.Rows) ([]dto.EmployeeResponse, error) {
	employeeMap := make(map[string]*dto.EmployeeResponse)
	var orderedIDs []string

	for rows.Next() {
		var (
			e dto.EmployeeResponse
			m dto.MobileResponse

			mobileID    *string
			mobileNum   *string
			mobileType  *string
			joiningDate *time.Time
			deletedAt   *time.Time
			updatedAt   time.Time
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
			&e.Version,
			&updatedAt,
			&deletedAt,
			&mobileID,
			&mobileNum,
			&mobileType,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := employeeMap[e.ID]; !exists {
			e.Mobiles = []dto.MobileResponse{}

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			}

			e.UpdatedAt = updatedAt.Format(time.RFC3339)

			if deletedAt != nil {
				e.DeletedAt = deletedAt.Format(time.RFC3339)
			} else {
				e.DeletedAt = ""
			}

			employeeMap[e.ID] = &e
			orderedIDs = append(orderedIDs, e.ID)
		}

		if mobileID != nil {
			m.ID = *mobileID
			m.Number = *mobileNum
			m.Type = *mobileType
			employeeMap[e.ID].Mobiles = append(employeeMap[e.ID].Mobiles, m)
		}
	}

	var employees []dto.EmployeeResponse
	for _, id := range orderedIDs {
		employees = append(employees, *employeeMap[id])
	}

	return employees, nil
}

func validatePagination(filter dto.EmployeeFilterRequest) (int, int, bool, error) {
	if (filter.Limit == nil && filter.Offset != nil) ||
		(filter.Limit != nil && filter.Offset == nil) {
		return 0, 0, false, errors.New("limit and offset must be provided together")
	}

	if filter.Limit != nil && filter.Offset != nil {
		return *filter.Limit, *filter.Offset, true, nil
	}

	return 0, 0, false, nil
}

func buildCountQuery(filter dto.EmployeeFilterRequest) (string, []interface{}) {
	query := `
		SELECT COUNT(DISTINCT e.id)
		FROM employees e
		WHERE 1=1
	`

	var args []interface{}
	idx := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND e.name ILIKE $%d", idx)
		args = append(args, "%"+filter.Search+"%")
		idx++
	}

	if len(filter.Designation) > 0 {
		query += fmt.Sprintf(" AND e.designation = ANY($%d)", idx)
		args = append(args, filter.Designation)
		idx++
	}

	if len(filter.Department) > 0 {
		query += fmt.Sprintf(" AND e.department = ANY($%d)", idx)
		args = append(args, filter.Department)
		idx++
	}

	if filter.Status == "active" {
		query += " AND e.is_active = true"
	} else if filter.Status == "inactive" {
		query += " AND e.is_active = false"
	}

	return query, args
}

func buildDataQuery(
	filter dto.EmployeeFilterRequest,
	paginate bool,
	limit int,
	offset int,
) (string, []interface{}) {

	query := `
		SELECT 
			e.id, e.name, e.email, e.designation, e.department, e.city,
			e.country, e.img_url, e.is_active, e.joining_date, e.version, e.updated_at, e.deleted_at,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE 1=1
	`

	var args []interface{}
	idx := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND e.name ILIKE $%d", idx)
		args = append(args, "%"+filter.Search+"%")
		idx++
	}

	if len(filter.Designation) > 0 {
		query += fmt.Sprintf(" AND e.designation = ANY($%d)", idx)
		args = append(args, filter.Designation)
		idx++
	}

	if len(filter.Department) > 0 {
		query += fmt.Sprintf(" AND e.department = ANY($%d)", idx)
		args = append(args, filter.Department)
		idx++
	}

	if filter.Status == "active" {
		query += " AND e.is_active = true"
	} else if filter.Status == "inactive" {
		query += " AND e.is_active = false"
	}

	// Sorting
	if filter.Search != "" {
		query += " ORDER BY e.name ASC, e.created_at DESC, e.id DESC"
	} else {
		query += " ORDER BY e.created_at DESC, e.id DESC"
	}

	// Pagination
	if paginate {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
		args = append(args, limit, offset)
	}

	return query, args
}

func (r *EmployeeRepository) executeCountQuery(
	ctx context.Context,
	query string,
	args []interface{},
) (int, error) {

	var total int
	err := r.PgDB.QueryRow(ctx, query, args...).Scan(&total)
	return total, err
}

func (r *EmployeeRepository) executeDataQuery(
	ctx context.Context,
	query string,
	args []interface{},
) (pgx.Rows, error) {
	return r.PgDB.Query(ctx, query, args...)
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
	emp dto.CreateEmployeeRequest,
) (*dto.EmployeeResponse, error) {

	var result *dto.EmployeeResponse

	err := r.withTx(ctx, func(tx pgx.Tx) error {

		var empID string

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

		res, err := r.getEmployeeByIDTx(ctx, tx, empID)
		if err != nil {
			return err
		}

		result = res

		return nil
	})

	return result, err
}

// Update Emoplyee
func (r *EmployeeRepository) UpdateEmployeeWithMobiles(
	ctx context.Context,
	id string,
	emp dto.UpdateEmployeeRequest,
) (*dto.EmployeeResponse, error) {

	var result *dto.EmployeeResponse

	err := r.withTx(ctx, func(tx pgx.Tx) error {

		// ✅ Convert JoiningDate (string → *time.Time)
		var joiningDate *time.Time
		if emp.JoiningDate != "" {
			t, err := time.Parse("2006-01-02", emp.JoiningDate)
			if err != nil {
				return err
			}
			joiningDate = &t
		}

		// ✅ 1. Update employee with optimistic locking
		cmd, err := tx.Exec(ctx, `
			UPDATE employees SET
				name = $1,
				designation = $2,
				department = $3,
				is_active = $4,
				img_url = $5,
				email = $6,
				city = $7,
				country = $8,
				joining_date = $9,
				version = version + 1,
				updated_at = NOW()
			WHERE id = $10 AND deleted_at IS NULL
		`,
			emp.Name,
			emp.Designation,
			emp.Department,
			emp.IsActive,
			emp.ImgURL,
			emp.Email,
			emp.City,
			emp.Country,
			joiningDate,
			id,
		)
		if err != nil {
			return err
		}

		// ❌ Conflict or not found
		if cmd.RowsAffected() == 0 {
			return utils.ErrConflict
		}

		// ✅ 2. Replace mobiles (simple approach)
		_, err = tx.Exec(ctx,
			`DELETE FROM mobiles WHERE employee_id = $1`,
			id,
		)
		if err != nil {
			return err
		}

		for _, m := range emp.Mobiles {
			_, err := tx.Exec(ctx, `
				INSERT INTO mobiles (employee_id, number, type)
				VALUES ($1, $2, $3)
			`,
				id, m.Number, m.Type,
			)
			if err != nil {
				return err
			}
		}

		// ✅ 3. Fetch updated employee INSIDE TX
		updatedEmp, err := r.getEmployeeByIDTx(ctx, tx, id)
		if err != nil {
			return err
		}

		result = updatedEmp
		return nil
	})

	return result, err
}

// Delete Emoplyee
func (r *EmployeeRepository) DeleteEmployee(ctx context.Context, id string) error {

	tx, err := r.PgDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	cmd, err := tx.Exec(ctx, `
	UPDATE employees
	SET 
		deleted_at = NOW(),
		version = version + 1,
		updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
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

func (r *EmployeeRepository) GetEmployeeByID(
	ctx context.Context,
	id string,
) (*dto.EmployeeResponse, error) {
	return r.getEmployeeByIDInternal(ctx, r.PgDB, id)
}

func (r *EmployeeRepository) getEmployeeByIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*dto.EmployeeResponse, error) {
	return r.getEmployeeByIDInternal(ctx, tx, id)
}

func (r *EmployeeRepository) getEmployeeByIDInternal(
	ctx context.Context,
	q Querier,
	id string,
) (*dto.EmployeeResponse, error) {

	rows, err := q.Query(ctx, `
		SELECT 
			e.id, e.name, e.email, e.designation, e.department, e.city,
			e.country, e.img_url, e.is_active, e.joining_date,
			e.version, e.updated_at, e.deleted_at,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE e.id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employee *dto.EmployeeResponse

	for rows.Next() {
		var (
			e dto.EmployeeResponse
			m dto.MobileResponse

			mobileID   *string
			mobileNum  *string
			mobileType *string

			joiningDate *time.Time
			updatedAt   time.Time
			deletedAt   *time.Time
		)

		err := rows.Scan(
			&e.ID, &e.Name, &e.Email,
			&e.Designation, &e.Department, &e.City,
			&e.Country, &e.ImgURL, &e.IsActive,
			&joiningDate,
			&e.Version,
			&updatedAt,
			&deletedAt,
			&mobileID, &mobileNum, &mobileType,
		)
		if err != nil {
			return nil, err
		}

		if employee == nil {
			e.Mobiles = []dto.MobileResponse{}

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			}

			e.UpdatedAt = updatedAt.Format(time.RFC3339)

			if deletedAt != nil {
				e.DeletedAt = deletedAt.Format(time.RFC3339)
			} else {
				e.DeletedAt = ""
			}

			employee = &e
		}

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

	if employee == nil {
		return nil, fmt.Errorf("employee not found")
	}

	return employee, nil
}
