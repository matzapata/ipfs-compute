package domain

type Source struct {
	ExecutablePath string
	AssetsPath     string
	SpecPath       string
}

type SourceSpecification struct {
	Env []string `json:"env"`
}

type ISourceService interface {
	GetSource() (*Source, error)
	GetSourceSpecification() (*SourceSpecification, error)
}

type ISourceRepository interface {
	GetSource() (*Source, error)
	GetSourceSpecification() (*SourceSpecification, error)
}
