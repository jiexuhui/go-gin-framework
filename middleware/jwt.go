package middleware

import (
	"errors"
	"livefun/global"
	"livefun/global/response"
	"livefun/model"
	"livefun/model/request"
	"livefun/service"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JWTAuth 验证JWT
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 我们这里jwt鉴权取头部信息 x-token 登录时回返回token信息 这里前端需要把token存储到cookie或者本地localSstorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		token := c.Request.Header.Get("x-token")
		if token == "" {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "未登录或非法访问", c)
			c.Abort()
			return
		}
		if service.IsBlacklist(token) {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, "您的帐户异地登陆或令牌失效", c)
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err == ErrTokenExpired {
				response.Result(response.ERROR, gin.H{
					"reload": true,
				}, "授权已过期", c)
				c.Abort()
				return
			}
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, err.Error(), c)
			c.Abort()
			return
		}
		if _, err = service.FindUserByUUID(claims.UUID.String()); err != nil {
			response.Result(response.ERROR, gin.H{
				"reload": true,
			}, err.Error(), c)
			c.Abort()
		}
		if claims.ExpiresAt-time.Now().Unix() < claims.BufferTime {
			claims.ExpiresAt = time.Now().Unix() + 60*60*24*7
			newToken, _ := j.CreateToken(*claims)
			newClaims, _ := j.ParseToken(newToken)
			c.Header("new-token", newToken)
			c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt, 10))
			if global.LF_CONFIG.App.UseMultipoint {
				err, RedisJwtToken := service.GetRedisJWT(newClaims.Username)
				if err != nil {

					global.LF_LOG.Error("get redis jwt failed", zap.Any("err", err))
				} else {
					service.JsonInBlacklist(model.JwtBlacklist{Jwt: RedisJwtToken})
					//当之前的取成功时才进行拉黑操作
				}
				// 无论如何都要记录当前的活跃状态
				_ = service.SetRedisJWT(newToken, newClaims.Username)
			}
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// JWT struct
type JWT struct {
	SigningKey []byte
}

// TokenExpired TokenExpired
var (
	ErrTokenExpired     = errors.New("Token is expired")
	ErrTokenNotValidYet = errors.New("Token not active yet")
	ErrTokenMalformed   = errors.New("That's not even a token")
	ErrTokenInvalid     = errors.New("Couldn't handle this token")
)

// NewJWT NewJWT
func NewJWT() *JWT {
	return &JWT{
		[]byte(global.LF_CONFIG.JWT.SigningKey),
	}
}

// CreateToken 创建一个token
func (j *JWT) CreateToken(claims request.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*request.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrTokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrTokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			} else {
				return nil, ErrTokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, ErrTokenInvalid
}

// 更新token
//func (j *JWT) RefreshToken(tokenString string) (string, error) {
//	jwt.TimeFunc = func() time.Time {
//		return time.Unix(0, 0)
//	}
//	token, err := jwt.ParseWithClaims(tokenString, &request.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return j.SigningKey, nil
//	})
//	if err != nil {
//		return "", err
//	}
//	if claims, ok := token.Claims.(*request.CustomClaims); ok && token.Valid {
//		jwt.TimeFunc = time.Now
//		claims.StandardClaims.ExpiresAt = time.Now().Unix() + 60*60*24*7
//		return j.CreateToken(*claims)
//	}
//	return "", TokenInvalid
//}
