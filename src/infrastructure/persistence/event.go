package persistence

import (
	"context"
	"database/sql"
	"everlasting/src/domain/event"
	sq "github.com/Masterminds/squirrel"
	"time"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/infrastructure/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type EventPersistence struct {
	db     *sqlx.DB
	logger *logger.AppLogger
}

func NewEventPersistence(db *sqlx.DB, logger *logger.AppLogger) (result *EventPersistence) {
	return &EventPersistence{db, logger}
}

// GetOneByID retrieves an event by its ID
func (e *EventPersistence) GetOneByID(ctx context.Context, id identity.ID) (*event.Event, error) {
	query := `SELECT 
		id, title, description, 'date', 'time', location,
		category, messages, max_messages, image, status,
		organizer, created_at, updated_at 
	FROM 
		events 
	WHERE id = $1`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_get_one_by_id", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var eventData event.Event
	err = stmt.QueryRow(id).Scan(
		&eventData.ID, &eventData.Title, &eventData.Description, &eventData.Date, &eventData.Time, &eventData.Location,
		&eventData.Category, &eventData.Messages, &eventData.MaxMessages, &eventData.Image, &eventData.Status,
		&eventData.Organizer, &eventData.CreatedAt, &eventData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrEventNotFound
		}
		e.logger.Error(ctx, "postgres:events_get_one_by_id", err.Error())
		return nil, err
	}
	return &eventData, nil
}

// GetOneByTitle retrieves an event by its title
func (e *EventPersistence) GetOneByTitle(ctx context.Context, title string) (*event.Event, error) {
	query := `SELECT 
		id, title, description, time, location,
		category, messages, max_messages, image, status,
		organizer, created_at, updated_at 
	FROM 
		events 
	WHERE title = $1`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_get_one_by_title", err.Error())
		return nil, err
	}
	defer stmt.Close()

	var eventData event.Event
	err = stmt.QueryRow(title).Scan(
		&eventData.ID, &eventData.Title, &eventData.Description, &eventData.Time, &eventData.Location,
		&eventData.Category, &eventData.Messages, &eventData.MaxMessages, &eventData.Image, &eventData.Status,
		&eventData.Organizer, &eventData.CreatedAt, &eventData.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrEventNotFound
		}
		e.logger.Error(ctx, "postgres:events_get_one_by_title", err.Error())
		return nil, err
	}
	return &eventData, nil
}

// Create inserts a new event into the database
func (e *EventPersistence) Create(ctx context.Context, event *event.Event) (*event.Event, error) {
	now := time.Now().UTC()
	query := `INSERT INTO 
		events (
			id, title, description, date, time, location,
			category, messages, max_messages, image, status,
			organizer, created_at, updated_at
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) 
		RETURNING 
			id, created_at, updated_at`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_create", err.Error())
		return nil, err
	}
	defer stmt.Close()

	vals := []interface{}{
		event.ID,
		event.Title,
		event.Description,
		event.Date,
		event.Time,
		event.Location,
		event.Category,
		event.Messages,
		event.MaxMessages,
		event.Image,
		event.Status,
		event.Organizer,
		now,
		now,
	}

	var id string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRow(vals...).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_create", err.Error())
		return nil, err
	}

	event.CreatedAt = createdAt
	event.UpdatedAt = updatedAt
	event.ID = id

	return e.GetOneByID(ctx, identity.FromStringOrNil(id))
}

// Update updates an existing event
func (e *EventPersistence) Update(ctx context.Context, event *event.Event) (*event.Event, error) {
	now := time.Now().UTC()
	query := `UPDATE 
		events
	SET
		title=$2, description=$3, date=$4, time=$5, location=$6,
		category=$7, max_messages=$8, image=$9,
		status=$10, organizer=$11, updated_at=$12
	WHERE
		id=$1
	RETURNING
		updated_at`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_update", err.Error())
		return nil, err
	}
	defer stmt.Close()

	vals := []interface{}{
		event.ID,
		event.Date,
		event.Title,
		event.Description,
		event.Time,
		event.Location,
		event.Category,
		event.MaxMessages,
		event.Image,
		event.Status,
		event.Organizer,
		now,
	}

	var updatedAt time.Time
	err = stmt.QueryRow(vals...).Scan(&updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errDomain.ErrEventNotFound
		}
		e.logger.Error(ctx, "postgres:events_update", err.Error())
		return nil, err
	}

	event.UpdatedAt = updatedAt
	return e.GetOneByID(ctx, identity.FromStringOrNil(event.ID))
}

// Delete soft deletes an event by updating its status
func (e *EventPersistence) Delete(ctx context.Context, id identity.ID) error {
	query := `UPDATE events SET status = 'deleted', updated_at = $2 WHERE id = $1`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_delete", err.Error())
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC()
	result, err := stmt.Exec(id, now)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_delete", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		e.logger.Error(ctx, "postgres:events_delete", err.Error())
		return err
	}

	if rowsAffected == 0 {
		return errDomain.ErrEventNotFound
	}

	return nil
}

