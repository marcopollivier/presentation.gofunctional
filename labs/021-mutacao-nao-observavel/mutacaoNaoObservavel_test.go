package main

import (
	"slices"
	"testing"
)

// Os testes escrevem o que a tela mostra — e um deles mostra o que a tela não
// consegue: que a permissão de escrita alcança o array inteiro de origem, não
// só o pedaço que foi passado.

// O ponto do slide: SortedPrint reordena o slice do chamador.
func TestSortedPrintMutaOChamador(t *testing.T) {
	x := []int{9, 1, 3}

	SortedPrint(x)

	if !slices.Equal(x, []int{1, 3, 9}) {
		t.Errorf("x = %v depois de SortedPrint, esperado [1 3 9] — a mutação É o ponto", x)
	}
}

// A mutação existe, mas morre dentro da função. Do lado de fora, nada aconteceu.
func TestSafeSortedPrintNaoEhObservavel(t *testing.T) {
	x := []int{9, 1, 3}

	SafeSortedPrint(x)

	if !slices.Equal(x, []int{9, 1, 3}) {
		t.Errorf("x = %v depois de SafeSortedPrint, esperado [9 1 3] intacto", x)
	}
}

// A ordem das duas demos no main não é arbitrária: depois de SortedPrint o
// slice já está ordenado, e aí SafeSortedPrint pareceria fazer a mesma coisa.
// Este teste fixa a armadilha para ninguém "arrumar" o main sem perceber.
func TestOContrasteSoApareceComEntradaDesordenada(t *testing.T) {
	jaOrdenado := []int{1, 3, 9}

	SortedPrint(jaOrdenado)

	if !slices.Equal(jaOrdenado, []int{1, 3, 9}) {
		t.Errorf("x = %v, esperado [1 3 9]", jaOrdenado)
	}
	// Com entrada já ordenada, a função que muta e a que não muta deixam o
	// mesmo resultado: o contraste depende do cenário, não só do código.
}

// O alcance da permissão: passar um pedaço entrega o array de trás inteiro.
// `fatia` tem 3 elementos, mas escreve dentro do array de 4 que o chamador
// nunca ofereceu por completo.
func TestAMutacaoAlcancaOArrayDeOrigem(t *testing.T) {
	original := []int{9, 1, 3, 0}
	fatia := original[:3]

	SortedPrint(fatia)

	if !slices.Equal(original, []int{1, 3, 9, 0}) {
		t.Errorf("original = %v, esperado [1 3 9 0] — a escrita atravessou a fatia", original)
	}
}
