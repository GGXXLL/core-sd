package core_sd

import (
	"context"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/events"
	"github.com/go-kit/kit/sd"
)

type Module struct {
	registrar  sd.Registrar
	dispatcher contract.Dispatcher
}

func NewRegistrarModule(r sd.Registrar, d contract.Dispatcher) Module {
	m := Module{
		registrar:  r,
		dispatcher: d,
	}
	subscribe(d, r)

	return m
}

func subscribe(d contract.Dispatcher, reg sd.Registrar) {
	d.Subscribe(events.Listen(core.OnGRPCServerStart, func(ctx context.Context, event interface{}) error {
		reg.Register()
		return nil
	}))
	d.Subscribe(events.Listen(core.OnGRPCServerShutdown, func(ctx context.Context, event interface{}) error {
		reg.Deregister()
		return nil
	}))

	d.Subscribe(events.Listen(core.OnHTTPServerStart, func(ctx context.Context, event interface{}) error {
		reg.Register()
		return nil
	}))

	d.Subscribe(events.Listen(core.OnHTTPServerShutdown, func(ctx context.Context, event interface{}) error {
		reg.Deregister()
		return nil
	}))
}
