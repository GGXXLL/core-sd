package core_sd

import (
	"context"
	"github.com/DoNewsCode/core/di"

	"github.com/DoNewsCode/core"
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/events"
	"github.com/go-kit/kit/sd"
)

type moduleIn struct {
	di.In

	Registrar  sd.Registrar
	Dispatcher contract.Dispatcher
	Subscribe  SubscribeFunc `optional:"true"`
}

type Module struct {
	registrar  sd.Registrar
	dispatcher contract.Dispatcher
}

type SubscribeFunc func(contract.Dispatcher, sd.Registrar)

func NewRegistrarModule(in moduleIn) Module {
	m := Module{
		registrar:  in.Registrar,
		dispatcher: in.Dispatcher,
	}
	if in.Subscribe == nil {
		in.Subscribe = DefaultSubscribe
	}
	in.Subscribe(m.dispatcher, m.registrar)

	return m
}

func DefaultSubscribe(d contract.Dispatcher, reg sd.Registrar) {
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
