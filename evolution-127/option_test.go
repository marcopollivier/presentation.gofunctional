package evolution

import "testing"

func TestHalfIfEven(t *testing.T) {
	tests := []struct {
		in   int
		want int
	}{
		{10, 5}, // par -> metade
		{7, -1}, // ímpar -> Chain falha -> default
		{4, 2},  // par -> metade
	}
	for _, tt := range tests {
		if got := HalfIfEven(tt.in); got != tt.want {
			t.Errorf("HalfIfEven(%d) = %d; want %d", tt.in, got, tt.want)
		}
	}
}

// A prova de que o encadeamento com métodos genéricos funciona: Map muda o
// tipo (int -> string) e ainda encadeia GetOrElse.
func TestChainableMapChangesType(t *testing.T) {
	got := Some(21).
		Map(func(n int) int { return n * 2 }).
		Map(func(n int) string {
			if n == 42 {
				return "resposta"
			}
			return "outro"
		}).
		GetOrElse("vazio")
	if got != "resposta" {
		t.Errorf("cadeia Map/Map/GetOrElse = %q; want %q", got, "resposta")
	}
}
