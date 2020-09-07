package conf

type Config struct {
	Models map[string][]string
}

func NewConfig() Config {
	mod := make(map[string][]string)
	mod["users"] = []string{"users1", "users2"}

	return Config{Models: mod}
}
