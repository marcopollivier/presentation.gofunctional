//go:build go1.22

package main

import (
	"fmt"
	"testing"
)

// Mesmo pacote, mesmo laço, LINGUAGEM 1.22 — fixada pelo `//go:build go1.22`.
// É o "depois".
//
// Se você rodar os testes num toolchain anterior ao 1.22, este arquivo é
// simplesmente excluído do build e só o "antes" roda. Nada quebra.

// closuresModernas reproduz exatamente o mesmo laço do arquivo do "antes".
func closuresModernas() []func() []int {
	var fs []func() []int
	var vistos []int
	for i := 0; i < 3; i++ {
		fs = append(fs, func() []int {
			vistos = append(vistos, i)
			return vistos
		})
	}
	return fs
}

func TestLinguagem122TemUmIPorIteracao(t *testing.T) {
	var got []int
	for _, f := range closuresModernas() {
		got = f()
	}

	quer := "[0 1 2]"
	if fmt.Sprint(got) != quer {
		t.Errorf("com linguagem 1.22 as closures viram %v, esperado %s", got, quer)
	}
}
