package consul

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)
// provideClient provide the consul.Client by *api.Client
func provideClient(client *api.Client) consul.Client {
	return consul.NewClient(client)
}