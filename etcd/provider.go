package etcd

func Providers() []interface{} {
	return []interface{}{provideClient, provideRegistrar, provideInstancer}
}
