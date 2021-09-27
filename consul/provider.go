package consul

import "github.com/ggxxll/core-sd/internal"

// Providers provide
// 		consul.Client
// 		sd.Registrar
// 		sd.Instancer
// 		sd.Endpointer
// 		lb.Balancer
func Providers() []interface{} {
	return []interface{}{provideRegistrar, provideInstancer, provideClient, internal.ProvideMore}
}
