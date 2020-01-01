package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var (
	SecretKey = "SecretKey"
)

func SetJwtSecretKey(key string) {
	SecretKey = key
}

func GenerateJwt(cid, uid, uname string, expire int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"cid": cid,
		"uid": uid,
		"uname": uname,
		"exp": time.Now().Add(time.Second * time.Duration(expire)),
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to generate token, cid = %s, uid = %s, uname = %s", cid, uid, uname)
	}
	return tokenString, nil
}

func ExplainJwt(t string, cid, uid, uname *string) error {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		*cid = claims["cid"].(string)
		*uid = claims["uid"].(string)
		*uname = claims["uname"].(string)
		return nil
	}
	return fmt.Errorf("invalid token")
}


