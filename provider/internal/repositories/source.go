package repositories

type Source struct {
	ExecutablePath string
	AssetsPath     string
	SpecPath       string
}

type SourceSpecification struct {
	Env []string `json:"env"`
}

type SourceRepository interface {
	GetSource() (*Source, error)
	GetSourceSpecification() (*SourceSpecification, error)
}
