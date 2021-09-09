package core_sd_test

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/di"
	"github.com/DoNewsCode/core/otetcd"
	"github.com/DoNewsCode/core/srvhttp"

	core_sd "github.com/ggxxll/core-sd"
	"github.com/ggxxll/core-sd/etcd"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	kitetcd "github.com/go-kit/kit/sd/etcd"
	"github.com/go-kit/kit/sd/lb"
)

func TestEtcdRegistrar(t *testing.T) {
	c := core.Default(
		core.WithInline("name", "consul_test"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("http.addr", "192.168.82.116:8888"),
		core.WithInline("grpc.disable", true),
		core.WithInline("etcd.default.endpoints", []string{"127.0.0.1:2379"}),
	)
	defer c.Shutdown()

	c.Provide(etcd.Providers())
	c.Provide(otetcd.Providers())
	c.Provide(di.Deps{
		func() *etcd.RegistrarOptions {
			return &etcd.RegistrarOptions{Service: kitetcd.Service{
				Key:   "/services/foosvc/127.0.0.1:8888",
				Value: "http://127.0.0.1:8888/live",
			}}
		},
		func() *etcd.InstancerOption {
			return &etcd.InstancerOption{Prefix: "/services/foosvc"}
		},
	})

	c.AddModuleFunc(core_sd.NewRegistrarModule)
	c.AddModule(srvhttp.HealthCheckModule{})

	ctx, canel := context.WithCancel(context.Background())
	defer canel()
	go func() {
		err := c.Serve(ctx)
		t.Log(err)
	}()
	time.Sleep(1 * time.Second)
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

func barFactory(string) (endpoint.Endpoint, io.Closer, error) {
	return endpoint.Nop, nil, nil
}
