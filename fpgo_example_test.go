package gofunctional

import "testing"

func TestProcessFpGo(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{"válido e >= 10", "21", 42},   // 21 passa, dobra -> 42
		{"positivo mas < 10", "7", -1}, // Chain falha -> default
		{"não positivo", "0", -1},      // parse falha na validação -> default
		{"não é número", "abc", -1},    // parse falha no Atoi -> default
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessFpGo(tt.in); got != tt.want {
				t.Errorf("ProcessFpGo(%q) = %d; want %d", tt.in, got, tt.want)
			}
		})
	}
}