// HardDelete permanently removes an event from the database
func (e *EventPersistence) HardDelete(ctx context.Context, id identity.ID) error {
	query := `DELETE FROM events WHERE id = $1`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_hard_delete", err.Error())
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_hard_delete", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		e.logger.Error(ctx, "postgres:events_hard_delete", err.Error())
		return err
	}

	if rowsAffected == 0 {
		return errDomain.ErrEventNotFound
	}

	return nil
}

// GetByQuery retrieves events based on query parameters
func (e *EventPersistence) GetByQuery(ctx context.Context, query *event.Query) (*event.Events, error) {
	// Generate filter query
	filter := e.generateFilterQuery(query)

	// Generate data collection query
	q := sq.Select("id", "title", "description", "time", "location",
		"category", "messages", "max_messages", "image", "status",
		"organizer", "created_at", "updated_at").
		From("events").
		Where(filter)

	// Per page limitation
	perPage := 20
	if query.PerPage > 0 {
		perPage = int(query.PerPage)
	}
	q = q.Limit(uint64(perPage))

	// Ordering
	q = q.OrderBy("created_at DESC")

	// Return as $1, $2 format
	q = q.PlaceholderFormat(sq.Dollar)

	sql, args, err := q.ToSql()
	if err != nil {
		e.logger.Error(ctx, "postgres:events_get_by_query", err.Error())
		return nil, err
	}

	stmt, err := e.generateStatement(ctx, sql)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_get_by_query", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_get_by_query", err.Error())
		return nil, err
	}
	defer rows.Close()

	// Collect result
	collection := make([]event.Event, 0)
	var nextCursor int64
	for rows.Next() {
		var eventData event.Event
		if err := rows.Scan(&eventData.ID, &eventData.Title, &eventData.Description, &eventData.Time,
			&eventData.Location, &eventData.Category, &eventData.Messages, &eventData.MaxMessages,
			&eventData.Image, &eventData.Status, &eventData.Organizer, &eventData.CreatedAt,
			&eventData.UpdatedAt); err != nil {
			e.logger.Error(ctx, "postgres:events_get_by_query", err.Error())
			return nil, err
		}
		collection = append(collection, eventData)
		nextCursor = eventData.CreatedAt.UnixMilli()
	}

	return &event.Events{
		Collection: collection,
		Pagination: event.Pagination{
			NextCursor: nextCursor,
		},
	}, nil
}

// UpdateMessagesCount updates the messages count for an event
func (e *EventPersistence) UpdateMessagesCount(ctx context.Context, id identity.ID, messagesCount int) error {
	query := `UPDATE events SET messages = $2, updated_at = $3 WHERE id = $1`

	stmt, err := e.generateStatement(ctx, query)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_update_messages_count", err.Error())
		return err
	}
	defer stmt.Close()

	now := time.Now().UTC()
	result, err := stmt.Exec(id, messagesCount, now)
	if err != nil {
		e.logger.Error(ctx, "postgres:events_update_messages_count", err.Error())
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		e.logger.Error(ctx, "postgres:events_update_messages_count", err.Error())
		return err
	}

	if rowsAffected == 0 {
		return errDomain.ErrEventNotFound
	}

	return nil
}

// generateFilterQuery creates the WHERE clause based on query parameters
func (e *EventPersistence) generateFilterQuery(query *event.Query) sq.And {
	result := sq.And{}

	// Exclude deleted events by default
	result = append(result, sq.NotEq{"status": "deleted"})

	if query.Category != "" {
		result = append(result, sq.Eq{"category": query.Category})
	}

	if query.Status != "" && query.Status != "deleted" {
		result = append(result, sq.Eq{"status": query.Status})
	}

	if query.Organizer != "" {
		result = append(result, sq.Eq{"organizer": query.Organizer})
	}

	if query.Location != "" {
		result = append(result, sq.ILike{"location": "%" + query.Location + "%"})
	}

	if query.Cursor > 0 {
		result = append(result, sq.Lt{"created_at": time.UnixMilli(query.Cursor).UTC()})
	}

	if query.DateFrom != "" {
		dtFromTime, err := time.Parse("2006-01-02", query.DateFrom)
		if err == nil {
			result = append(result, sq.GtOrEq{"created_at": dtFromTime})
		}
	}

	if query.DateUntil != "" {
		dtUntilTime, err := time.Parse("2006-01-02", query.DateUntil)
		if err == nil {
			// Set to end of day
			dtUntilTime = dtUntilTime.Add(24*time.Hour - time.Second)
			result = append(result, sq.LtOrEq{"created_at": dtUntilTime})
		}
	}

	return result
}

// generateStatement prepares a statement, handling transaction context if present
func (e *EventPersistence) generateStatement(ctx context.Context, sql string) (*sqlx.Stmt, error) {
	stmt, err := e.db.Preparex(sql)

	// Check if operation is in transaction context
	tx, ok := ctx.Value(unitofwork.TransactionContextKey).(*sqlx.Tx)
	if ok && tx != nil {
		stmt, err = tx.Preparex(sql)
	}

	return stmt, err
}
