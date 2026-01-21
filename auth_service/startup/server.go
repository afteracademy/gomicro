package startup

import (
	"context"
	"time"

	"github.com/afteracademy/gomicro/auth-service/config"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/postgres"
	"github.com/afteracademy/goserve/v2/redis"
)

type Shutdown = func()

func Server() {
	env := config.NewEnv(".env", true)
	router, _, shutdown := create(env)
	defer shutdown()
	router.Start(env.ServerHost, env.ServerPort)
}

func create(env *config.Env) (micro.Router, Module, Shutdown) {
	context := context.Background()

	dbConfig := postgres.DbConfig{
		User:        env.DBUser,
		Pwd:         env.DBUserPwd,
		Host:        env.DBHost,
		Port:        env.DBPort,
		Name:        env.DBName,
		MinPoolSize: env.DBMinPoolSize,
		MaxPoolSize: env.DBMaxPoolSize,
		Timeout:     time.Duration(env.DBQueryTimeout) * time.Second,
	}

	db := postgres.NewDatabase(context, dbConfig)
	db.Connect()

	redisConfig := redis.Config{
		Host: env.RedisHost,
		Port: env.RedisPort,
		Pwd:  env.RedisPwd,
		DB:   env.RedisDB,
	}

	store := redis.NewStore(context, &redisConfig)
	store.Connect()

	natsConfig := micro.Config{
		NatsUrl:            env.NatsUrl,
		NatsServiceName:    env.NatsServiceName,
		NatsServiceVersion: env.NatsServiceVersion,
		Timeout:            time.Second * 10,
	}

	natsClient := micro.NewNatsClient(&natsConfig)

	module := NewModule(context, env, db, store, natsClient)

	router := micro.NewRouter(env.GoMode, natsClient)
	router.RegisterValidationParsers(network.CustomTagNameFunc())
	router.LoadRootMiddlewares(module.RootMiddlewares())
	router.LoadControllers(module.Controllers())

	shutdown := func() {
		db.Disconnect()
		store.Disconnect()
		natsClient.Disconnect()
	}

	return router, module, shutdown
}
