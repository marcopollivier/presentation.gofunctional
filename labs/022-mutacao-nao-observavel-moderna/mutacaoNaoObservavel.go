package main

import (
	"fmt"
	"slices"
)

// Lab 022 — a mesma mutação não observável, na versão moderna. Bloco 2, com um
// pé no 3.
//
// É o "depois" do 021. Lá, guardar a promessa de não mexer no slice do chamador
// custava três linhas e uma variável temporária:
//
//	c := make([]int, len(x))
//	copy(c, x)
//	sort.Ints(c)
//
// Aqui custa uma, e a stdlib faz o clone por você:
//
//	c := slices.Sorted(slices.Values(x))
//
// Esse é o ponto do slide, e ele desmonta o argumento de 2023. Naquela época a
// defesa da mutação era ergonômica: `sort.Ints(x)` era uma linha, o caminho
// imutável era um parágrafo. O 1.21 trouxe `slices` para a stdlib, o 1.23
// trouxe os iteradores, e a versão que NÃO muta virou a mais curta das duas.
// O atrito que sustentava "funcional em Go é forçação de barra" caiu — e caiu
// sem release note nenhuma dizendo a palavra "funcional".
//
// Repare no nome: SortedPrint é o MESMO nome do 021, com a MESMA assinatura,
// func([]int) — e comportamento oposto. Lá reescrevia o slice do chamador,
// aqui não toca nele. Nem o nome nem a assinatura contam a verdade; só o
// corpo. É a tese do capítulo levada ao limite.
//
// Roda com: make run — no Go 1.27, ver o README.

func main() {
	var x = []int{9, 1, 3}

	fmt.Printf("x antes da chamada   %v\n", x)
	SortedPrint(x)
	fmt.Printf("x depois da chamada  %v\n", x)
}

// SortedPrint ordena e imprime sem tocar no argumento.
//
// `slices.Values(x)` transforma o slice num iterador (`iter.Seq`), e
// `slices.Sorted` consome esse iterador para um slice NOVO, já ordenado. O
// original nunca é escrito — não porque alguém tomou cuidado, mas porque não
// existe caminho para escrever nele nesta formulação.
//
// A mutação continua existindo: `slices.Sorted` ordena o slice que ela mesma
// alocou. Continua não sendo observável, que é o assunto do 021.
//
// Ela também aloca, e essa alocação é exatamente o que o benchmark do bloco 4
// mede. Guarde: em 2023 esse custo foi lido como "o preço do paradigma", e o
// post-mortem em exemplos/07-benchmark-sort mostra que a conta estava errada.
func SortedPrint(x []int) {
	c := slices.Sorted(slices.Values(x))
	fmt.Printf("x depois da chamada  %v\n", c)
}
