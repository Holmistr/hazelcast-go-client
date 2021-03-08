package client

import "github.com/hazelcast/hazelcast-go-client/v4/internal/cluster"

type Config struct {
	ClientName  string
	ClusterName string
	Network     cluster.NetworkConfig
}

type ConfigProvider interface {
	Config() (Config, error)
}

type ConfigBuilder interface {
	SetClientName(name string) ConfigBuilder
	SetClusterName(name string) ConfigBuilder
	Network() cluster.NetworkConfigBuilder
	Config() (Config, error)
}

type ConfigBuilderImpl struct {
	config               Config
	networkConfigBuilder *cluster.NetworkConfigBuilderImpl
}

func NewConfigBuilderImpl() *ConfigBuilderImpl {
	return &ConfigBuilderImpl{
		networkConfigBuilder: cluster.NewNetworkConfigBuilderImpl(),
	}
}

func (c *ConfigBuilderImpl) SetClientName(name string) ConfigBuilder {
	c.config.ClientName = name
	return c
}

func (c *ConfigBuilderImpl) SetClusterName(name string) ConfigBuilder {
	c.config.ClusterName = name
	return c
}

func (c *ConfigBuilderImpl) Network() cluster.NetworkConfigBuilder {
	if c.networkConfigBuilder == nil {
		c.networkConfigBuilder = &cluster.NetworkConfigBuilderImpl{}
	}
	return c.networkConfigBuilder
}

func (c ConfigBuilderImpl) Config() (Config, error) {
	if c.networkConfigBuilder != nil {
		if networkConfig, err := c.networkConfigBuilder.Config(); err != nil {
			return Config{}, err
		} else {
			c.config.Network = networkConfig
		}
	}
	return c.config, nil
}
