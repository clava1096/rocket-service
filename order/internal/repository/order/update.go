package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
	repoModel "github.com/clava1096/rocket-service/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, order model.Order) (model.Order, error) {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		SetMap(sq.Eq{
			"part_uuids":       order.PartUUIDs,
			"total_price":      order.TotalPrice,
			"status":           order.Status,
			"transaction_uuid": order.TransactionUUID,
			"payment_method":   order.PaymentMethod,
			"updated_at":       time.Now(),
		}).
		Where(sq.Eq{"uuid": order.UUID}).
		Suffix(`RETURNING uuid, user_uuid, part_uuids, 
			total_price, status, transaction_uuid,
			payment_method, updated_at, created_at`)

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return model.Order{}, repoModel.ErrSqlFailedBuildQuery
	}

	var savedOrder repoModel.Order
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&savedOrder.UUID, &savedOrder.UserUUID, &savedOrder.PartUUIDs,
		&savedOrder.TotalPrice, &savedOrder.Status, &savedOrder.TransactionUUID,
		&savedOrder.PaymentMethod, &savedOrder.CreatedAt, &savedOrder.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Order{}, model.ErrOrderNotFound
		}
		return model.Order{}, fmt.Errorf("update order: %w", err)
	}

	return converter.OrderFromRepoModel(savedOrder), nil
}
