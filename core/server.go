package core

import (
	"fmt"
	"livefun/global"
	"livefun/initialize"
	"time"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

// RunWindowsServer start server
func RunWindowsServer() {
	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.LF_CONFIG.App.Addr)
	time.Sleep(10 * time.Microsecond)

	s := initServer(address, Router)
	time.Sleep(10 * time.Microsecond)
	global.LF_LOG.Info("server run success on ", zap.String("address", address))
	fmt.Printf(`
		Welcome!
		当前版本:V1.0.0 
		127.0.0.1%v
	`, address)
	s.ListenAndServe()
}
