package response

import "livefun/model"

// RegisterResponse 注册返回
type RegisterResponse struct {
	Userid   string `json:"uuid"`
	Username string `json:"userName"`
	Password string `json:"-"`
	NickName string `json:"nickName"`
}

// LoginResponse 登录返回
type LoginResponse struct {
	User      model.User `json:"user"`
	Token     string     `json:"token"`
	ExpiresAt int64      `json:"expiresAt"`
}
