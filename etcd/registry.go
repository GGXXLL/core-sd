package etcd

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/ggxxll/core-sd/internal"
	"github.com/ggxxll/core-sd/util"
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
	HTTPService *Service
	GRPCService *Service
}

type registrarIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  etcdv3.Client
	Options *RegistrarOptions `optional:"true"`
}

func provideRegistrar(in registrarIn) (sd.Registrar, error) {
	if in.Options == nil {
		in.Options = &RegistrarOptions{}
	}
	r := util.Registrar{}
	if in.Options.HTTPService != nil {
		r.HTTP = etcdv3.NewRegistrar(in.Client, etcdv3.Service(*in.Options.HTTPService), in.Logger)
	} else {
		if !in.Conf.Bool("http.disable") {
			addr, err := internal.ParseAddr(in.Conf.String("http.addr"))
			if err != nil {
				return nil, err
			}
			r.HTTP = etcdv3.NewRegistrar(in.Client, etcdv3.Service(Service{
				Key:   internal.ServiceKey(in.AppName, in.Env, "http") + "/" + addr,
				Value: "http://" + addr,
				TTL:   nil,
			}), in.Logger)
		}

	}
	if in.Options.GRPCService != nil {
		r.GRPC = etcdv3.NewRegistrar(in.Client, etcdv3.Service(*in.Options.GRPCService), in.Logger)
	} else {
		if !in.Conf.Bool("grpc.disable") {
			addr, err := internal.ParseAddr(in.Conf.String("grpc.addr"))
			if err != nil {
				return nil, err
			}
			r.GRPC = etcdv3.NewRegistrar(in.Client, etcdv3.Service(Service{
				Key:   internal.ServiceKey(in.AppName, in.Env, "grpc") + "/" + addr,
				Value: addr,
				TTL:   nil,
			}), in.Logger)
		}
	}

	return &r, nil
}
