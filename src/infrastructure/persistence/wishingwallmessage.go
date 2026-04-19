package persistence

import (
	"context"
	"database/sql"
	"everlasting/src/domain/event"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/infrastructure/pkg/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

type WishingWallMessagePersistence struct {
	db     *sqlx.DB
	logger *logger.AppLogger
}

func NewWishingWallMessagePersistence(db *sqlx.DB, logger *logger.AppLogger) *WishingWallMessagePersistence {
	return &WishingWallMessagePersistence{
		db:     db,
		logger: logger,
	}
}

// Create inserts a new wishing wall message into the database
func (w *WishingWallMessagePersistence) Create(ctx context.Context, message *event.WishingWallMessage) (*event.WishingWallMessage, error) {
	query := `INSERT INTO 
		wishing_wall_message (
			id, name, message, photo, event_id, created_at, updated_at
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7) 
		RETURNING 
			id, name, message, photo, event_id, created_at, updated_at`

	stmt, err := w.generateStatement(ctx, query)
	if err != nil {
		w.logger.Error(ctx, "postgres:wishing_wall_message_create", err.Error())
		return nil, err
	}
	defer stmt.Close()

	now := time.Now().UTC()
	vals := []interface{}{
		message.ID,
		message.Name,
		message.Message,
		message.Photo,
		message.EventID,
		now,
		now,
	}

	var messageData event.WishingWallMessage
	err = stmt.QueryRow(vals...).Scan(
		&messageData.ID,
		&messageData.Name,
		&messageData.Message,
		&messageData.Photo,
		&messageData.EventID,
		&messageData.CreatedAt,
		&messageData.UpdatedAt,
	)

	if err != nil {
		w.logger.Error(ctx, "postgres:wishing_wall_message_create", err.Error())
		return nil, err
	}

	return &messageData, nil
}

// GetAllByEventID retrieves all wishing wall messages by event ID
func (w *WishingWallMessagePersistence) GetAllByEventID(ctx context.Context, id identity.ID) ([]event.WishingWallMessage, error) {
	query := `SELECT 
		name, message, photo, created_at, updated_at 
	FROM 
		wishing_wall_message 
	WHERE 
		event_id = $1
	ORDER BY 
		created_at DESC`

	stmt, err := w.generateStatement(ctx, query)
	if err != nil {
		w.logger.Error(ctx, "postgres:wishing_wall_message_get_all_by_event_id", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return []event.WishingWallMessage{}, nil
		}
		w.logger.Error(ctx, "postgres:wishing_wall_message_get_all_by_event_id", err.Error())
		return nil, err
	}
	defer rows.Close()

	var messages []event.WishingWallMessage
	for rows.Next() {
		var message event.WishingWallMessage
		err := rows.Scan(
			&message.Name,
			&message.Message,
			&message.Photo,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			w.logger.Error(ctx, "postgres:wishing_wall_message_get_all_by_event_id", err.Error())
			return nil, err
		}
		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		w.logger.Error(ctx, "postgres:wishing_wall_message_get_all_by_event_id", err.Error())
		return nil, err
	}

	return messages, nil
}

// generateStatement prepares a statement, handling transaction context if present
func (w *WishingWallMessagePersistence) generateStatement(ctx context.Context, sql string) (*sqlx.Stmt, error) {
	stmt, err := w.db.Preparex(sql)

	// Check if operation is in transaction context
	tx, ok := ctx.Value(unitofwork.TransactionContextKey).(*sqlx.Tx)
	if ok && tx != nil {
		stmt, err = tx.Preparex(sql)
	}

	return stmt, err
}
