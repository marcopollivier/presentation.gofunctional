package ponteiros

import (
	"slices"
	"testing"
)

// Example é a demo de palco: a sequência que antes vivia num main(), agora
// com a saída verificada pelo `go test`. Os dois blocos são espelhados — só a
// terceira linha de cada um revela a diferença.
func Example() {
	x := []int{9, 1, 3}

	Print(x)
	SafeSortedPrint(x)
	Print(x) // o original sobreviveu

	SortedPrint(x)
	Print(x) // o original mudou

	// Output:
	// [9 1 3]
	// [1 3 9]
	// [9 1 3]
	// [1 3 9]
	// [1 3 9]
}

// O contrato de cada uma, exercitado sem cerimônia: o que interessa não é o
// que elas imprimem, é o que sobra no slice do chamador.
func TestContratoDeMutacao(t *testing.T) {
	tests := []struct {
		name string
		fn   func([]int)
		muta bool
	}{
		{"SortedPrint muta o argumento", SortedPrint, true},
		{"SafeSortedPrint preserva o argumento", SafeSortedPrint, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := []int{9, 1, 3}
			x := slices.Clone(original)

			tt.fn(x)

			if mutou := !slices.Equal(x, original); mutou != tt.muta {
				t.Errorf("mutou o argumento = %v; want %v (got %v, original %v)",
					mutou, tt.muta, x, original)
			}
		})
	}
}
