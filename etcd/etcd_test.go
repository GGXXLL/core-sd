package etcd_test

import (
	"context"
	"fmt"

	"os"
	"strings"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/otetcd"
	"github.com/DoNewsCode/core/srvhttp"

	core_sd "github.com/ggxxll/core-sd"
	"github.com/ggxxll/core-sd/etcd"

	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
)

func TestEtcdRegistrar(t *testing.T) {
	if os.Getenv("ETCD_ADDR") == "" {
		fmt.Println("set ETCD_ADDR to run test")
		return
	}
	serverIp := os.Getenv("SERVER_IP")
	if serverIp == "" {
		serverIp = "127.0.0.1"
	}

	c := core.Default(
		core.WithInline("name", "consul_test"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("http.addr", serverIp+":8888"),
		core.WithInline("grpc.disable", true),
		core.WithInline("etcd.default.endpoints", strings.Split(os.Getenv("ETCD_ADDR"), ",")),
	)
	defer c.Shutdown()

	c.Provide(etcd.Providers())
	c.Provide(otetcd.Providers())
	c.Provide(di.Deps{
		func() *etcd.RegistrarOptions {
			return &etcd.RegistrarOptions{Service: etcdv3.Service{
				Key:   fmt.Sprintf("/services/foosvc/%s:8888", serverIp),
				Value: fmt.Sprintf("http://%s:8888/live", serverIp),
			}}
		},
		func() *etcd.InstancerOption {
			return &etcd.InstancerOption{Prefix: "/services/foosvc"}
		},
	})

	// do registrar, note: must be called before serve start.
	// first method
	c.AddModuleFunc(core_sd.NewRegistrarModule)
	// second method
	//c.Invoke(core_sd.DefaultSubscribe)

	c.AddModule(srvhttp.HealthCheckModule{})

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
	time.Sleep(1*time.Second)
}
