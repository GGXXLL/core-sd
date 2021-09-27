package zk

import (
	"time"

	"github.com/DoNewsCode/core/contract"
	"github.com/DoNewsCode/core/di"
	"github.com/go-kit/kit/sd/zk"
	"github.com/go-kit/log"
	stdzk "github.com/go-zookeeper/zk"
)

// Option unified import for zk.Option of kit.
// Otherwise, user need imports package like this:
// 	import (
//		"github.com/ggxxll/core-sd/zk"
//		kitzk "github.com/go-kit/kit/sd/zk"
//  )
type Option func() zk.Option

type ClientOptions struct {
	Endpoints     []string
	ClientOptions []Option
}

// ACL returns an Option specifying a non-default ACL for creating parent nodes.
func ACL(acl []stdzk.ACL) Option {
	return func() zk.Option {
		return zk.ACL(acl)
	}
}

// Credentials returns an Option specifying a user/password combination which
// the client will use to authenticate itself with.
func Credentials(user, pass string) Option {
	return func() zk.Option {
		return zk.Credentials(user, pass)
	}
}

// ConnectTimeout returns an Option specifying a non-default connection timeout
// when we try to establish a connection to a ZooKeeper server.
func ConnectTimeout(t time.Duration) Option {
	return func() zk.Option {
		return zk.ConnectTimeout(t)
	}
}

// SessionTimeout returns an Option specifying a non-default session timeout.
func SessionTimeout(t time.Duration) Option {
	return func() zk.Option {
		return zk.SessionTimeout(t)
	}
}

// Payload returns an Option specifying non-default data values for each znode
// created by CreateParentNodes.
func Payload(payload [][]byte) Option {
	return func() zk.Option {
		return zk.Payload(payload)
	}
}

// EventHandler returns an Option specifying a callback function to handle
// incoming zk.Event payloads (ZooKeeper connection events).
func EventHandler(handler func(stdzk.Event)) Option {
	return func() zk.Option {
		return zk.EventHandler(handler)
	}
}

type clientIn struct {
	di.In

	Conf   contract.ConfigAccessor
	Logger log.Logger

	Options *ClientOptions
}

func provideClient(in clientIn) (zk.Client, func(), error) {
	var opts = make([]zk.Option, len(in.Options.ClientOptions))
	for i, opt := range in.Options.ClientOptions {
		opts[i] = opt()
	}
	cli, err := zk.NewClient(in.Options.Endpoints, in.Logger, opts...)
	if err != nil {
		return nil, nil, err
	}
	return cli, func() {
		cli.Stop()
	}, nil
}
