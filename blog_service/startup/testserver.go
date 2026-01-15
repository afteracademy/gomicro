package startup

import (
	"net/http/httptest"

	"github.com/afteracademy/gomicro/blog-service/config"
	"github.com/afteracademy/goserve/arch/micro"
)

type Teardown = func()

func TestServer() (micro.Router, Module, Teardown) {
	env := config.NewEnv("../.test.env", false)
	router, module, shutdown := create(env)
	ts := httptest.NewServer(router.GetEngine())
	teardown := func() {
		ts.Close()
		shutdown()
	}
	return router, module, teardown
}
