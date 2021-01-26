package v1

import (
	"fmt"
	"livefun/global"
	"livefun/global/response"
	"livefun/middleware"
	"livefun/model"
	"livefun/model/request"
	resp "livefun/model/response"
	"livefun/service"
	"livefun/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// Register 注册
func Register(c *gin.Context) {
	var R request.RegisterStruct
	err := c.ShouldBind(&R)

	UserVerify := utils.Rules{
		"Username": {utils.NotEmpty()},
		"Password": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(R, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	user := &model.User{Username: R.Username, NickName: R.NickName, Password: R.Password}
	userInfo, err := service.Register(*user)
	if err != nil {
		response.FailWithDetailed(response.ERROR, userInfo, fmt.Sprintf("%v", err), c)
	} else {
		response.OkDetailed(userInfo, "注册成功", c)
	}
}

// Login 登录
// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body model.User true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func Login(c *gin.Context) {
	var L model.User
	_ = c.ShouldBind(&L)

	UserVerify := utils.Rules{
		"Username": {utils.NotEmpty()},
		"Password": {utils.NotEmpty()},
	}

	UserVerifyErr := utils.Verify(L, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	U := &model.User{Username: L.Username, Password: L.Password}
	if user, err := service.Login(U); err != nil {
		response.FailWithMessage(fmt.Sprintf("用户名密码错误或%v", err), c)
	} else {
		tokenNext(c, *user)
	}
}

// 登录以后签发jwt
func tokenNext(c *gin.Context, user model.User) {
	j := &middleware.JWT{
		SigningKey: []byte(global.LF_CONFIG.JWT.SigningKey), // 唯一签名
	}
	clams := request.CustomClaims{
		UUID:       user.Userid,
		ID:         user.ID,
		NickName:   user.NickName,
		Username:   user.Username,
		BufferTime: 60 * 60 * 24, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,       // 签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*7, // 过期时间 7天
			Issuer:    "liveFun",                      // 签名的发行者
		},
	}
	token, err := j.CreateToken(clams)
	if err != nil {
		response.FailWithMessage("获取token失败", c)
		return
	}
	if !global.LF_CONFIG.App.UseMultipoint {
		response.OkWithData(resp.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
		}, c)
		return
	}
	err, jwtStr := service.GetRedisJWT(user.Username)
	if err == redis.Nil {
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithData(resp.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
		}, c)
	} else if err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		var blackJWT model.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := service.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := service.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithData(resp.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: clams.StandardClaims.ExpiresAt * 1000,
		}, c)
	}
}

// Users 用户列表
func Users(c *gin.Context) {
	pageInfo := request.PageInfo{
		Page:     1,
		PageSize: 10,
	}
	users, total, err := service.Users(pageInfo)
	if err != nil {
		response.FailWithDetailed(response.ERROR, users, fmt.Sprintf("%v", err), c)
	} else {
		response.OkWithData(resp.PageResult{
			List:     users,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, c)
	}
}
