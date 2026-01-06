package model

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/clava1096/rocket-service/order/internal/model"
	repoModel "github.com/clava1096/rocket-service/order/internal/repository/model"
)

func (r *repository) Delete(ctx context.Context, uuid string) error {
	builderDelete := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", model.OrderStatusCancelled).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return repoModel.ErrSqlFailedBuildQuery
	}

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return repoModel.ErrSqlFailedBuildQuery
	}

	if result.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
