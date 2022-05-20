package consul

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"

	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/log"

	"github.com/hashicorp/consul/api"
)

// RegisterOption wrap original options
type RegisterOption struct {
	AgentServiceRegistration            *api.AgentServiceRegistration
	AgentServiceRegistrationInterceptor AgentServiceRegistrationInterceptor
}

type registrarIn struct {
	di.In

	AppName contract.AppName
	Env     contract.Env
	Conf    contract.ConfigAccessor
	Logger  log.Logger

	Client  consul.Client
	Options *RegisterOption `optional:"true"`
}

type AgentServiceRegistrationInterceptor func(*api.AgentServiceRegistration)

func provideRegistrar(in registrarIn) (sd.Registrar, error) {
	if in.Options == nil {
		in.Options = &RegisterOption{
			AgentServiceRegistration: nil,
			AgentServiceRegistrationInterceptor: func(registration *api.AgentServiceRegistration) {

			},
		}
	}
	ID := in.AppName.String() + ":" + in.Env.String()
	if in.Options.AgentServiceRegistration == nil {
		checks := make([]*api.AgentServiceCheck, 0)
		endpoints := make(map[string]string)
		if addr := in.Conf.String("http.addr"); addr != "" && !in.Conf.Bool("http.disable") {
			endpoints["http"] = "//" + addr
			checks = append(checks, &api.AgentServiceCheck{
				CheckID:  ID + ":http",
				Status:   api.HealthPassing,
				Interval: "10s",
				HTTP:     "http://" + addr + "/live",
			})
		}
		if addr := in.Conf.String("grpc.addr"); addr != "" && !in.Conf.Bool("grpc.disable") {
			endpoints["grpc"] = "//" + addr
			checks = append(checks, &api.AgentServiceCheck{
				CheckID:  ID + ":grpc",
				Status:   api.HealthPassing,
				Interval: "10s",
				GRPC:     addr,
			})
		}
		if len(endpoints) == 0 {
			return nil, fmt.Errorf("endpoints is empty")
		}
		var (
			addr string
			port uint64
		)
		addresses := make(map[string]api.ServiceAddress, len(endpoints))
		for name, endpoint := range endpoints {
			raw, err := url.Parse(endpoint)
			if err != nil {
				return nil, err
			}
			addr = raw.Hostname()
			port, _ = strconv.ParseUint(raw.Port(), 10, 16)
			addresses[name] = api.ServiceAddress{Address: endpoint, Port: int(port)}
		}

		in.Options.AgentServiceRegistration = &api.AgentServiceRegistration{
			ID:   ID,
			Name: ID,
			Tags: []string{
				fmt.Sprintf("version=%s", in.Conf.String("version")),
			},
			Port:            int(port),
			Address:         addr,
			TaggedAddresses: addresses,
			Checks:          checks,
		}
	}
	in.Options.AgentServiceRegistrationInterceptor(in.Options.AgentServiceRegistration)

	reg := consul.NewRegistrar(in.Client, in.Options.AgentServiceRegistration, in.Logger)
	return reg, nil
}
