package user

import (
	"context"

	"github.com/clava1096/rocket-service/iam/internal/model"
	"github.com/clava1096/rocket-service/iam/internal/repository/converter"
	repoModel "github.com/clava1096/rocket-service/iam/internal/repository/model"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Create(ctx context.Context, user model.User) (model.User, error) {
	builderInsert := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("uuid", "login", "email",
			"password_hash", "notification_methods", "createdAt", "updatedAt").
		Values(user.ID, user.Login, user.Email,
			user.PasswordHash, user.NotificationMethods, user.CreatedAt, user.UpdatedAt).
		Suffix("RETURNING uuid, login, email," +
			" password_hash, notification_methods, createdAt, updatedAt")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return model.User{}, repoModel.ErrSqlFailedBuildQuery
	}

	var savedUser repoModel.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&savedUser.UUID, &savedUser.Login,
		&savedUser.Email, &savedUser.PasswordHash, &savedUser.NotificationMethods,
		&savedUser.CreatedAt, &savedUser.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		switch pgErr.ConstraintName {
		case "users_login_key": //TODO возможно не будет работать, уточнить
			return model.User{}, model.ErrUserLoginExists
		case "users_email_key":
			return model.User{}, model.ErrUserEmailExists
		default:
			return model.User{}, model.ErrUserAlreadyExists
		}
	}

	return converter.UserFromRepoModel(savedUser), nil
}
