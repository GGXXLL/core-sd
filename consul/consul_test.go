package consul_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/srvgrpc"
	"github.com/DoNewsCode/core/srvhttp"

	core_sd "github.com/ggxxll/core-sd"
	"github.com/ggxxll/core-sd/consul"

	"github.com/go-kit/kit/sd/lb"

	"github.com/hashicorp/consul/api"
)

func provideConsulClient(conf contract.ConfigAccessor) (*api.Client, error) {
	return api.NewClient(&api.Config{
		Address: conf.String("consul.addr"),
	})
}

func TestConsul(t *testing.T) {
	if os.Getenv("CONSUL_ADDR") == "" {
		fmt.Println("set CONSUL_ADDR to run test")
		return
	}

	serverIp := os.Getenv("SERVER_IP")
	if serverIp == "" {
		serverIp = "127.0.0.1"
	}

	c := core.Default(
		core.WithInline("name", "app"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("log.level", "none"),
		core.WithInline("http.addr", serverIp+":8888"),
		core.WithInline("grpc.addr", serverIp+":9999"),
		core.WithInline("consul.addr", os.Getenv("CONSUL_ADDR")),
	)
	defer c.Shutdown()

	c.Provide(consul.Providers())
	c.Provide(di.Deps{
		provideConsulClient,
	})

	// do registrar, note: must be called before serve start.
	// first method
	c.AddModuleFunc(core_sd.NewRegistrarModule)
	// second method
	//c.Invoke(core_sd.DefaultSubscribe)

	c.AddModule(srvhttp.HealthCheckModule{})
	c.AddModule(srvgrpc.HealthCheckModule{})

	var g sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(1)
	go func() {
		defer g.Done()
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
	g.Wait()
}
