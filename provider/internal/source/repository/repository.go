package source_repository

type Source struct {
	ExecutablePath string
	AssetsPath     string
	SpecPath       string
}

type SourceRepository interface {
	GetSource() (*Source, error)
}
