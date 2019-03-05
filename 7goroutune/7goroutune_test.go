package a

import (
	"testing"
)

// go test -run=xxx -bench=. -benchtime="3s" -cpuprofile profile_cpu.out
/*
https://graphviz.gitlab.io/_pages/Download/Download_windows.html
# go tool pprof -svg profile_cpu.out > profile_cpu.svg
# go tool pprof -pdf profile_cpu.out > profile_cpu.pdf
*/
func BenchmarkGoroutune(b *testing.B) {
	// 在report中包含内存分配信息    BenchmarkStringJoin1-4 300000 4351 ns/op 32 B/op 2 allocs/op
	/*
		-4表示4个CPU线程执行；
		300000表示总共执行了30万次；
		4531ns/op，表示每次执行耗时4531纳秒；
		32B/op表示每次执行分配了32字节内存；
		2 allocs/op表示每次执行分配了2次对象。
	*/
	b.ReportAllocs()
	p := NewPool()
	for i := 0; i < b.N; i++ {
		Runpool(p, i)
	}
}
func BenchmarkStringJoin2(b *testing.B) {
	b.ReportAllocs()
	input := []string{"Hello", "World"}
	join := func(strs []string, delim string) string {
		if len(strs) == 2 {
			return strs[0] + delim + strs[1]
		}
		return ""
	}
	for i := 0; i < b.N; i++ {
		result := join(input, " ")
		if result != "Hello World" {
			b.Error("Unexpected result: " + result)
		}
	}
}
func BenchmarkStringJoin2B(b *testing.B) {
	b.ReportAllocs()
	join := func(strs []string, delim string) string {
		if len(strs) == 2 {
			return strs[0] + delim + strs[1]
		}
		return ""
	}
	for i := 0; i < b.N; i++ {
		input := []string{"Hello", "World"}
		result := join(input, " ")
		if result != "Hello World" {
			b.Error("Unexpected result: " + result)
		}
	}
}

/*
只关心性能最差的那个

go test -run=xxx -bench=BenchmarkStringJoin2B$ -cpuprofile profile_2b.out
go test -run=xxx -bench=BenchmarkStringJoin2$ -cpuprofile profile_2.out
go tool pprof -svg profile_2b.out > profile_2b.svg
go tool pprof -svg profile_2.out > profile_2.svg

go test -run=xxx -bench=BenchmarkGoroutune$ -cpuprofile grt.out
go tool pprof -pdf grt.out > grt.pdf
*/
