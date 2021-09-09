package etcd

import (
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcd"
)



type RegistrarOptions struct {
	Service etcd.Service
}

type registrarIn struct {
	di.In

	Logger log.Logger

	Client  etcd.Client
	Options *RegistrarOptions
}

func provideRegistrar(in registrarIn) sd.Registrar {
	return etcd.NewRegistrar(in.Client, in.Options.Service, in.Logger)
}
