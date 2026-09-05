package crossversion

import (
	"fmt"
	"testing"
)

// Estes dois testes valem mais que os benchmarks.
//
// Eles demonstram o erro de 2023 SEM cronômetro, sem variância, sem "depende da
// máquina". São determinísticos: rodam igual em qualquer versão do Go e em
// qualquer hardware. Um benchmark convence quem confia no seu setup; um teste
// que imprime o estado antes e depois convence quem não confia.
//
// Rode com -v para ver a saída no palco:
//
//	go test -v -run 'TestLegacy|TestImmutable' ./...

const preview = 5 // quantos elementos mostrar; o suficiente para a plateia ler

// TestLegacyBenchmarkIsBiased mostra a causa do erro num único passo.
//
// Uma chamada a SortInPlace basta: o slice do chamador sai ordenado. No loop do
// benchmark de 2023, isso significa que a iteração 1 ordenava dados aleatórios
// e as iterações 2..N ordenavam dados JÁ ORDENADOS — o melhor caso do pdqsort.
// O lado imutável, que nunca escrevia no original, recebia dados aleatórios
// sempre.
//
// Não foi descuido de medição. Foi não-localidade: `func SortInPlace(s []int)`
// não tem `*`, não tem anotação, não tem nada que avise que o argumento será
// destruído. O erro estava na assinatura, não no laço.
func TestLegacyBenchmarkIsBiased(t *testing.T) {
	data := makeData(1000)

	before := fmt.Sprint(data[:preview])
	sortedBefore := isSorted(data)

	// Exatamente uma iteração do benchmark de 2023.
	SortInPlace(data)

	after := fmt.Sprint(data[:preview])
	sortedAfter := isSorted(data)

	t.Logf("iteração 1 recebeu: %s...  (ordenado? %v)", before, sortedBefore)
	t.Logf("iteração 2 receberá: %s...  (ordenado? %v)", after, sortedAfter)

	if sortedBefore {
		t.Fatal("os dados de partida já estavam ordenados: o teste não prova nada")
	}
	if !sortedAfter {
		t.Fatal("SortInPlace não ordenou o slice do chamador")
	}

	// A afirmação que interessa: o chamador foi reescrito. É isso que
	// contaminou todas as iterações seguintes.
	t.Log("a partir daqui, todas as iterações do benchmark de 2023 mediam o melhor caso")
}

// TestImmutablePreservesInput é o contraste: a versão pura não CONSEGUE
// cometer esse erro.
//
// Não porque alguém tomou cuidado, mas porque a entrada é inalcançável para
// escrita — SortImmutable trabalha sobre memória que ela mesma alocou. O
// benchmark permanece honesto por construção, não por disciplina de quem o
// escreveu.
func TestImmutablePreservesInput(t *testing.T) {
	data := makeData(1000)
	before := fmt.Sprint(data[:preview])

	// Cem iterações, o quanto quiser: a entrada não muda.
	for i := 0; i < 100; i++ {
		got := SortImmutable(data)
		if !isSorted(got) {
			t.Fatalf("iteração %d: o resultado não saiu ordenado", i)
		}
	}

	after := fmt.Sprint(data[:preview])

	t.Logf("antes de 100 chamadas: %s...", before)
	t.Logf("depois de 100 chamadas: %s...", after)

	if before != after {
		t.Fatalf("a entrada mudou: %s -> %s", before, after)
	}
	if isSorted(data) {
		t.Fatal("a entrada saiu ordenada — a mutação escapou")
	}

	t.Log("toda iteração recebe a mesma qualidade de entrada: o viés é impossível aqui")
}

// isSorted evita `slices.IsSorted` (1.21+) e `sort.IntsAreSorted` por clareza:
// o laço explícito é o que se lê em voz alta no palco.
func isSorted(s []int) bool {
	for i := 1; i < len(s); i++ {
		if s[i-1] > s[i] {
			return false
		}
	}
	return true
}
