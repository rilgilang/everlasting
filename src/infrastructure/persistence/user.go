package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/domain/user"
	"everlasting/src/infrastructure/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type UserPersistence struct {
	db     *sqlx.DB
	logger *logger.AppLogger
}

func NewUserPersistence(db *sqlx.DB, logger *logger.AppLogger) *UserPersistence {
	return &UserPersistence{
		db:     db,
		logger: logger,
	}
}

// GetOneByEmail implements user.UserRepository.
func (r *UserPersistence) GetOneByEmail(ctx context.Context, email user.Email) (result *user.User, err error) {
	query := `SELECT 
		id, email, name, role, status, 
		created_at, updated_at, ciphertext 
	FROM 
		users 
	WHERE 
		email = $1`

	stmt, err := r.db.Preparex(query)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_get_one_by_email", err.Error())
		return result, err
	}

	defer stmt.Close()
	if err != nil {
		r.logger.Error(ctx, "postgres:user_get_one_by_email", err.Error())
		return result, err
	}
	defer stmt.Close()

	var rs user.User
	err = stmt.QueryRow(email).Scan(
		&rs.ID, &rs.Email, &rs.Name, &rs.Role, &rs.Status,
		&rs.CreatedAt, &rs.UpdatedAt, &rs.CipherText,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrUserNotFound
		}
		r.logger.Error(ctx, "postgres:user_get_one_by_email", err.Error())
		return nil, err
	}
	return &rs, err
}

func (r *UserPersistence) GetOneByID(ctx context.Context, id identity.ID) (result *user.User, err error) {
	query := `SELECT 
		id, email, name, role, status, 
		created_at, updated_at, ciphertext 
	FROM 
		users 
	WHERE 
		id = $1`

	stmt, err := r.db.Preparex(query)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_get_one_by_id", err.Error())
		return result, err
	}

	defer stmt.Close()
	if err != nil {
		r.logger.Error(ctx, "postgres:user_get_one_by_id", err.Error())
		return result, err
	}
	defer stmt.Close()

	var rs user.User
	err = stmt.QueryRow(id).Scan(
		&rs.ID, &rs.Email, &rs.Name, &rs.Role, &rs.Status,
		&rs.CreatedAt, &rs.UpdatedAt, &rs.CipherText,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrUserNotFound
		}
		r.logger.Error(ctx, "postgres:user_get_one_by_id", err.Error())
		return nil, err
	}
	return &rs, err
}

func (r *UserPersistence) Create(ctx context.Context, newUser *user.User) (result *user.User, err error) {
	query := "INSERT INTO users (id, name, email, role, ciphertext ,status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"

	stmt, err := r.db.Preparex(query)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_create", err.Error())
		return result, err
	}

	defer stmt.Close()

	vals := []interface{}{
		newUser.ID,
		newUser.Name,
		newUser.Email,
		newUser.Role,
		newUser.CipherText,
		newUser.Status,
		time.Time(newUser.CreatedAt),
		time.Time(newUser.UpdatedAt),
	}
	var id identity.ID
	err = stmt.QueryRow(vals...).Scan(&id)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_create", err.Error())
		return result, err
	}

	return r.GetOneByID(ctx, id)
}

// UpdateByID implements user.UserRepository.
func (r *UserPersistence) UpdateByID(ctx context.Context, user *user.User, id identity.ID) (result *user.User, err error) {
	query := `UPDATE 
		users 
	SET 
		name=$2, email=$3, ciphertext=$4, status=$5, 
		updated_at=$6 
	WHERE 
		id=$1 RETURNING id`

	vals := []interface{}{
		user.ID,
		user.Name,
		user.Email,
		user.CipherText,
		user.Status,
		time.Time(user.UpdatedAt),
	}

	stmt, err := r.db.Preparex(query)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_update", err.Error())
		return result, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(vals...)
	if err != nil {
		r.logger.Error(ctx, "postgres:user_update", err.Error())
		return result, err
	}
	return r.GetOneByID(ctx, user.ID)
}

