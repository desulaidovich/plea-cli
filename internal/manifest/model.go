package manifest

type Config struct {
	Name     string `yaml:"name"      json:"name"`
	Module   string `yaml:"module"    json:"module"`
	Output   string `yaml:"output"    json:"output"`
	LogLevel string `yaml:"log_level" json:"log_level"`
	Verbose  bool   `yaml:"verbose"   json:"verbose"`
}
