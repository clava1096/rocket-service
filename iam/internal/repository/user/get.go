package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/clava1096/rocket-service/iam/internal/repository/converter"
	repoModel "github.com/clava1096/rocket-service/iam/internal/repository/model"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/json"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.User, error) {
	builderGet := sq.Select("*").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderGet.ToSql()
	if err != nil {
		return model.User{}, repoModel.ErrSqlFailedBuildQuery
	}

	var userFromTable repoModel.User

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&userFromTable.UUID, &userFromTable.Login, &userFromTable.Email,
		&userFromTable.PasswordHash, &userFromTable.NotificationMethods,
		&userFromTable.CreatedAt, &userFromTable.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, nil
		}
		return model.User{}, fmt.Errorf("error while select user: %w", err)
	}

	return converter.UserFromRepoModel(userFromTable), nil
}

func (r *repository) GetByLogin(ctx context.Context, login string) (model.User, error) {
	builderGet := sq.Select(
		"uuid",
		"login",
		"email",
		"password_hash",
		"notification_methods",
		"created_at",
		"updated_at",
	).
		From("users").
		Where(sq.Eq{"login": login}).
		Limit(1)

	query, args, err := builderGet.ToSql()
	if err != nil {
		return model.User{}, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	var notifMethodsRaw []byte

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.PasswordHash,
		&notifMethodsRaw,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return model.User{}, model.ErrUserNotFound
	}

	if err != nil {
		return model.User{}, fmt.Errorf("failed to execute query: %w", err)
	}

	if len(notifMethodsRaw) > 0 {
		err = json.Unmarshal(notifMethodsRaw, &user.NotificationMethods)
		if err != nil {
			return model.User{}, fmt.Errorf("failed to unmarshal notification methods: %w", err)
		}
	} else {
		user.NotificationMethods = []model.NotificationMethod{}
	}

	return user, nil
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	builder := sq.Select(
		"uuid", "login", "email", "password_hash",
		"notification_methods", "created_at", "updated_at",
	).From("users").Where(sq.Eq{"id": id}).Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user model.User
	var notifRaw []byte

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.ID, &user.Login, &user.Email, &user.PasswordHash,
		&notifRaw, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}

	if len(notifRaw) > 0 {
		_ = json.Unmarshal(notifRaw, &user.NotificationMethods)
	}

	return &user, nil
}
