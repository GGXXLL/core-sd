package internal

import (
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"io"
)

type moreOut struct {
	di.Out

	Endpointer sd.Endpointer
	Balancer   lb.Balancer
}

type moreIn struct {
	di.In

	Logger log.Logger

	Instancer sd.Instancer
	Factory   sd.Factory `optional:"true"`
	EndpointerOptions   []sd.EndpointerOption `optional:"true"`
}

func ProvideMore(in moreIn) moreOut {
	if in.Factory == nil {
		in.Factory = func(instance string) (endpoint.Endpoint, io.Closer, error) {
			return endpoint.Nop, nil, nil
		}
	}
	endpointer := sd.NewEndpointer(in.Instancer, in.Factory, in.Logger, in.EndpointerOptions...)
	balancer := lb.NewRoundRobin(endpointer)

	return moreOut{
		Endpointer: endpointer,
		Balancer:   balancer,
	}
}
