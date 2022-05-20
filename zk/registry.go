package zk

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ggxxll/core-sd/internal"
	"github.com/ggxxll/core-sd/util"
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
	HTTPService *Service
	GRPCService *Service
}

type registrarIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  zk.Client
	Options *RegistrarOptions `optional:"true"`
}

func provideRegistrar(in registrarIn) sd.Registrar {
	if in.Options == nil {
		in.Options = &RegistrarOptions{}
	}
	r := util.Registrar{}
	if in.Options.HTTPService != nil {
		r.HTTP = zk.NewRegistrar(in.Client, zk.Service(*in.Options.HTTPService), in.Logger)
	} else {
		if !in.Conf.Bool("http.disable") {
			path := "/" + internal.ServiceKey(in.AppName, in.Env, "http")
			svc := Service{
				Path: path,
				Name: in.AppName.String() + ":" + in.Env.String(),
				Data: nil,
			}
			r.HTTP = zk.NewRegistrar(in.Client, zk.Service(svc), in.Logger)
		}

	}
	if in.Options.GRPCService != nil {
		r.GRPC = zk.NewRegistrar(in.Client, zk.Service(*in.Options.GRPCService), in.Logger)
	} else {
		if !in.Conf.Bool("grpc.disable") {
			path := "/" + internal.ServiceKey(in.AppName, in.Env, "grpc")
			svc := Service{
				Path: path,
				Name: in.AppName.String() + ":" + in.Env.String(),
				Data: nil,
			}
			r.GRPC = zk.NewRegistrar(in.Client, zk.Service(svc), in.Logger)
		}
	}
	return &r
}
