package app

import (
	"context"
	"log"

	iamV1API "github.com/clava1096/rocket-service/iam/internal/api/iam/v1"
	"github.com/clava1096/rocket-service/iam/internal/config"
	"github.com/clava1096/rocket-service/iam/internal/repository"
	repoSession "github.com/clava1096/rocket-service/iam/internal/repository/session"
	repoUser "github.com/clava1096/rocket-service/iam/internal/repository/user"
	"github.com/clava1096/rocket-service/iam/internal/service"
	"github.com/clava1096/rocket-service/iam/internal/service/auth"
	"github.com/clava1096/rocket-service/iam/internal/service/user"
	"github.com/clava1096/rocket-service/platform/pkg/cache"
	"github.com/clava1096/rocket-service/platform/pkg/cache/redis"
	"github.com/clava1096/rocket-service/platform/pkg/logger"
	iamV1 "github.com/clava1096/rocket-service/shared/pkg/proto/user/v1"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

type diContainer struct {
	iamV1API iamV1.UserServiceServer

	iamService     service.IAMService
	sessionService service.SessionService

	userRepository    repository.UserRepository
	sessionRepository repository.SessionRepository

	postgresPool *pgxpool.Pool

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) IamV1API(ctx context.Context) iamV1.UserServiceServer {

	if d.iamV1API == nil {
		d.iamV1API = iamV1API.NewAPI(d.IamService(ctx), d.SessionService(ctx))
	}

	return d.iamV1API
}

func (d *diContainer) SessionService(ctx context.Context) service.SessionService {

	if d.sessionService == nil {
		d.sessionService = auth.NewService(d.SessionRepository(ctx), d.UserRepository(ctx))
	}

	return d.sessionService
}

func (d *diContainer) IamService(ctx context.Context) service.IAMService {

	if d.iamService == nil {
		d.iamService = user.NewService(
			d.UserRepository(ctx),
		)
	}

	return d.iamService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.UserRepository {
	if d.userRepository == nil {
		d.userRepository = repoUser.NewRepository(d.PostgresPool(ctx))
	}

	return d.userRepository
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Fatalf("Error connecting to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("Error pinging database: %v", err)
		}
		d.postgresPool = pool
	}
	return d.postgresPool
}

func (d *diContainer) SessionRepository(_ context.Context) repository.SessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = repoSession.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}
