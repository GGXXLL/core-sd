package consul_test

import (
	"context"
	"fmt"
	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/srvgrpc"
	"github.com/DoNewsCode/core/srvhttp"
	core_sd "github.com/ggxxll/core-sd"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"io"
	"os"
	"testing"
	"time"

	"github.com/ggxxll/core-sd/consul"

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
		func() sd.Factory {
			return func(instance string) (endpoint.Endpoint, io.Closer, error) {
				return endpoint.Nop, nil, nil
			}
		},
	})

	// do registrar, note: must be called before serve start.
	// first method
	c.AddModuleFunc(core_sd.NewRegistrarModule)
	// second method
	//c.Invoke(core_sd.DefaultSubscribe)

	c.AddModule(srvhttp.HealthCheckModule{})
	c.AddModule(srvgrpc.HealthCheckModule{})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		_ = c.Serve(ctx)
	}()
	time.Sleep(1 * time.Second)

	c.Invoke(func(b lb.Balancer) {
		retry := lb.Retry(3, 3*time.Second, b)

		// And now retry can be used like any other endpoint.
		req := struct{}{}
		if _, err := retry(ctx, req); err != nil {
			t.Fatal(err)
		}
	})
	cancel()
	time.Sleep(1 * time.Second)
}
