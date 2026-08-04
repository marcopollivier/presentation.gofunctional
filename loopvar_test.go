package gofunctional

import "testing"

func call(fns []func() int) []int {
	out := make([]int, len(fns))
	for i, fn := range fns {
		out[i] = fn()
	}
	return out
}

func TestLoopVarPerIteration(t *testing.T) {
	// Go 1.22+: cada closure captura sua própria cópia da variável de loop.
	got := call(LoopVarPerIteration())
	want := []int{0, 1, 2}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("LoopVarPerIteration() = %v; want %v (rodando em Go < 1.22?)", got, want)
		}
	}
}

func TestLoopVarShared(t *testing.T) {
	// O comportamento pré-1.22 reproduzido de propósito: variável compartilhada.
	got := call(LoopVarShared())
	want := []int{3, 3, 3}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("LoopVarShared() = %v; want %v", got, want)
		}
	}
}
