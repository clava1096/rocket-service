package model

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
	repoModel "github.com/clava1096/rocket-service/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Order, error) {

	builderGet := sq.Select("*").
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderGet.ToSql()
	if err != nil {
		return model.Order{}, repoModel.ErrSqlFailedBuildQuery
	}
	var orderFromTable repoModel.Order

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&orderFromTable.UUID, &orderFromTable.UserUUID, &orderFromTable.PartUUIDs,
		&orderFromTable.TotalPrice, &orderFromTable.Status, &orderFromTable.TransactionUUID,
		&orderFromTable.PaymentMethod, &orderFromTable.CreatedAt, &orderFromTable.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("error while select order: %w", err)
	}

	return converter.OrderFromRepoModel(orderFromTable), nil
}
