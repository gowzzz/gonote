package main

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)
type User struct{
	Id uint
	Name string
	Commons int
}
var gormdb *gorm.DB

// func newGorm() *gorm.DB {
// 	var err error
// 	args := []string{MYSQL_USER, ":", MYSQL_PASSWORD, "@", "tcp(", MYSQL_HOST, ")/", MYSQL_DB, "?charset=utf8&parseTime=True&loc=Local"}
// 	argBuf := bytes.Buffer{}
// 	for _, arg := range args {
// 		argBuf.WriteString(arg)
// 	}
// 	argsStr := argBuf.String()
// 	if GormDb == nil { //加锁是为了并发，加锁前判断是为了减少操作锁的消耗
// 		lock.Lock()
// 		defer lock.Unlock()
// 		if GormDb == nil {
// 			GormDb, err = gorm.Open("mysql", argsStr)
// 			if err != nil {
// 				GormDb = nil
// 				funcName, file, line, ok := runtime.Caller(0)
// 				if ok {
// 					log.Printf(" %s(%s:%d): %s\n", file, runtime.FuncForPC(funcName).Name(), line, err)
// 				}
// 			}
// 		}
// 	}
// 	return GormDb
// 	//只有main.go中才有必要		defer db.Close()
// }

// func GetGormDB() {
func main() {
	if gormdb==nil{

	}else{

	}
  db, err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1)/wz?charset=utf8&parseTime=True&loc=Local")
  if err!=nil{
	  panic(err)
  }
  defer db.Close()

  db=db.AutoMigrate(&User{})
  if db.Error!=nil{
	panic(db.Error)
  }
}