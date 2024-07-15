package repositories

type IFileCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Exists(key string) bool
}
