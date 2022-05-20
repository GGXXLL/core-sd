package util

import "github.com/go-kit/kit/sd"

type Registrar struct {
	HTTP sd.Registrar
	GRPC sd.Registrar
}

func (r *Registrar) Register() {
	if r.HTTP != nil {
		r.HTTP.Register()
	}
	if r.GRPC != nil {
		r.GRPC.Register()
	}
}

func (r *Registrar) Deregister() {
	if r.HTTP != nil {
		r.HTTP.Deregister()
	}
	if r.GRPC != nil {
		r.GRPC.Deregister()
	}
}
