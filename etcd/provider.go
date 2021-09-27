package etcd

import "github.com/ggxxll/core-sd/internal"

// Providers provide
// 		etcdv3.Client
// 		sd.Registrar
// 		sd.Instancer
// 		sd.Endpointer
// 		lb.Balancer
func Providers() []interface{} {
	return []interface{}{provideClient, provideRegistrar, provideInstancer, internal.ProvideMore}
}
