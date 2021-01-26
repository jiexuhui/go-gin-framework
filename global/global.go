package global

import (
	"livefun/config"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	LF_CONFIG config.Server
	LF_DB     *gorm.DB
	LF_VP     *viper.Viper
	LF_LOG    *zap.Logger
	LF_REDIS  *redis.Client
)
