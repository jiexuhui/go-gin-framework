package config

// Server struct
type Server struct {
	App   App   `mapstructure:"app" json:"app" yaml:"app"`
	Zap   Zap   `mapstructure:"zap" json:"zap" yaml:"zap"`
	Mysql Mysql `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	JWT   JWT   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis Redis `mapstructure:"redis" json:"redis" yaml:"redis"`
}
