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
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	dom := "gy"    // the operation that the user performs on the resource.
	e := casbin.NewEnforcer("./dom_model.conf", "./dom_policy.csv")
	subs := e.GetAllSubjects()
	fmt.Println("subs:", subs)

	objs := e.GetAllObjects()
	fmt.Println("objs:", objs)

	acts := e.GetAllActions()
	fmt.Println("acts:", acts)

	roles := e.GetAllRoles()
	fmt.Println("roles:", roles)

	// ----------------------
	e.AddPolicySafe("zhuangjia", "jn", "data1", "read")
	e.RemovePolicySafe("zhuangjia", "jn", "data1", "read")
	e.AddGroupingPolicy("alice", "admin", "jn")
	e.RemoveGroupingPolicy("alice", "admin", "jn")

	polis := e.GetPolicy()
	fmt.Println("polis:", polis)

	grops := e.GetGroupingPolicy()
	fmt.Println("grops:", grops)

	//  e.AddFunction()//添加自定义方法

	return
	time.Sleep(1 * time.Second)
	// e.LoadModel()
	// e.LoadPolicy()
	if ok, err := e.EnforceSafe(sub, dom, obj, act); ok {
		fmt.Println("gy ok")
	} else {
		fmt.Println("gy deny:", err)
	}
	if ok, err := e.EnforceSafe(sub, "jn", obj, act); ok {
		fmt.Println("jn ok")
	} else {
		fmt.Println("jn deny:", err)
	}
	// p, zhuangjia, jn, data2, read
	// g, alice, admin, gy

	e.LoadPolicy()
	e.SavePolicy() //没有该方法不会写入文件
	time.Sleep(1 * time.Second)
	if ok, err := e.EnforceSafe(sub, "jn", obj, act); ok {
		fmt.Println("jn ok")
	} else {
		fmt.Println("jn deny:", err)
	}
	// e.RemovePolicy("alice", "data1", "read")
	// time.Sleep(1 * time.Second)
	// if e.Enforce(sub, obj, act) == true {
	// 	fmt.Println("ok")
	// } else {
	// 	fmt.Println("deny")
	// }
}
