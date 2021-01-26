package model

type JwtBlacklist struct {
	Base
	Jwt string `gorm:"type:text;comment:jwt"`
}
