package conf

type Config struct {
	Models map[string][]string
	Rule   Rule
}

func NewConfig() Config {
	mod := make(map[string][]string)
	mod["users"] = []string{"users1", "users2"}

	return Config{Models: mod, Rule: UseRangeRule()}
}

func UseRangeRule() *RangeRule {
	rangerule := NewRangeRule("id")
	rangerule.AddRange(500, 0, "users1")
	rangerule.AddRange(1000, 501, "users2")

	return rangerule
}
