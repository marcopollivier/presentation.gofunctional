package main

import (
	"fmt"
	"slices"
	"testing"
)

// O par do 021: lá os testes provavam que a mutação escapava; aqui provam que
// não há por onde ela escapar.

// O ponto do slide: o argumento sai da função como entrou.
func TestSortedPrintNaoMutaOChamador(t *testing.T) {
	x := []int{9, 1, 3}

	SortedPrint(x)

	if !slices.Equal(x, []int{9, 1, 3}) {
		t.Errorf("x = %v depois de SortedPrint, esperado [9 1 3] intacto", x)
	}
}

// E não é sorte com uma entrada específica.
func TestSortedPrintNaoMutaNenhumaEntrada(t *testing.T) {
	casos := [][]int{
		{},
		{1},
		{9, 1, 3},
		{5, 5, 5},
		{-2, 10, 0, -7},
	}

	for _, c := range casos {
		original := slices.Clone(c)

		SortedPrint(c)

		if !slices.Equal(c, original) {
			t.Errorf("entrada %v virou %v depois de SortedPrint", original, c)
		}
	}
}

// A versão de uma linha entrega o mesmo resultado que as três linhas do 021.
// É a comparação que o slide faz: mudou a ergonomia, não a semântica.
func TestOResultadoEhOMesmoDoCaminhoAntigo(t *testing.T) {
	x := []int{9, 1, 3}

	moderno := slices.Sorted(slices.Values(x))

	antigo := make([]int, len(x))
	copy(antigo, x)
	slices.Sort(antigo)

	if !slices.Equal(moderno, antigo) {
		t.Errorf("moderno = %v, antigo = %v — deveriam coincidir", moderno, antigo)
	}
	if !slices.Equal(x, []int{9, 1, 3}) {
		t.Errorf("x = %v, esperado intacto pelos dois caminhos", x)
	}
}

// slices.Sorted devolve um slice NOVO: escrever nele não alcança a origem.
// É o oposto exato do TestAMutacaoAlcancaOArrayDeOrigem, no 021.
func TestOResultadoNaoCompartilhaArrayComAEntrada(t *testing.T) {
	x := []int{9, 1, 3}

	c := slices.Sorted(slices.Values(x))
	c[0] = 999

	if slices.Contains(x, 999) {
		t.Errorf("x = %v — a escrita no resultado alcançou a entrada", x)
	}
}

// Example é a demo de palco, com a saída verificada.
//
// A segunda linha sai de DENTRO de SortedPrint e mostra o resultado ordenado;
// a terceira é o main mostrando que `x` não mudou. As duas carregam o mesmo
// rótulo — se o texto do Printf mudar, este Example falha e é só atualizar.
func Example() {
	x := []int{9, 1, 3}

	fmt.Printf("x antes da chamada   %v\n", x)
	SortedPrint(x)
	fmt.Printf("x depois da chamada  %v\n", x)

	// Output:
	// x antes da chamada   [9 1 3]
	// x depois da chamada  [1 3 9]
	// x depois da chamada  [9 1 3]
}
