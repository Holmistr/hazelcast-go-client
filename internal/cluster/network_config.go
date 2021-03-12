package cluster

import (
	"fmt"
	"github.com/hazelcast/hazelcast-go-client/v4/internal/core"
	"time"
)

type NetworkConfig interface {
	Addrs() []string
	SmartRouting() bool
	ConnectionTimeout() time.Duration
}

type NetworkConfigImpl struct {
	addresses         []string
	smartRouting      bool
	connectionTimeout time.Duration
}

func NewNetworkConfigImpl() *NetworkConfigImpl {
	defaultAddr := fmt.Sprintf("%s:%d", core.DefaultHost, core.DefaultPort)
	return &NetworkConfigImpl{
		addresses:         []string{defaultAddr},
		connectionTimeout: 5 * time.Second,
	}
}

func (n NetworkConfigImpl) Addrs() []string {
	return n.addresses
}

func (n NetworkConfigImpl) ConnectionTimeout() time.Duration {
	return n.connectionTimeout
}

func (n NetworkConfigImpl) SmartRouting() bool {
	return n.smartRouting
}

type NetworkConfigBuilder interface {
	SetAddresses(addr ...string) NetworkConfigBuilder
	Config() (NetworkConfig, error)
}

type NetworkConfigProvider interface {
	Addresses() []string
	ConnectionTimeout() time.Duration
}

type NetworkConfigBuilderImpl struct {
	networkConfig *NetworkConfigImpl
	err           error
}

func NewNetworkConfigBuilderImpl() *NetworkConfigBuilderImpl {
	return &NetworkConfigBuilderImpl{
		networkConfig: NewNetworkConfigImpl(),
	}
}

func (n *NetworkConfigBuilderImpl) SetAddresses(addresses ...string) NetworkConfigBuilder {
	selfAddresses := make([]string, len(addresses))
	for i, addr := range addresses {
		if err := checkAddress(addr); err != nil {
			n.err = err
			return n
		}
		selfAddresses[i] = addr
	}
	n.networkConfig.addresses = selfAddresses
	return n
}

func (n *NetworkConfigBuilderImpl) SetConnectionTimeout(timeout time.Duration) NetworkConfigBuilder {
	n.networkConfig.connectionTimeout = timeout
	return n
}

func (n NetworkConfigBuilderImpl) Config() (NetworkConfig, error) {
	if n.err != nil {
		return n.networkConfig, n.err
	}
	return n.networkConfig, nil
}

func checkAddress(addr string) error {
	return nil
}