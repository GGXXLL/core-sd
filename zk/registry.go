package zk

import (
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/zk"
	"github.com/go-kit/log"
)
// Service unified import for zk.Service of kit.
// Otherwise, user need imports package like this:
// 	import (
//		"github.com/ggxxll/core-sd/zk"
//		kitzk "github.com/go-kit/kit/sd/zk"
//  )
type Service zk.Service

// RegistrarOptions wraps args of zk.NewRegistrar func.
type RegistrarOptions struct {
	Service Service
}

type registrarIn struct {
	di.In

	Logger log.Logger

	Client  zk.Client
	Options *RegistrarOptions
}

func provideRegistrar(in registrarIn) sd.Registrar {
	service := zk.Service(in.Options.Service)
	return zk.NewRegistrar(in.Client, service, in.Logger)
}
