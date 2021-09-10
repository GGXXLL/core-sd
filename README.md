# core-sd

[DoNewsCode/core](https://github.com/DoNewsCode/core) is a service container that elegantly bootstrap and coordinate twelve-factor apps in Go.
This is a service-discovery plugin using the [DoNewsCode/core](https://github.com/DoNewsCode/core) framework and [go-kit](https://github.com/go-kit/kit).

[![Go Report Card](https://goreportcard.com/badge/github.com/ggxxll/core-sd)](https://goreportcard.com/report/github.com/ggxxll/core-sd)
[![Core Release](https://img.shields.io/github/release/DoNewsCode/core.svg)](https://github.com/DoNewsCode/core/releases/latest)

#### Support Backend
- [x] etcd-v3
- [x] consul

#### Usage


#### consul-registry

```go
package main

import (
    "context"

    "github.com/DoNewsCode/core"
    "github.com/DoNewsCode/core/contract"
    "github.com/DoNewsCode/core/di"
    "github.com/DoNewsCode/core/srvgrpc"
    "github.com/DoNewsCode/core/srvhttp"
    core_sd "github.com/ggxxll/core-sd"
    "github.com/ggxxll/core-sd/consul"
    "github.com/hashicorp/consul/api"
)

func main() {
    c := core.Default(
        core.WithInline("name", "foo"),
        core.WithInline("version", "0.0.1"),
        core.WithInline("http.addr", "127.0.0.1:8888"),
        core.WithInline("grpc.addr", "127.0.0.1:9999"),
        core.WithInline("consul.addr", "127.0.0.1:8500"),
    )
    defer c.Shutdown()

    // provide sd.Registrar and consul.Client of kit.
    c.Provide(consul.Providers())
    // provide backend of consul client.
    c.Provide(di.Deps{
        func(conf contract.ConfigAccessor) (*api.Client, error) {
            return api.NewClient(&api.Config{
                Address: conf.String("consul.addr"),
            })
        },
    })

    // provide health check module.
    c.AddModule(
        srvhttp.HealthCheckModule{},
        srvgrpc.HealthCheckModule{},
    )
    // do registrar, note: must be called before serve start.
    // first method
    c.AddModuleFunc(core_sd.NewRegistrarModule)
    // second method
    //c.Invoke(core_sd.DefaultSubscribe)


    // start server
    c.Serve(context.Background())
}
```


#### consul-discovery

```go
package core_sd_test

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ggxxll/core-sd/consul"

	"github.com/go-kit/kit/endpoint"
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
	c := core.Default(
		core.WithInline("name", "foo"),
		core.WithInline("version", "0.0.1"),
		core.WithInline("http.addr", "127.0.0.1:8888"),
		core.WithInline("grpc.addr", "127.0.0.1:9999"),
		core.WithInline("consul.addr", "127.0.0.1:8500"),
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
		// provide self 
        func() sd.Factory {
            return func(instance string) (endpoint.Endpoint, io.Closer, error) {
                return endpoint.Nop, nil, nil
            }
        },
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

    c.Invoke(func(b lb.Balancer) {
        retry := lb.Retry(3, 3*time.Second, b)

        // And now retry can be used like any other endpoint.
        req := struct{}{}
        if _, err := retry(ctx, req); err != nil {
            t.Fatal(err)
        }
    })
}
```