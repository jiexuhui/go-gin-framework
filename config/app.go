package config

type App struct {
	Appname       string `mapstructure:"appname" json:"appname" yaml:"appname"`
	AppSecret     string `mapstructure:"appsecret" json:"appsecret" yaml:"appsecret"`
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
}
