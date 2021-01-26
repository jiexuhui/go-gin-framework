package main

import (
	"livefun/core"
	"livefun/global"
	"livefun/initialize"
)

func main() {
	initialize.Gorm()
	// 程序结束前关闭数据库链接
	db, _ := global.LF_DB.DB()
	defer db.Close()
	core.RunWindowsServer()
}
