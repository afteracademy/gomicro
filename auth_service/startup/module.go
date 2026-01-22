package startup

import (
	"context"

	"github.com/afteracademy/gomicro/auth-service/api/auth"
	authMW "github.com/afteracademy/gomicro/auth-service/api/auth/middleware"
	"github.com/afteracademy/gomicro/auth-service/api/user"
	"github.com/afteracademy/gomicro/auth-service/config"
	"github.com/afteracademy/goserve/v2/micro"
	coreMW "github.com/afteracademy/goserve/v2/middleware"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/postgres"
	"github.com/afteracademy/goserve/v2/redis"
)

type Module micro.Module[module]

type module struct {
	Context     context.Context
	Env         *config.Env
	DB          postgres.Database
	Store       redis.Store
	NatsClient  micro.NatsClient
	UserService user.Service
	AuthService auth.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []micro.Controller {
	return []micro.Controller{
		auth.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.AuthService, m.UserService),
		user.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.UserService),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(), // NOTE: this should be the first handler to be mounted
		coreMW.NewNotFound(),
	}
}

func (m *module) AuthenticationProvider() network.AuthenticationProvider {
	return authMW.NewAuthenticationProvider(m.AuthService, m.UserService)
}

func (m *module) AuthorizationProvider() network.AuthorizationProvider {
	return authMW.NewAuthorizationProvider(m.AuthService)
}

func NewModule(
	context context.Context,
	env *config.Env,
	db postgres.Database,
	store redis.Store,
	natsClient micro.NatsClient,
) Module {
	userService := user.NewService(db)
	authService := auth.NewService(db, env, userService)
	return &module{
		Context:     context,
		Env:         env,
		DB:          db,
		Store:       store,
		NatsClient:  natsClient,
		UserService: userService,
		AuthService: authService,
	}
}
