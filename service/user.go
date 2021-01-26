package service

import (
	"errors"
	"livefun/global"
	"livefun/model"
	"livefun/model/request"
	"livefun/model/response"
	"livefun/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

// Register 新建用户
func Register(u model.User) (userInter model.User, err error) {
	if !errors.Is(global.LF_DB.Where("username = ?", u.Username).First(&userInter).Error, gorm.ErrRecordNotFound) {
		return userInter, errors.New("用户名已注册")
	}

	// 附加uuid 密码md5加密  注册
	pwd := genPassword(u.Password)
	u.Password = utils.MD5V([]byte(pwd))
	u.Userid = uuid.NewV4()
	err = global.LF_DB.Create(&u).Error
	return u, err
}

func genPassword(password string) (encrypwd string) {
	return model.LeftSalt + password + model.RightSalt
}

// Login 登录
func Login(u *model.User) (userInter *model.User, err error) {
	u.Password = genPassword(u.Password)
	err = global.LF_DB.Where("username = ? AND password = ?", u.Username, u.Password).First(userInter).Error
	return
}

// Users list err error, list interface{},
func Users(info request.PageInfo) (users []response.RegisterResponse, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.LF_DB.Model(&model.User{})
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&users).Error
	return
}

// FindUserByUUID 查询用户
// @title    FindUserByUuid
// @description   Get user information by uuid, 通过uuid获取用户信息
// @auth                     （2020/04/05  20:22）
// @param     Userid query string true "用户id"
// @return    err             error
// @return    user            *model.User
func FindUserByUUID(uuid string) (user *model.User, err error) {
	var u model.User
	if err = global.LF_DB.Where("`userid` = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}
