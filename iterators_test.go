package gofunctional

import "testing"

func TestSumSquaresOfEvens(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	// pares: 2,4,6 -> quadrados: 4,16,36 -> soma: 56
	const want = 56

	if got := SumSquaresOfEvens(nums); got != want {
		t.Errorf("SumSquaresOfEvens = %d; want %d", got, want)
	}
	// As duas implementações têm que concordar — é o mesmo cálculo.
	if got := SumSquaresOfEvensSlices(nums); got != want {
		t.Errorf("SumSquaresOfEvensSlices = %d; want %d", got, want)
	}
}

// Os benchmarks mostram o ponto: a versão lazy não aloca slices intermediários
// (alocações constantes, dos closures do pipeline), enquanto a versão que
// materializa aloca proporcionalmente ao tamanho da entrada.
//   go test -bench=SumSquares -benchmem ./...

func benchInput(n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = i
	}
	return s
}

func BenchmarkSumSquaresLazy(b *testing.B) {
	nums := benchInput(10_000)
	var sink int
	for b.Loop() {
		sink = SumSquaresOfEvens(nums)
	}
	_ = sink
}

func BenchmarkSumSquaresSlices(b *testing.B) {
	nums := benchInput(10_000)
	var sink int
	for b.Loop() {
		sink = SumSquaresOfEvensSlices(nums)
	}
	_ = sink
}
