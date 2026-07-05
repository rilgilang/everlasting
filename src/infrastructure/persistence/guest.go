package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"everlasting/src/domain/guest"
	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/infrastructure/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	errDomain "everlasting/src/domain/error"
)

type GuestPersistence struct {
	db     *sqlx.DB
	logger *logger.AppLogger
}

func NewGuestPersistence(db *sqlx.DB, logger *logger.AppLogger) *GuestPersistence {
	return &GuestPersistence{db: db, logger: logger}
}

func (g *GuestPersistence) GetOneByID(ctx context.Context, id string) (*guest.Guest, error) {
	query := `SELECT
		id, event_id, name, phone_number, address,
		status, is_invitation_sent, last_invitation_sent,
		created_at, updated_at, deleted_at
	FROM guests
	WHERE id = $1 AND deleted_at IS NULL`

	stmt, err := g.generateStatement(ctx, query)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_get_one_by_id", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var guestData guest.Guest
	var lastInvitationSent sql.NullTime
	err = stmt.QueryRow(id).Scan(
		&guestData.ID, &guestData.EventId, &guestData.Name,
		&guestData.PhoneNumber, &guestData.Address,
		&guestData.Status, &guestData.IsInvitationSent, &lastInvitationSent,
		&guestData.CreatedAt, &guestData.UpdatedAt, &guestData.DeletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrGuestNotFound
		}
		g.logger.Error(ctx, "postgres:guest_get_one_by_id", err.Error())
		return nil, err
	}

	if lastInvitationSent.Valid {
		guestData.LastInvitationSent = &lastInvitationSent.Time
	}

	return &guestData, nil
}

