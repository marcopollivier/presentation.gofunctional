//go:build go1.21

package main

import (
	"fmt"
	"testing"
)

// Este arquivo é compilado com a LINGUAGEM 1.21, garantido pelo `//go:build
// go1.21` da primeira linha — independente do que o go.mod diga. É o "antes",
// fixado por escrito.
//
// O par com loopvar_depois_test.go é o artefato mais direto deste lab: os dois
// comportamentos coexistem no MESMO pacote, na MESMA compilação, e `go test`
// prova os dois de uma vez. A diferença entre eles é uma linha de build tag.

// closures reproduz exatamente o laço do main.
func closures() []func() []int {
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

func TestLinguagem121CompartilhaAVariavel(t *testing.T) {
	var got []int
	for _, f := range closures() {
		got = f()
	}

	quer := "[3 3 3]"
	if fmt.Sprint(got) != quer {
		t.Errorf("com linguagem 1.21 as closures viram %v, esperado %s — o footgun É o ponto", got, quer)
	}
}
