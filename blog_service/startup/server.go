package startup

import (
	"context"
	"time"

	"github.com/afteracademy/gomicro/blog-service/config"
	"github.com/afteracademy/goserve/v2/micro"
	"github.com/afteracademy/goserve/v2/mongo"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/redis"
	"github.com/gin-gonic/gin"
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

	dbConfig := mongo.DbConfig{
		User:        env.DBUser,
		Pwd:         env.DBUserPwd,
		Host:        env.DBHost,
		Port:        env.DBPort,
		Name:        env.DBName,
		MinPoolSize: env.DBMinPoolSize,
		MaxPoolSize: env.DBMaxPoolSize,
		Timeout:     time.Duration(env.DBQueryTimeout) * time.Second,
	}

	db := mongo.NewDatabase(context, dbConfig)
	db.Connect()

	if env.GoMode != gin.TestMode {
		EnsureDbIndexes(db)
	}

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
		Timeout:            time.Duration(env.NatsTimeoutSec) * time.Second,
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