func (g *GuestPersistence) Create(ctx context.Context, guestData *guest.Guest) (*guest.Guest, error) {
	query := `INSERT INTO guests
		(id, event_id, name, phone_number, address, status,
		 is_invitation_sent, last_invitation_sent, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	stmt, err := g.generateStatement(ctx, query)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_create", err.Error())
		return nil, err
	}
	defer stmt.Close()

	now := time.Now().UTC()

	vals := []interface{}{
		guestData.ID, guestData.EventId, guestData.Name,
		guestData.PhoneNumber, guestData.Address, guestData.Status,
		guestData.IsInvitationSent, nil,
		now, now,
	}

	_, err = stmt.ExecContext(ctx, vals...)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_create", err.Error())
		return nil, err
	}

	return g.GetOneByID(ctx, guestData.ID)
}

func (g *GuestPersistence) Update(ctx context.Context, guestData *guest.Guest) (*guest.Guest, error) {
	query := `UPDATE guests SET
		event_id=$2, name=$3, phone_number=$4, address=$5,
		status=$6, is_invitation_sent=$7, last_invitation_sent=$8,
		updated_at=$9
	WHERE id=$1 AND deleted_at IS NULL
	RETURNING updated_at`

	stmt, err := g.generateStatement(ctx, query)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_update", err.Error())
		return nil, err
	}
	defer stmt.Close()

	now := time.Now().UTC()
	var lastInvitationSent interface{}
	if !guestData.LastInvitationSent.IsZero() {
		lastInvitationSent = guestData.LastInvitationSent
	}

	vals := []interface{}{
		guestData.ID, guestData.EventId, guestData.Name,
		guestData.PhoneNumber, guestData.Address, guestData.Status,
		guestData.IsInvitationSent, lastInvitationSent,
		now,
	}

	var updatedAt time.Time
	err = stmt.QueryRow(vals...).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrGuestNotFound
		}
		g.logger.Error(ctx, "postgres:guest_update", err.Error())
		return nil, err
	}

	guestData.UpdatedAt = updatedAt
	return g.GetOneByID(ctx, guestData.ID)
}

func (g *GuestPersistence) Delete(ctx context.Context, id string) error {
	now := time.Now().UTC()
	query := `UPDATE guests SET deleted_at=$2, updated_at=$2 WHERE id=$1 AND deleted_at IS NULL`

	stmt, err := g.generateStatement(ctx, query)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_delete", err.Error())
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id, now)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_delete", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_delete", err.Error())
		return err
	}

	if rowsAffected == 0 {
		return errDomain.ErrGuestNotFound
	}

	return nil
}

func (g *GuestPersistence) GetByQuery(ctx context.Context, query *guest.Query) (*guest.Guests, error) {
	filter := g.generateFilterQuery(query)

	total, err := g.getTotalCount(ctx, filter)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_get_by_query", err.Error())
		return nil, err
	}

	q := sq.Select(
		"id", "event_id", "name", "phone_number", "address",
		"status", "is_invitation_sent", "last_invitation_sent",
		"created_at", "updated_at", "deleted_at",
	).From("guests").Where(filter)

	q, page, perPage := g.applyLimitAndOffset(q, query)
	q = g.applySorting(q, query)
	q = q.PlaceholderFormat(sq.Dollar)

	toSql, args, err := q.ToSql()
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_get_by_query", err.Error())
		return nil, err
	}

	stmt, err := g.generateStatement(ctx, toSql)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_get_by_query", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		g.logger.Error(ctx, "postgres:guest_get_by_query", err.Error())
		return nil, err
	}
	defer rows.Close()

	collection := make([]guest.Guest, 0)
	for rows.Next() {
		var gd guest.Guest
		var lastInvitationSent sql.NullTime
		if err := rows.Scan(
			&gd.ID, &gd.EventId, &gd.Name,
			&gd.PhoneNumber, &gd.Address,
			&gd.Status, &gd.IsInvitationSent, &lastInvitationSent,
			&gd.CreatedAt, &gd.UpdatedAt, &gd.DeletedAt,
		); err != nil {
			g.logger.Error(ctx, "postgres:guest_get_by_query", err.Error())
			return nil, err
		}
		if lastInvitationSent.Valid {
			gd.LastInvitationSent = &lastInvitationSent.Time
		}
		collection = append(collection, gd)
	}

	return &guest.Guests{
		Collection: collection,
		Pagination: guest.Pagination{
			CurrentPage: page,
			MaxPage:     int64(math.Ceil(float64(total) / float64(perPage))),
			TotalData:   int64(total),
		},
	}, nil
}

func (g *GuestPersistence) generateFilterQuery(query *guest.Query) sq.And {
	result := sq.And{}
	result = append(result, sq.Eq{"deleted_at": nil})

	if query.EventId != "" {
		result = append(result, sq.Eq{"event_id": query.EventId})
	}

	if query.Status != "" {
		result = append(result, sq.Eq{"status": query.Status})
	}

	if query.Q != "" {
		result = append(result, sq.ILike{"name": fmt.Sprintf("%%%s%%", query.Q)})
	}

	return result
}

func (g *GuestPersistence) getTotalCount(ctx context.Context, filter sq.And) (uint, error) {
	q := sq.Select("count(id)").
		From("guests").
		Where(filter).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, err
	}

	stmt, err := g.generateStatement(ctx, sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var total uint
	err = stmt.QueryRowContext(ctx, args...).Scan(&total)
	return total, err
}

func (g *GuestPersistence) applyLimitAndOffset(q sq.SelectBuilder, query *guest.Query) (sq.SelectBuilder, int64, int64) {
	page := int64(1)
	if query.Page > 0 {
		page = query.Page
	}

	perPage := int64(20)
	if query.PerPage > 0 {
		perPage = query.PerPage
	}

	offset := (page - 1) * perPage
	return q.Limit(uint64(perPage)).Offset(uint64(offset)), page, perPage
}

func (g *GuestPersistence) applySorting(q sq.SelectBuilder, query *guest.Query) sq.SelectBuilder {
	sortBy := string(guest.GuestSortByCreatedAt)
	if query.SortBy != "" {
		sortBy = string(query.SortBy)
	}

	order := string(guest.GuestOrderDesc)
	if query.Order != "" {
		order = string(query.Order)
	}

	return q.OrderBy(fmt.Sprintf("%s %s", sortBy, order))
}

func (g *GuestPersistence) generateStatement(ctx context.Context, sql string) (*sqlx.Stmt, error) {
	stmt, err := g.db.Preparex(sql)

	tx, ok := ctx.Value(unitofwork.TransactionContextKey).(*sqlx.Tx)
	if ok && tx != nil {
		stmt, err = tx.Preparex(sql)
	}

	return stmt, err
}
