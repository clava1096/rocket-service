package user

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"
)

func (r *repository) ExistsByLoginOrEmail(ctx context.Context, login, email string) (bool, error) {
	builderExist := sq.Select("login", "email").
		From("user").
		Where(sq.Or{
			sq.Eq{"login": login},
			sq.Eq{"email": email},
		}).
		Limit(1)

	query, args, err := builderExist.ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	err = r.db.QueryRow(ctx, query, args...).Scan(&exists)
	// todo probably this code won't find identical email\login
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}
