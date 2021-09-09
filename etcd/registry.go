package etcd

import (
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
)

type RegistrarOptions struct {
	Service etcdv3.Service
}

type registrarIn struct {
	di.In

	Logger log.Logger

	Client  etcdv3.Client
	Options *RegistrarOptions
}

func provideRegistrar(in registrarIn) sd.Registrar {
	return etcdv3.NewRegistrar(in.Client, in.Options.Service, in.Logger)
}
