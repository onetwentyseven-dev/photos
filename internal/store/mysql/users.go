package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"photos"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var userTableColumns = []string{
	"id", "name", "email", "ts_created", "ts_updated",
}

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) User(ctx context.Context, id uuid.UUID) (*photos.User, error) {
	var user photos.User
	query, args, err := sq.Select(userTableColumns...).From("users").Where(sq.Eq{"id": id}).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user by email query: %w", err)
	}

	err = r.db.GetContext(ctx, &user, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (r *UserRepository) UserByEmail(ctx context.Context, email string) (*photos.User, error) {

	var user photos.User
	query, args, err := sq.Select(userTableColumns...).From("users").Where(sq.Eq{"email": email}).Limit(1).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user by email query: %w", err)
	}

	err = r.db.GetContext(ctx, &user, query, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &user, err

}

func (r *UserRepository) CreateUser(ctx context.Context, user *photos.User) error {

	now := time.Now()
	user.TSCreated = now
	user.TSUpdated = now

	_, err := sq.Insert("users").SetMap(structs.Map(user)).RunWith(r.db).ExecContext(ctx)
	return err

}
