package startup

import (
	"context"

	"github.com/afteracademy/gomicro/blog-service/api/auth"
	authMW "github.com/afteracademy/gomicro/blog-service/api/auth/middleware"
	"github.com/afteracademy/gomicro/blog-service/api/author"
	"github.com/afteracademy/gomicro/blog-service/api/blog"
	"github.com/afteracademy/gomicro/blog-service/api/blogs"
	"github.com/afteracademy/gomicro/blog-service/api/editor"
	"github.com/afteracademy/gomicro/blog-service/api/health"
	"github.com/afteracademy/gomicro/blog-service/config"
	"github.com/afteracademy/goserve/v2/micro"
	coreMW "github.com/afteracademy/goserve/v2/middleware"
	"github.com/afteracademy/goserve/v2/mongo"
	"github.com/afteracademy/goserve/v2/network"
	"github.com/afteracademy/goserve/v2/redis"
)

type Module micro.Module[module]

type module struct {
	Context       context.Context
	Env           *config.Env
	DB            mongo.Database
	Store         redis.Store
	NatsClient    micro.NatsClient
	AuthService   auth.Service
	BlogService   blog.Service
	HealthService health.Service
}

func (m *module) GetInstance() *module {
	return m
}

func (m *module) Controllers() []micro.Controller {
	return []micro.Controller{
		health.NewController(m.HealthService),
		blog.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), m.BlogService),
		blogs.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), blogs.NewService(m.DB, m.Store)),
		author.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), author.NewService(m.DB, m.BlogService)),
		editor.NewController(m.AuthenticationProvider(), m.AuthorizationProvider(), editor.NewService(m.DB, m.AuthService)),
	}
}

func (m *module) RootMiddlewares() []network.RootMiddleware {
	return []network.RootMiddleware{
		coreMW.NewErrorCatcher(), // NOTE: this should be the first handler to be mounted
		coreMW.NewNotFound(),
	}
}

func (m *module) AuthenticationProvider() network.AuthenticationProvider {
	return authMW.NewAuthenticationProvider(m.AuthService)
}

func (m *module) AuthorizationProvider() network.AuthorizationProvider {
	return authMW.NewAuthorizationProvider(m.AuthService)
}

func NewModule(context context.Context, env *config.Env, db mongo.Database, store redis.Store, natsClient micro.NatsClient) Module {
	authService := auth.NewService(natsClient)
	blogService := blog.NewService(db, store, authService)
	healthService := health.NewService()

	return &module{
		Context:       context,
		Env:           env,
		DB:            db,
		Store:         store,
		NatsClient:    natsClient,
		AuthService:   authService,
		BlogService:   blogService,
		HealthService: healthService,
	}
}
