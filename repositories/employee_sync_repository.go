package repositories

import (
	"context"
	"static-api/dto"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncRepository struct {
	db *pgxpool.Pool
}

func NewSyncRepository(db *pgxpool.Pool) *SyncRepository {
	return &SyncRepository{db: db}
}

func (r *SyncRepository) Sync(ctx context.Context, req dto.SyncRequest) ([]dto.Employee, dto.Cursor, bool, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, dto.Cursor{}, false, err
	}
	defer tx.Rollback(ctx)

	var cursorTime *time.Time
	var cursorID *string

	if req.Cursor != nil {

		if req.Cursor.CursorTime != "" {
			t, err := time.Parse(time.RFC3339, req.Cursor.CursorTime)
			if err != nil {
				return nil, dto.Cursor{}, false, err
			}
			cursorTime = &t
		}

		if req.Cursor.CursorID != "" {
			cursorID = &req.Cursor.CursorID
		}
	}

	// =========================
	// 1. PUSH (Employees + Mobiles)
	// =========================

	for _, e := range req.Employees {

		var deletedAt *time.Time
		if e.DeletedAt != nil {
			now := time.Now().UTC()
			deletedAt = &now
		}

		// Employee upsert
		_, err := tx.Exec(ctx, `
		INSERT INTO employees (
			id, name, designation, department, is_active,
			img_url, email, city, country, joining_date,
			version, updated_at, deleted_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),$12)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			designation = EXCLUDED.designation,
			department = EXCLUDED.department,
			is_active = EXCLUDED.is_active,
			img_url = EXCLUDED.img_url,
			email = EXCLUDED.email,
			city = EXCLUDED.city,
			country = EXCLUDED.country,
			joining_date = EXCLUDED.joining_date,
			version = EXCLUDED.version,
			updated_at = NOW(),
			deleted_at = EXCLUDED.deleted_at
		WHERE employees.version < EXCLUDED.version;
		`,
			e.ID, e.Name, e.Designation, e.Department, e.IsActive,
			e.ImgURL, e.Email, e.City, e.Country, e.JoiningDate,
			e.Version, deletedAt,
		)

		if err != nil {
			return nil, dto.Cursor{}, false, err
		}

		// 🔥 Replace mobiles safely
		_, err = tx.Exec(ctx, `
			UPDATE mobiles
			SET deleted_at = NOW(), updated_at = NOW()
			WHERE employee_id = $1;
		`, e.ID)

		if err != nil {
			return nil, dto.Cursor{}, false, err
		}

		for _, m := range e.Mobiles {

			var mDeletedAt *time.Time
			if m.DeletedAt != nil {
				now := time.Now().UTC()
				mDeletedAt = &now
			}

			_, err := tx.Exec(ctx, `
			INSERT INTO mobiles (
				id, employee_id, type, number,
				version, updated_at, deleted_at
			)
			VALUES ($1,$2,$3,$4,$5,NOW(),$6)
			ON CONFLICT (id) DO UPDATE SET
				type = EXCLUDED.type,
				number = EXCLUDED.number,
				version = EXCLUDED.version,
				updated_at = NOW(),
				deleted_at = NULL
			WHERE mobiles.version < EXCLUDED.version;
			`,
				m.ID, e.ID, m.Type, m.Number,
				m.Version, mDeletedAt,
			)

			if err != nil {
				return nil, dto.Cursor{}, false, err
			}
		}
	}

	// =========================
	// 2. PULL (Employees + Mobiles)
	// =========================

	if req.Limit == 0 {
		req.Limit = 50
	}

	rows, err := tx.Query(ctx, `
	SELECT id, name, designation, department, is_active,
	       img_url, email, city, country, joining_date,
	       version, updated_at, deleted_at
	FROM employees
	WHERE 
		($1 IS NULL 
		OR updated_at > $1 
		OR (updated_at = $1 AND id > $2))
	ORDER BY updated_at, id
	LIMIT $3;
	`, cursorTime, cursorID, req.Limit+1)

	if err != nil {
		return nil, dto.Cursor{}, false, err
	}
	defer rows.Close()

	var employees []dto.Employee

	for rows.Next() {
		var e dto.Employee
		err := rows.Scan(
			&e.ID, &e.Name, &e.Designation, &e.Department,
			&e.IsActive, &e.ImgURL, &e.Email,
			&e.City, &e.Country, &e.JoiningDate,
			&e.Version, &e.UpdatedAt, &e.DeletedAt,
		)
		if err != nil {
			return nil, dto.Cursor{}, false, err
		}

		// fetch mobiles
		mRows, _ := tx.Query(ctx, `
			SELECT id, employee_id, type, number, version, updated_at, deleted_at
			FROM mobiles WHERE employee_id=$1;
		`, e.ID)

		for mRows.Next() {
			var m dto.Mobile
			mRows.Scan(&m.ID, &m.EmployeeID, &m.Type, &m.Number,
				&m.Version, &m.UpdatedAt, &m.DeletedAt)
			e.Mobiles = append(e.Mobiles, m)
		}

		employees = append(employees, e)
	}

	hasMore := len(employees) > req.Limit
	if hasMore {
		employees = employees[:req.Limit]
	}

	var nextCursor dto.Cursor
	if len(employees) > 0 {
		last := employees[len(employees)-1]
		nextCursor = dto.Cursor{
			CursorTime: last.UpdatedAt.Format(time.RFC3339),
			CursorID:   last.ID,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, dto.Cursor{}, false, err
	}

	return employees, nextCursor, hasMore, nil
}
