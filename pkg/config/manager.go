package config

import (
	"fmt"
)

const groupPrefix = "group-"
const servicePrefix = "service-"

type Manager struct{}

func NewConfigManager() *Manager {
	return &Manager{}
}

func (cm *Manager) CreateServiceConfig(alias string, cfg *ServiceConfig) error {
	config, err := MarshalConfig(cfg)
	if err != nil {
		return err
	}
	return WriteConfig(servicePrefix+alias, config)
}

func (cm *Manager) CreateGroupConfig(alias string, cfg *GroupConfig) error {
	config, err := MarshalConfig(cfg)
	if err != nil {
		return err
	}
	return WriteConfig(groupPrefix+alias, config)
}

func (cm *Manager) ReadServiceConfig(alias string) (*ServiceConfig, error) {
	data, err := ReadConfig(servicePrefix + alias)
	if err != nil {
		return nil, fmt.Errorf("error reading service config: %w", err)
	}
	return UnmarshalAndValidate[ServiceConfig](data)
}

func (cm *Manager) ReadGroupConfig(alias string) (*GroupConfig, error) {
	data, err := ReadConfig(groupPrefix + alias)
	if err != nil {
		return nil, fmt.Errorf("error reading group config: %w", err)
	}
	return UnmarshalAndValidate[GroupConfig](data)
}

func (cm *Manager) DeleteServiceConfig(alias string) error {
	return DeleteConfig(servicePrefix + alias)
}

func (cm *Manager) DeleteGroupConfig(alias string) error {
	return DeleteConfig(groupPrefix + alias)
}
