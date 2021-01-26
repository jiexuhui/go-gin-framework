package request

// RegisterStruct User register structure
type RegisterStruct struct {
	Username  string `json:"userName"`
	Password  string `json:"passWord"`
	NickName  string `json:"nickName" gorm:"default:'QMPlusUser'"`
	HeaderImg string `json:"headerImg" gorm:"default:'http://www.henrongyi.top/avatar/lufu.jpg'"`
}
