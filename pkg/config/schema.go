package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"regexp"
)

type Config interface {
	Validate() error
}

type ServiceConfig struct {
	Namespace   string `yaml:"namespace" validate:"required"`
	ServiceName string `yaml:"service" validate:"required"`
	PortMapping string `yaml:"portMapping" validate:"required,portmapping"`
}

type GroupConfig struct {
	Services []string `yaml:"services" validate:"required,dive,required"`
}

// validate is the validator instance used to validate structs
var validate *validator.Validate

func init() {
	// Initialize the validator
	validate = validator.New(validator.WithRequiredStructEnabled())
	// Register a custom portMapping validation
	err := validate.RegisterValidation("portmapping", validatePortMapping)
	if err != nil {
		return
	}
}

// MarshalConfig converts a Config struct to YAML data
func MarshalConfig[T Config](cfg *T) ([]byte, error) {
	if err := (*cfg).Validate(); err != nil {
		return nil, err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, fmt.Errorf("error marshaling config: %w", err)
	}
	return data, nil
}

// UnmarshalAndValidate converts YAML data into a Config struct
func UnmarshalAndValidate[T Config](data []byte) (*T, error) {
	var cfg T
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg ServiceConfig) Validate() error {
	return validate.Struct(cfg)
}

func (cfg GroupConfig) Validate() error {
	return validate.Struct(cfg)
}

// validatePortMapping is a custom validator function for port mappings
func validatePortMapping(fl validator.FieldLevel) bool {
	// A simple regex pattern to validate port mapping, e.g., "8080:80"
	portMappingPattern := regexp.MustCompile(`^\d+:\d+$`)
	return portMappingPattern.MatchString(fl.Field().String())
}
