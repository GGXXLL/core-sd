package util

import (
	"github.com/go-kit/kit/sd"
	"sync"
)

type Registrar struct {
	HTTP sd.Registrar
	GRPC sd.Registrar

	onceHTTPRegister   sync.Once
	onceHTTPDeregister sync.Once
	onceGRPCRegister   sync.Once
	onceGRPCDeregister sync.Once
}

func (r *Registrar) Register() {
	if r.HTTP != nil {
		r.onceHTTPRegister.Do(r.HTTP.Register)
	}
	if r.GRPC != nil {
		r.onceGRPCRegister.Do(r.GRPC.Register)
	}
}

func (r *Registrar) Deregister() {
	if r.HTTP != nil {
		r.onceHTTPDeregister.Do(r.HTTP.Deregister)
	}
	if r.GRPC != nil {
		r.onceGRPCDeregister.Do(r.GRPC.Deregister)
	}
}
