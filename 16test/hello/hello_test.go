package hello

import (
	"bytes"
	"testing"
	"text/template"
)

// go test -bench BenchmarkHello

/*
BenchmarkHello-4   	20000000	        72.5 ns/op	       5 B/op	       1 allocs/op
PASS
ok  	github.com/gonote/16test/hello	1.749s
结果表示以 72.5 ns/op 的速度运行了 20000000 次，耗时 1.749s 秒

多跑几组
go test -bench BenchmarkHello -cpu 1,2,4 -count=5
*/
func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Hello()
	}
}

// 想并行地运行测试用例，我们可以用 b.RunParallel

/*
go test -bench=BenchmarkTemplateParallel -cpu 1,2,4
*/
func BenchmarkTemplateParallel(b *testing.B) {
	templ := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			buf.Reset()
			templ.Execute(&buf, "World")
		}
	})
}

// go test -run TestGroupedParallel -v
func TestGroupedParallel(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{
			Name: "A=1",
		},
		{
			Name: "A=2",
		},
		{
			Name: "B=1",
		},
	}
	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
		})
	}
}
