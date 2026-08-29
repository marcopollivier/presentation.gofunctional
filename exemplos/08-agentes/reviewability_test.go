package agentes

import (
	"slices"
	"testing"
)

// A função pura testa-se sem cerimônia: table-driven, sem setup, sem mock.
// É o loop de feedback que um agente consegue fechar sozinho.
func TestMovingAveragePure(t *testing.T) {
	tests := []struct {
		name   string
		series []float64
		window int
		want   []float64
	}{
		{"janela 3", []float64{1, 2, 3, 4, 5}, 3, []float64{2, 3, 4}},
		{"janela = tamanho", []float64{2, 4, 6}, 3, []float64{4}},
		{"janela grande demais", []float64{1, 2}, 3, nil},
		{"janela inválida", []float64{1, 2, 3}, 0, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MovingAveragePure(tt.series, tt.window); !slices.Equal(got, tt.want) {
				t.Errorf("MovingAveragePure(%v, %d) = %v; want %v", tt.series, tt.window, got, tt.want)
			}
		})
	}
}

// A versão com estado só está "correta" se você garantir a ordem e o reset do
// global. Este teste passa hoje, mas a corretude é frágil: depende do
// accumulator estar zerado no lugar certo. É essa fragilidade que a palestra
// mostra no palco.
func TestMovingAverageStateful(t *testing.T) {
	got := MovingAverageStateful([]float64{1, 2, 3, 4, 5}, 3)
	if want := []float64{2, 3, 4}; !slices.Equal(got, want) {
		t.Errorf("MovingAverageStateful = %v; want %v", got, want)
	}
}
