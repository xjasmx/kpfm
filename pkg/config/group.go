package config

type Group struct {
	Alias  string
	Config *GroupConfig
}

func NewGroup(alias string, services []string) *Group {
	return &Group{
		Alias: alias,
		Config: &GroupConfig{
			Services: services,
		},
	}
}

func (g *Group) Save() error {
	cm := NewConfigManager()
	return cm.CreateGroupConfig(g.Alias, g.Config)
}

func LoadGroup(alias string) (*Group, error) {
	cm := NewConfigManager()
	cfg, err := cm.ReadGroupConfig(alias)
	if err != nil {
		return nil, err
	}

	return &Group{
		Alias:  alias,
		Config: cfg,
	}, nil
}
