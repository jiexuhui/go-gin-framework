package initialize

import (
	"livefun/global"
	"livefun/model"
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Gorm 初始化数据库并产生数据库全局变量
func Gorm() {
	switch global.LF_CONFIG.App.DbType {
	case "mysql":
		GormMysql()
	default:
		GormMysql()
	}
}

var err error

// GormMysql 初始化Mysql数据库
func GormMysql() {
	m := global.LF_CONFIG.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	gormConfig := config(m.LogMode)
	if global.LF_DB, err = gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		global.LF_LOG.Error("MySQL启动异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(global.LF_DB)
		sqlDB, _ := global.LF_DB.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	}
}

// GormDBTables 注册数据库表专用
func GormDBTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.User{},
	)
	if err != nil {
		global.LF_LOG.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	global.LF_LOG.Info("register table success")
}

// config 根据配置决定是否开启日志
func config(mod bool) (c *gorm.Config) {
	if mod {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	} else {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	}
	return
}
