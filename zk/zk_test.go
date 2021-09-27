package zk_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/srvgrpc"
	"github.com/DoNewsCode/core/srvhttp"

	core_sd "github.com/ggxxll/core-sd"
	"github.com/ggxxll/core-sd/zk"

	"github.com/go-kit/kit/sd/lb"
)

func TestZk(t *testing.T) {
	if os.Getenv("ZK_ADDR") == "" {
		fmt.Println("set ZK_ADDR to run test")
		return
	}

	serverIp := os.Getenv("SERVER_IP")
	if serverIp == "" {
		serverIp = "127.0.0.1"
	}

	var path = "/core-sd.zk/test"

	c := core.Default(
		core.WithInline("name", "foo"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("log.level", "none"),
		core.WithInline("http.addr", serverIp+":8888"),
		core.WithInline("grpc.addr", serverIp+":9999"),
		core.WithInline("zookeeper.addr", strings.Split(os.Getenv("ZK_ADDR"), ",")),
	)
	defer c.Shutdown()

	c.Provide(zk.Providers())
	c.Provide(di.Deps{
		// provide this option for create zk.Client
		func(conf contract.ConfigAccessor) *zk.ClientOptions {
			payload := [][]byte{[]byte("Payload"), []byte("Test")}
			return &zk.ClientOptions{
				Endpoints: conf.Strings("zookeeper.addr"),
				ClientOptions: []zk.Option{
					zk.Payload(payload),
				},
			}
		},
		// provide this option for define the registrar info
		func() *zk.RegistrarOptions {
			return &zk.RegistrarOptions{
				Service: zk.Service{
					Path: path,
					Name: "core-sd",
					Data: []byte("just some payload"),
				},
			}
		},
		// provide this option for discover service info by path
		func(appName contract.AppName, conf contract.ConfigAccessor) *zk.InstancerOption {
			return &zk.InstancerOption{
				Path: path,
			}
		},
		func() core_sd.SubscribeFunc {
			return core_sd.SubscribeGRPC
		},
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
