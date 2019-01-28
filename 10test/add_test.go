package add

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct{ a, b, r int }{
		{1, 2, 3},
		{10, 10, 20},
		{11, 12, 23},
		{1000, 2000, 3000},
	}
	for _, tt := range tests {
		if a := Add(tt.a, tt.b); a != tt.r {
			t.Errorf("add(%d,%d),got %d,expected %d", tt.a, tt.b, a, tt.r)
		}
	}
}
func BenchmarkAdd(t *testing.B) {
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Add(1, 2)
	}
}
