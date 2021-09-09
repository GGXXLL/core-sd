package core_sd_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/srvgrpc"
	"github.com/DoNewsCode/core/srvhttp"

	core_sd "github.com/ggxxll/core-sd"
	"github.com/ggxxll/core-sd/consul"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
)

func provideConsulClient(conf contract.ConfigAccessor) (*api.Client, error) {
	return api.NewClient(&api.Config{
		Address: conf.String("consul.addr"),
	})
}

func TestConsulRegistrar(t *testing.T) {
	if os.Getenv("CONSUL_ADDR") == "" {
		fmt.Println("set CONSUL_ADDR to run test")
		return
	}

	serverIp := os.Getenv("SERVER_IP")
	if serverIp == "" {
		serverIp = "127.0.0.1"
	}

	c := core.Default(
		core.WithInline("name", "foo"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("http.addr", serverIp+":8888"),
		core.WithInline("grpc.addr", serverIp+":9999"),
		core.WithInline("consul.addr", os.Getenv("CONSUL_ADDR")),
	)
	defer c.Shutdown()

	c.Provide(consul.Providers())
	c.Provide(di.Deps{
		provideConsulClient,
	})
	c.Provide(di.Deps{
		func(appName contract.AppName, conf contract.ConfigAccessor) *consul.InstancerOption {
			return &consul.InstancerOption{
				Service: appName.String(),
				Tags: []string{
					fmt.Sprintf("version=%s", conf.String("version")),
				},
				PassingOnly: false,
			}
		},
	})

	c.AddModuleFunc(core_sd.NewRegistrarModule)
	c.AddModule(srvhttp.HealthCheckModule{})
	c.AddModule(srvgrpc.HealthCheckModule{})

	ctx, canel := context.WithCancel(context.Background())
	defer canel()
	go func() {
		_ = c.Serve(ctx)
	}()
	time.Sleep(5 * time.Second)

	c.Invoke(func(in sd.Instancer, logger log.Logger) {
		endpointer := sd.NewEndpointer(in, barFactory, logger)
		balancer := lb.NewRoundRobin(endpointer)
		retry := lb.Retry(3, 3*time.Second, balancer)

		// And now retry can be used like any other endpoint.
		req := struct{}{}
		if _, err := retry(ctx, req); err != nil {
			t.Fatal(err)
		}
	})
}
