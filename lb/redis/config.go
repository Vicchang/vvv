package redis

type Config struct {
	Type   string `yaml:"type"`
	Nodes  string `yaml:"nodes"`
	Master string `yaml:"master"`
}
