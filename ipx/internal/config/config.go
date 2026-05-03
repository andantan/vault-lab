package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RPCURL  string `yaml:"rpc_url"`
	ChainID int64  `yaml:"chain_id"`
	Keys    []Key  `yaml:"keys"`
}

type Key struct {
	Alias      string `yaml:"alias"`
	Address    string `yaml:"address"`
	PublicKey  string `yaml:"public_key"`
	PrivateKey string `yaml:"private_key"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err = yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.RPCURL == "" {
		return nil, fmt.Errorf("config: rpc_url is required")
	}
	if cfg.ChainID == 0 {
		return nil, fmt.Errorf("config: chain_id is required")
	}

	return &cfg, nil
}

func (c *Config) KeyByAddress(address string) (*Key, error) {
	for i := range c.Keys {
		if strings.EqualFold(c.Keys[i].Address, address) {
			return &c.Keys[i], nil
		}
	}
	return nil, fmt.Errorf("config: key with address %q not found", address)
}
