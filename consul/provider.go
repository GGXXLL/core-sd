package consul

import "github.com/ggxxll/core-sd/internal"

func Providers() []interface{} {
	return []interface{}{provideRegistrar, provideInstancer, provideClient, internal.ProvideMore}
}
