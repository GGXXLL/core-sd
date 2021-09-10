package etcd

import "github.com/ggxxll/core-sd/internal"

func Providers() []interface{} {
	return []interface{}{provideClient, provideRegistrar, provideInstancer, internal.ProvideMore}
}
