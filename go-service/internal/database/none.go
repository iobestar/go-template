package database

type NoneConfig struct {
}

func NewNoneConfig() NoneConfig {
	return NoneConfig{}
}

func (c NoneConfig) Url() string {
	return "none://"
}
