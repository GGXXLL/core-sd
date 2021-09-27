package zk

import "github.com/ggxxll/core-sd/internal"


// Providers provide
// 		zk.Client
// 		sd.Registrar
// 		sd.Instancer
// 		sd.Endpointer
// 		lb.Balancer
func Providers() []interface{} {
	return []interface{}{provideRegistrar, provideInstancer, provideClient, internal.ProvideMore}
}
