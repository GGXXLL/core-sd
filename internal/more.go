package internal

import (
	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/lb"
	"github.com/go-kit/log"
	"io"
	"net/url"
	"strings"
)

type moreOut struct {
	di.Out

	Endpointer sd.Endpointer
	Balancer   lb.Balancer
}

type moreIn struct {
	di.In

	Logger log.Logger

	Instancer         sd.Instancer
	Factory           sd.Factory            `optional:"true"`
	EndpointerOptions []sd.EndpointerOption `optional:"true"`
}

// ProvideMore makes it easier to get sd.Endpointer and lb.Balancer
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

func ParseAddr(s string) (string, error) {
	if !strings.HasPrefix(s, "//") {
		s = "//" + s
	}
	raw, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	hostName := raw.Hostname()
	if hostName == "" {
		return "localhost:" + raw.Port(), nil
	}
	return raw.Host, nil
}

func ServiceKey(appName contract.AppName, env contract.Env, mode string) string {
	return strings.Join([]string{"services", appName.String(), env.String(), mode}, "/")
}
