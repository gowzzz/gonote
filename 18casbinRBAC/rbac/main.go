package main

// http://www.cnblogs.com/wang_yb/archive/2018/11/20/9987397.html
import (
	"fmt"
	"time"

	"github.com/casbin/casbin"
)

var modepath = "./rbac_model.conf"
var csvpath = "./rbac_policy.csv"

/*
// 禁用AutoSave机制
e.EnableAutoSave(false)

// 因为禁用了AutoSave，当前策略的改变只在内存中生效
// 这些策略在持久层中仍是不变的
e.AddPolicy(...)
e.RemovePolicy(...)

// 开启AutoSave机制
e.EnableAutoSave(true)

*/
func main() {
	sub := "alice" // the user that wants to access a resource.
	obj := "data2" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	e := casbin.NewEnforcer("./rbac_model.conf", "./rbac_policy.csv")
	time.Sleep(1 * time.Second)
	// e.LoadModel()
	// e.LoadPolicy()
	if e.Enforce(sub, obj, act) == true {
		fmt.Println("ok")
	} else {
		fmt.Println("deny")
	}
	// g, alice, data2_admin
	// p, data2_admin, data2, write
	// e.AddPolicy("data2_admin", "data2", "write")
	e.AddGroupingPolicy("alice", "data2_admin")

	// e.SavePolicy() //没有该方法不会写入文件
	time.Sleep(1 * time.Second)
	if e.Enforce(sub, obj, act) == true {
		fmt.Println("ok")
	} else {
		fmt.Println("deny")
	}
	// e.RemovePolicy("alice", "data1", "read")
	// time.Sleep(1 * time.Second)
	// if e.Enforce(sub, obj, act) == true {
	// 	fmt.Println("ok")
	// } else {
	// 	fmt.Println("deny")
	// }
}
