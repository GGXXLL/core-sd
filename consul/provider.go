package consul

func Providers() []interface{} {
	return []interface{}{provideRegistrar, provideInstancer, provideClient}
}
