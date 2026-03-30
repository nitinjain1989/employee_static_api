package repositories

import (
	"context"
	"fmt"
	"static-api/dto"
	"static-api/utils"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SyncRepository struct {
	db *pgxpool.Pool
}

func NewSyncRepository(db *pgxpool.Pool) *SyncRepository {
	return &SyncRepository{db: db}
}

func (r *SyncRepository) Sync(
	ctx context.Context,
	req dto.SyncRequest,
) ([]dto.EmployeeResponse, dto.Cursor, bool, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, dto.Cursor{}, false, err
	}
	defer tx.Rollback(ctx)

	// =========================
	// 🔐 CURSOR HANDLING (SAFE)
	// =========================

	var cursorSeq int64
	if req.Cursor != nil {
		cursorSeq = req.Cursor.Seq
	}

	// =========================
	// 1. PUSH (UPSERT)
	// =========================
	for _, e := range req.Employees {

		if e.ID == "" {
			return nil, dto.Cursor{}, false, fmt.Errorf("employee id cannot be empty")
		}

		if !utils.IsValidUUID(e.ID) {
			return nil, dto.Cursor{}, false, fmt.Errorf("invalid employee id: %s", e.ID)
		}

		var deletedAt *time.Time
		if e.DeletedAt != nil {
			now := time.Now().UTC()
			deletedAt = &now
		}

		// Employee UPSERT (last write wins)
		_, err := tx.Exec(ctx, `
			INSERT INTO employees (
				id, name, designation, department, is_active,
				img_url, email, city, country, joining_date,
				version, deleted_at
			)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,1,$11)

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
				version = employees.version + 1,
				deleted_at = EXCLUDED.deleted_at;
		`,
			e.ID, e.Name, e.Designation, e.Department, e.IsActive,
			e.ImgURL, e.Email, e.City, e.Country, e.JoiningDate,
			deletedAt,
		)
		if err != nil {
			return nil, dto.Cursor{}, false, err
		}

		// 🔥 HARD delete old mobiles
		_, err = tx.Exec(ctx, `
			DELETE FROM mobiles
			WHERE employee_id = $1;
		`, e.ID)
		if err != nil {
			return nil, dto.Cursor{}, false, err
		}

		// Insert latest mobiles (fresh state)
		for _, m := range e.Mobiles {

			_, err := tx.Exec(ctx, `
				INSERT INTO mobiles (
					employee_id, type, number
				)
				VALUES ($1,$2,$3);
			`,
				e.ID, m.Type, m.Number,
			)
			if err != nil {
				return nil, dto.Cursor{}, false, err
			}
		}
	}

	// =========================
	// 2. PULL (CORRECT PAGINATION)
	// =========================

	if req.Limit <= 0 {
		req.Limit = 50
	}

	rows, err := tx.Query(ctx, `
		SELECT 
			e.id, e.name, e.email, e.designation, e.department,
			e.city, e.country, e.img_url, e.is_active, e.joining_date,e.created_at,
			e.version, e.updated_at, e.deleted_at, e.updated_seq,
			m.id, m.number, m.type
		FROM employees e
		LEFT JOIN mobiles m ON m.employee_id = e.id
		WHERE e.updated_seq > $1
		ORDER BY e.updated_seq ASC
		LIMIT $2;
	`, cursorSeq, req.Limit+1)

	if err != nil {
		return nil, dto.Cursor{}, false, err
	}
	defer rows.Close()

	employees, lastSeq, count, err := mapEmployeesWithSeq(rows)
	if err != nil {
		return nil, dto.Cursor{}, false, err
	}

	if count == 0 {
		if err := tx.Commit(ctx); err != nil {
			return nil, dto.Cursor{}, false, err
		}

		return []dto.EmployeeResponse{}, dto.Cursor{
			Seq: cursorSeq, // 🔥 FIX
		}, false, nil
	}

	hasMore := count > req.Limit

	nextCursor := dto.Cursor{
		Seq: lastSeq,
	}

	// =========================
	// COMMIT
	// =========================
	if err := tx.Commit(ctx); err != nil {
		return nil, dto.Cursor{}, false, err
	}

	return employees, nextCursor, hasMore, nil
}

func mapEmployeesWithSeq(rows pgx.Rows) ([]dto.EmployeeResponse, int64, int, error) {

	empMap := make(map[string]*dto.EmployeeResponse)
	var lastSeq int64
	count := 0

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
			createdAt   time.Time
			updatedSeq  int64
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
			&createdAt,
			&e.Version,
			&updatedAt,
			&deletedAt,
			&updatedSeq,
			&mobileID,
			&mobileNum,
			&mobileType,
		)

		if err != nil {
			return nil, 0, 0, err
		}

		lastSeq = updatedSeq
		count++

		_, exists := empMap[e.ID]
		if !exists {

			e.Mobiles = []dto.MobileResponse{}

			if joiningDate != nil {
				e.JoiningDate = joiningDate.Format("2006-01-02")
			}

			e.UpdatedAt = updatedAt.Format(time.RFC3339)
			e.CreatedAt = createdAt.Format(time.RFC3339)
			if deletedAt != nil {
				e.DeletedAt = deletedAt.Format(time.RFC3339)
			} else {
				e.DeletedAt = ""
			}

			empMap[e.ID] = &e
		}

		if mobileID != nil {
			m.ID = *mobileID
			m.Number = *mobileNum
			m.Type = *mobileType
			empMap[e.ID].Mobiles = append(empMap[e.ID].Mobiles, m)
		}
	}

	// convert map → slice
	result := make([]dto.EmployeeResponse, 0, len(empMap))
	for _, v := range empMap {
		result = append(result, *v)
	}

	return result, lastSeq, count, nil
}
