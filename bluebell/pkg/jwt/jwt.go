package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

func getExpireDuration() time.Duration {
	return time.Hour * time.Duration(viper.GetInt("auth.jwt_expire"))
}

var mySercet = []byte("okokokok")
var mySecret = mySercet

type MyClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	fmt.Printf("【生成Token】使用的密钥: %s\n", string(mySecret)) // 调试输出

	c := MyClaims{
		UserID:   userID,
		Username: username, // 确保使用传入的username
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(getExpireDuration()).Unix(),
			Issuer:    "bluebell",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, err := token.SignedString(mySecret)
	if err != nil {
		fmt.Printf("【生成Token失败】错误: %v\n", err)
		return "", err
	}

	fmt.Printf("【生成Token成功】Token: %s\n", tokenString)
	return tokenString, nil
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	fmt.Printf("【解析Token】使用的密钥: %s\n", string(mySecret)) // 调试输出
	fmt.Printf("【解析Token】Token字符串: %s\n", tokenString)

	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySecret, nil
	})

	if err != nil {
		fmt.Printf("【解析Token失败】错误类型: %T, 错误信息: %v\n", err, err)

		// 详细错误分析
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("Token格式错误: %v", err)
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				return nil, fmt.Errorf("Token不可验证: %v", err)
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, fmt.Errorf("签名无效(密钥可能不匹配): %v", err)
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("Token已过期: %v", err)
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("Token尚未生效: %v", err)
			}
		}
		return nil, err
	}

	if token.Valid {
		fmt.Printf("【解析Token成功】UserID: %d, Username: %s\n", mc.UserID, mc.Username)
		return mc, nil
	}

	return nil, errors.New("invalid token")
}
