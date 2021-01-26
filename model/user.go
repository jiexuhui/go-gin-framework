package model

import (
	uuid "github.com/satori/go.uuid"
)

var (
	// LeftSalt LeftSalt
	LeftSalt  = "live" 
	// RightSalt RightSalt
	RightSalt = "fun"
)

// User struct
type User struct {
	Base
	Userid    uuid.UUID `json:"userId" gorm:"comment:用户UUID"`
	Username  string    `json:"userName" gorm:"comment:用户登录名"`
	Password  string    `json:"-"  gorm:"comment:用户登录密码"`
	NickName  string    `json:"nickName" gorm:"default:系统用户;comment:用户昵称" `
	HeaderImg string    `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"`
}
