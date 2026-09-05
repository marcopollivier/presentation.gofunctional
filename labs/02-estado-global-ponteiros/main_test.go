package main

import (
	"slices"
	"testing"
)

// Os testes provam o que a assinatura esconde.
//
// Repare que o teste interessante não é sobre o RETORNO de SortedPrint — é
// sobre o que sobrou no chamador depois dela. Quando você precisa escrever uma
// asserção sobre o estado de quem chamou, a função já contou que vaza.

func TestLenEquivaleALenValue(t *testing.T) {
	casos := []string{"", "Go", "Marco Ollivier", "áéíóú"}

	for _, s := range casos {
		// A cópia não é supersticão: no Go 1.21 a variável de loop é UMA só,
		// reusada a cada volta, e logo abaixo se tira o endereço dela. Do 1.22
		// em diante esta linha vira ruído. É o lab inteiro em miniatura.
		s := s
		if got, quer := Len(&s), LenValue(s); got != quer {
			t.Errorf("Len(&%q) = %d, LenValue(%q) = %d — deveriam ser iguais", s, got, s, quer)
		}
	}
}

// O estado a mais que só a versão com ponteiro tem: o nil.
// LenValue não consegue nem chegar aqui — não existe string nil.
func TestLenComNilEntraEmPanico(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Len(nil) deveria entrar em pânico — é o estado extra que o ponteiro traz")
		}
	}()

	Len(nil)
}

// O ponto do slide: SortedPrint reordena o slice de quem chamou.
func TestSortedPrintMutaOChamador(t *testing.T) {
	x := []int{3, 1, 2}

	SortedPrint(x)

	if !slices.Equal(x, []int{1, 2, 3}) {
		t.Errorf("x = %v depois de SortedPrint, esperado [1 2 3] — a mutação É o ponto", x)
	}
}

// Mesma saída visível, chamador intacto.
func TestSortedPrintSafeNaoMutaOChamador(t *testing.T) {
	x := []int{3, 1, 2}

	SortedPrintSafe(x)

	if !slices.Equal(x, []int{3, 1, 2}) {
		t.Errorf("x = %v depois de SortedPrintSafe, esperado [3 1 2] intacto", x)
	}
}

// Example é a demo de palco: as duas funções imprimem exatamente a mesma
// linha. A diferença não está na saída — está no que sobra.
func Example() {
	a := []int{3, 1, 2}
	b := []int{3, 1, 2}

	SortedPrintSafe(a)
	SortedPrint(b)

	// Output:
	// [1 2 3]
	// [1 2 3]
}
