package gofunctional

import "testing"

func TestHalfIfEven(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{10, 5},  // par -> metade
		{7, -1},  // ímpar -> Chain falha -> default
		{4, 2},   // par -> metade
	}
	for _, tt := range tests {
		if got := HalfIfEven(tt.in); got != tt.want {
			t.Errorf("HalfIfEven(%d) = %d; want %d", tt.in, got, tt.want)
		}
	}
}

func TestOptionMapChain(t *testing.T) {
	// None curto-circuita: Map sobre None continua None.
	got := MapOption(None[int](), func(n int) int { return n * 2 })
	if GetOrElse(got, -1) != -1 {
		t.Error("MapOption sobre None deveria continuar None")
	}
	// Some passa pelo Map.
	if v := GetOrElse(MapOption(Some(21), func(n int) int { return n * 2 }), -1); v != 42 {
		t.Errorf("MapOption(Some(21), *2) = %d; want 42", v)
	}
}
