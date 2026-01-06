package model

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/clava1096/rocket-service/order/internal/model"
	converter "github.com/clava1096/rocket-service/order/internal/repository/conveter"
	repoModel "github.com/clava1096/rocket-service/order/internal/repository/model"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) Create(ctx context.Context, order model.Order) (model.Order, error) {
	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("uuid", "user_uuid", "part_uuids",
			"total_price", "status", "transaction_uuid",
			"payment_method", "created_at", "updated_at").
		Values(order.UUID, order.UserUUID, order.PartUUIDs,
			order.TotalPrice, order.Status, order.TransactionUUID,
			order.PaymentMethod, order.CreatedAt, order.UpdatedAt).
		Suffix(`RETURNING uuid, user_uuid, part_uuids,
		total_price, status, transaction_uuid,
		payment_method, created_at, updated_at`)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return model.Order{}, repoModel.ErrSqlFailedBuildQuery // todo возможно не будет работать, подумать над тем как отлавливать такие ошибки
	}

	var savedOrder repoModel.Order

	err = r.db.QueryRow(ctx, query, args...).Scan(
		&savedOrder.UUID, &savedOrder.UserUUID, &savedOrder.PartUUIDs,
		&savedOrder.TotalPrice, &savedOrder.Status, &savedOrder.TransactionUUID,
		&savedOrder.PaymentMethod, &savedOrder.CreatedAt, &savedOrder.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Order{}, model.ErrThisOrderExists
		}
		return model.Order{}, fmt.Errorf("failed to create order: %w", err) //todo не понятно как лучше сделать, возможно лучше передать какие-то туда значения?
	}

	return converter.OrderFromRepoModel(savedOrder), nil
}
