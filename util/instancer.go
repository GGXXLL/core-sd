package util

import "github.com/go-kit/kit/sd"

type Instancer struct {
	HTTP sd.Instancer
	GRPC sd.Instancer
}

func (i *Instancer) Register(ch chan<- sd.Event) {
	if i.HTTP != nil {
		i.HTTP.Register(ch)
	}
	if i.GRPC != nil {
		i.GRPC.Register(ch)
	}
}
func (i *Instancer) Deregister(ch chan<- sd.Event) {
	if i.HTTP != nil {
		i.HTTP.Deregister(ch)
	}
	if i.GRPC != nil {
		i.GRPC.Deregister(ch)
	}
}
func (i *Instancer) Stop() {
	if i.HTTP != nil {
		i.HTTP.Stop()
	}
	if i.GRPC != nil {
		i.GRPC.Stop()
	}
}
