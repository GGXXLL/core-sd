package etcd

import (
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/log"
)

// Service unified import for zk.Service of kit.
// Otherwise, user need imports package like this:
// 	import (
//		"github.com/ggxxll/core-sd/etcd"
//		kitzk "github.com/go-kit/kit/sd/etcdv3"
//  )
type Service etcdv3.Service

// RegistrarOptions wraps args of etcdv3.NewRegistrar func.
type RegistrarOptions struct {
	Service Service
}

type registrarIn struct {
	di.In

	Logger log.Logger

	Client  etcdv3.Client
	Options *RegistrarOptions
}

func provideRegistrar(in registrarIn) sd.Registrar {
	service := etcdv3.Service(in.Options.Service)
	return etcdv3.NewRegistrar(in.Client, service, in.Logger)
}