func (u *UserPersistence) GetByQuery(ctx context.Context, query *user.Query) (result *user.Users, err error) {
	// Generate filter query
	filter := u.generateFilterQuery(query)
	// Get data count (for pagination purpose)
	total, err := u.getTotalCount(ctx, filter)
	if err != nil {
		u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
		return result, err
	}

	// Generate data collection query
	q := sq.Select("u.id", "u.name", "u.email", "u.role", "u.status",
		"u.created_at", "u.updated_at", "u.ciphertext",
	).
		From("users u").
		Where(filter)

	// Apply limit and offset
	q, page, perPage := u.applyLimitAndOffset(q, query)

	// Sort by
	q = u.applySorting(q, query)

	// Return as $1, $2 format
	q = q.PlaceholderFormat(sq.Dollar)

	sql, args, err := q.ToSql()

	if err != nil {
		u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
		return result, err
	}

	stmt, err := u.generateStatement(ctx, sql)
	if err != nil {
		u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
		return result, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
		return nil, err
	}

	defer rows.Close()

	// Collect result
	collection := make([]user.User, 0)
	for rows.Next() {
		var data user.User
		if err := rows.Scan(&data.ID, &data.Name, &data.Email, &data.Role, &data.Status,
			&data.CreatedAt, &data.UpdatedAt, &data.CipherText); err != nil {
			u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
			return nil, err
		}
		collection = append(collection, data)
	}

	return &user.Users{
		Collection: collection,
		Pagination: user.Pagination{
			CurrentPage: page,
			MaxPage:     int64(math.Ceil(float64(total) / float64(perPage))),
			TotalData:   int64(total),
		},
	}, err
}

func (u *UserPersistence) generateFilterQuery(query *user.Query) (result sq.And) {
	if query.Q != "" {
		result = append(result,
			sq.ILike{"u.name": fmt.Sprintf("%%%s%%", query.Q)},
		)
	}

	if query.Role != "" {
		result = append(result,
			sq.Eq{"u.role": query.Role},
		)
	}

	if query.Status != "" {
		result = append(result,
			sq.Eq{"u.status": query.Status},
		)
	}

	return result
}

func (u *UserPersistence) getTotalCount(ctx context.Context, filter sq.And) (result uint, err error) {
	type CountChannel struct {
		total uint
		err   error
	}

	countCh := make(chan *CountChannel)
	go func(filter sq.And, m chan *CountChannel) {
		defer func() {
			if err := recover(); err != nil {
				m <- &CountChannel{
					0, fmt.Errorf("recovered: %v", err),
				}
			}
		}()
		q := sq.Select("count(u.id)").
			From("users u").
			Where(filter).
			PlaceholderFormat(sq.Dollar)

		sql, args, err := q.ToSql()
		if err != nil {
			m <- &CountChannel{
				0, err,
			}
		}

		stmt, err := u.generateStatement(ctx, sql)
		if err != nil {
			u.logger.Error(ctx, "postgres:user_get_by_query", err.Error())
			return
		}
		defer stmt.Close()

		var result uint
		row := stmt.QueryRowContext(ctx, args...)
		err = row.Scan(&result)
		m <- &CountChannel{
			result, err,
		}
	}(filter, countCh)
	count := <-countCh
	if count.err != nil {
		return result, err
	}
	return count.total, err
}

func (u *UserPersistence) applyLimitAndOffset(q sq.SelectBuilder, query *user.Query) (result sq.SelectBuilder, currentPage, perPage int64) {
	var page int64 = 1
	if query.Page > 0 {
		page = query.Page
	}

	perPage = 20
	if query.PerPage > 0 {
		perPage = query.PerPage
	}

	offset := (page - 1) * perPage
	return q.Limit(uint64(perPage)).Offset(uint64(offset)), page, perPage
}

func (u *UserPersistence) applySorting(q sq.SelectBuilder, query *user.Query) (result sq.SelectBuilder) {
	sortBy := string(user.UserSortByCreatedAt)
	if query.SortBy != "" {
		sortBy = string(query.SortBy)
	}

	order := string(user.UserOrderDesc)
	if query.Order != "" {
		order = string(query.Order)
	}

	return q.OrderBy(fmt.Sprintf("%s %s", sortBy, order))
}

func (u *UserPersistence) generateStatement(ctx context.Context, sql string) (*sqlx.Stmt, error) {
	stmt, err := u.db.Preparex(sql)
	// In case opration is in transaction context
	tx, ok := ctx.Value(unitofwork.TransactionContextKey).(*sqlx.Tx)
	if ok {
		stmt, err = tx.Preparex(sql)
	}

	return stmt, err
}
