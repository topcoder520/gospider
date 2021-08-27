package data

type Store interface {
	Add(key string, value string) error
	BatchAdd(m map[string]string) error
	Get(key string) (string, error)
	Del(key string) error
	List(prefix string, limit ...string) ([]string, error)
}
