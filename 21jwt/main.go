package main

import (
	"errors"
	"fmt"
	// "strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SecretKey = "aaa"

func main() {
	// a := strings.Split("http://10.58.122.238/cameraimages/staff/staff_1554708362966589889.jpeg", "/")
	// fmt.Println(a[len(a)-1])
	// return

	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	// 类型断言时 数字都是float64 字符串是string类型。
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = "time.Now().Unix()"
	claims["aaa"] = 123
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		fmt.Println("err:", err)
	}
	fmt.Println(tokenString)
	err = ParseToken(tokenString)
	if err != nil {
		fmt.Println("err2:", err)
	}

}
func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	}
}

func ParseToken(tokenss string) (err error) {
	token, err := jwt.Parse(tokenss, secret())
	if err != nil {
		return
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return
	}
	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	fmt.Println("claim:", claim)
	fmt.Printf("%T:\n", claim["exp"])
	fmt.Printf("%T:\n", claim["iat"])
	fmt.Printf("%T:\n", claim["aaa"])
	id, ok := claim["exp"].(float64)
	if !ok {
		err = errors.New("cannot convert claim to id")
		return
	}
	fmt.Println("id:", id)

	return
}

// func GetIdFromClaims(key string, claims jwt.Claims) string {
// 	v := reflect.ValueOf(claims)
// 	if v.Kind() == reflect.Map {
// 		for _, k := range v.MapKeys() {
// 			value := v.MapIndex(k)

// 			if fmt.Sprintf("%s", k.Interface()) == key {
// 				return fmt.Sprintf("%v", value.Interface())
// 			}
// 		}
// 	}
// 	return ""
// }
