package ponteiros

import (
	"fmt"
	"slices"
	"sort"
)

// Bloco "o paradigma": por que "evitar mudança de estado" não é preciosismo.
//
// Um slice é um cabeçalho (ponteiro para o array, len, cap) passado POR VALOR.
// A cópia do cabeçalho é barata e inofensiva — mas ela aponta para o MESMO
// array. Então quem recebe um []int recebe, de graça, permissão de escrita na
// memória do chamador. A assinatura `func(x []int)` não avisa nada disso.
//
// SortedPrint e SafeSortedPrint têm a mesma assinatura, imprimem a mesma coisa
// e diferem só no que sobra depois. É o contraste do slide.

// SortedPrint ordena e imprime. O nome fala de imprimir; o efeito colateral de
// mutar o argumento não aparece na assinatura — esse é o ponto.
func SortedPrint(x []int) {
	sort.Ints(x)
	fmt.Println(x)
}

// SafeSortedPrint faz o mesmo trabalho visível sem tocar no slice do chamador.
// Note que ela TAMBÉM muta: pureza não é "nunca mutar", é a mutação não
// escapar. `c` nasce e morre aqui dentro.
func SafeSortedPrint(x []int) {
	c := slices.Clone(x)
	slices.Sort(c) // muta! mas ninguém além dela vê o `c`
	fmt.Println(c)
}

func Print(x []int) {
	fmt.Println(x)
}
