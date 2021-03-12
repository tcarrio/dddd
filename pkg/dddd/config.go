package dddd

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config provides configuration for dddd runner
type Config struct {
	Key        string `yaml:"key,omitempty"`
	Email      string `yaml:"email,omitempty"`
	Name       string `yaml:"name,omitempty"`
	Domain     string `yaml:"domain,omitempty"`
	IP         string `yaml:"ip,omitempty"`
	RecordType string `yaml:"recordType,omitempty"`
	LogLevel   string `yaml:"logLevel,omitempty"`
}

// ParseConfig is used to load a config from a specific file
func ParseConfig(location string) (Config, error) {
	var config Config

	file, err := os.Open(location)

	if err != nil {
		return config, fmt.Errorf("Failed to open file %s", location)
	}
	defer file.Close()

	contentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return config, fmt.Errorf("Failed to deserialize file to byte array")
	}

	err = yaml.Unmarshal(contentBytes, &config)
	if err != nil {
		return config, fmt.Errorf("Failed to parse YAML from configuration")
	}

	return config, nil
}

// MergeConfig combines the source's valid entries with the target
func MergeConfig(targetConfig *Config, sourceConfig Config) {
	key := sourceConfig.Key
	if len(key) > 0 {
		targetConfig.Key = key
	}

	email := sourceConfig.Email
	if len(key) > 0 {
		targetConfig.Email = email
	}

	name := sourceConfig.Name
	if len(key) > 0 {
		targetConfig.Name = name
	}

	domain := sourceConfig.Domain
	if len(key) > 0 {
		targetConfig.Domain = domain
	}

	ip := sourceConfig.IP
	if len(key) > 0 {
		targetConfig.IP = ip
	}

	recordType := sourceConfig.RecordType
	if len(key) > 0 {
		targetConfig.RecordType = recordType
	}

	logLevel := sourceConfig.LogLevel
	if len(key) > 0 {
		targetConfig.LogLevel = logLevel
	}
}

type missingFieldsTracker struct {
	fields []string
}

func newTracker() *missingFieldsTracker {
	return &missingFieldsTracker{
		fields: make([]string, 0),
	}
}

func (tracker *missingFieldsTracker) addField(field string) {
	tracker.fields = append(tracker.fields, field)
}

func (tracker *missingFieldsTracker) process() error {
	if len(tracker.fields) > 0 {
		return fmt.Errorf("Missing the fields: %v", tracker.fields)
	}
	return nil
}

// ValidateConfig determines whether the given config is sufficient
func ValidateConfig(config *Config) error {
	tracker := newTracker()

	if len(config.Email) == 0 {
		tracker.addField(("email"))
	}

	if len(config.Key) == 0 {
		tracker.addField("key")
	}

	if len(config.Domain) == 0 {
		tracker.addField("domain")
	}

	return tracker.process()
}
