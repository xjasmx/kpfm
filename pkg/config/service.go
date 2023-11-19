package config

type Service struct {
	Alias  string
	Config *ServiceConfig
}

func NewService(alias string, namespace string, serviceName string, portMapping string) *Service {
	return &Service{
		Alias: alias,
		Config: &ServiceConfig{
			Namespace:   namespace,
			ServiceName: serviceName,
			PortMapping: portMapping,
		},
	}
}

func (s *Service) Save() error {
	cm := NewConfigManager()
	return cm.CreateServiceConfig(s.Alias, s.Config)
}

func LoadService(alias string) (*Service, error) {
	cm := NewConfigManager()
	cfg, err := cm.ReadServiceConfig(alias)
	if err != nil {
		return nil, err
	}

	return &Service{
		Alias:  alias,
		Config: cfg,
	}, nil
}
