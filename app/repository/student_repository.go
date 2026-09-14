package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"api-students/app/model"
)

var (
	ErrNotFound  = errors.New("data tidak ditemukan")
	ErrDuplicate = errors.New("nim sudah terdaftar")
)

type StudentRepository interface {
	FindAll(ctx context.Context, q model.ListQuery) ([]model.Student, int, error)
	FindByID(ctx context.Context, id int) (model.Student, error)
	Create(ctx context.Context, s model.Student) (model.Student, error)
	Update(ctx context.Context, s model.Student) (model.Student, error)
	Delete(ctx context.Context, id int) error
}

var kolomUrut = map[string]string{
	"id":         "id",
	"nim":        "nim",
	"name":       "name",
	"grade":      "grade",
	"created_at": "created_at",
}

type studentPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewStudentRepository(pool *pgxpool.Pool) StudentRepository {
	return &studentPostgresRepository{pool: pool}
}

func (r *studentPostgresRepository) FindAll(
	ctx context.Context, q model.ListQuery,
) ([]model.Student, int, error) {
	kolom := "id"
	if k, ok := kolomUrut[q.Sort]; ok {
		kolom = k
	}

	arah := "ASC"
	if strings.ToUpper(q.Order) == "DESC" {
		arah = "DESC"
	}

	offset := (q.Page - 1) * q.Limit

	var conditions []string
	var args []any
	argIdx := 1

	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR nim ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q.Search+"%")
		argIdx++
	}

	if q.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIdx))
		args = append(args, *q.IsActive)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM students %s", whereClause)
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("menghitung total data: %w", err)
	}

	posisiLimit := len(args) + 1
	posisiOffset := len(args) + 2
	dataQuery := fmt.Sprintf(
		"SELECT id, nim, name, grade, is_active, COALESCE(owner_id, 0), created_at, updated_at FROM students %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
		whereClause, kolom, arah, posisiLimit, posisiOffset,
	)

	args = append(args, q.Limit, offset)

	rows, err := r.pool.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("mengambil data: %w", err)
	}
	defer rows.Close()

	var students []model.Student
	for rows.Next() {
		var s model.Student
		if err := rows.Scan(&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.OwnerID, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("membaca baris data: %w", err)
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterasi data: %w", err)
	}

	return students, total, nil
}

func (r *studentPostgresRepository) FindByID(
	ctx context.Context, id int,
) (model.Student, error) {
	query := "SELECT id, nim, name, grade, is_active, COALESCE(owner_id, 0), created_at, updated_at FROM students WHERE id = $1"

	var s model.Student
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.OwnerID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		return model.Student{}, fmt.Errorf("mengambil student by id: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Create(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	query := `
		INSERT INTO students (nim, name, grade, is_active, owner_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, nim, name, grade, is_active, COALESCE(owner_id, 0), created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive, s.OwnerID).Scan(
		&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.OwnerID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("menyimpan student: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Update(
	ctx context.Context, s model.Student,
) (model.Student, error) {
	query := `
		UPDATE students
		SET nim = $1, name = $2, grade = $3, is_active = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, nim, name, grade, is_active, COALESCE(owner_id, 0), created_at, updated_at
	`

	err := r.pool.QueryRow(ctx, query, s.NIM, s.Name, s.Grade, s.IsActive, s.ID).Scan(
		&s.ID, &s.NIM, &s.Name, &s.Grade, &s.IsActive, &s.OwnerID, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Student{}, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Student{}, ErrDuplicate
		}
		return model.Student{}, fmt.Errorf("memperbarui student: %w", err)
	}

	return s, nil
}

func (r *studentPostgresRepository) Delete(
	ctx context.Context, id int,
) error {
	cmdTag, err := r.pool.Exec(ctx, "DELETE FROM students WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("menghapus student: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}
