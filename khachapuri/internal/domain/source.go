package domain

type Source struct {
	SourcePath string
	ConfigPath string
}

type Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Service struct {
	Build struct {
		Context string `yaml:"context"`
		Command string `yaml:"command"`
	} `yaml:"build"`
	Executable string   `yaml:"executable"`
	Assets     []string `yaml:"assets"`
	Env        []Env    `yaml:"env"`
	EnvFile    []string `yaml:"env_file"`
}

type Provider struct {
	Name   string `yaml:"name"`
	Domain string `yaml:"domain"`
}

type SourceSpecification struct {
	Version   string             `yaml:"version"`
	Services  map[string]Service `yaml:"services"`
	Providers struct {
		Allowlist []Provider `yaml:"allowlist"`
	} `yaml:"providers"`
}

type ISourceService interface {
	GetSource() (*Source, error)
	GetSourceSpecification() (*SourceSpecification, error)
}

type ISourceRepository interface {
	GetSource() (*Source, error)
	GetSourceSpecification() (*SourceSpecification, error)
}
