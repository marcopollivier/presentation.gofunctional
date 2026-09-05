package main

import (
	"fmt"
	"slices"
	"sort"
)

// Lab 02 — estado compartilhado e ponteiros: o `*` que você vê não é o que
// machuca. Bloco 2.
//
// Dois jeitos de entregar permissão de escrita na memória do chamador:
//
//	Len(&nome)       explícito — o `&` no ponto de chamada avisa. E o ponteiro
//	                 nem era necessário: a função só lê.
//	SortedPrint(xs)  invisível — nenhum `&`, nenhum `*`, e mesmo assim a função
//	                 reordena o slice de quem chamou.
//
// O segundo caso é o que morde. Um slice é um cabeçalho (ponteiro para o
// array, len, cap) passado POR VALOR: a cópia do cabeçalho é barata e
// inofensiva, mas aponta para o MESMO array. `func(x []int)` não tem um único
// caractere que faça você desconfiar disso.
//
// Este lab roda no Go 1.21 de propósito — o Go de quando a palestra nasceu.
// `slices.Clone` e `slices.Sort` já existem (o pacote `slices` entrou na
// stdlib no 1.21), mas o `slices.Sorted(slices.Values(x))` de hoje ainda não:
// isso é 1.23. Ver `make version` para conferir o toolchain em uso.
//
// Len e SortedPrint são anti-exemplos: são o lado esquerdo do slide. Não os
// "conserte" — o par com LenValue e SortedPrintSafe é o conteúdo.

// Len pede um *string para fazer uma coisa que não precisa de ponteiro nenhum:
// ler o tamanho. Pedir ponteiro é pedir permissão de escrita; esta função
// pediu, não usa, e ainda obriga o chamador a ter uma variável endereçável.
// De brinde, ganha um estado que a versão por valor não tem: s == nil, que
// aqui vira panic.
func Len(s *string) int {
	return len(*s)
}

// LenValue faz o mesmo trabalho e não pede nada além do que lê. A cópia de uma
// string é o cabeçalho (ponteiro + len), 16 bytes — não copia os bytes do
// texto. O "ponteiro por performance" não se sustenta aqui.
func LenValue(s string) int {
	return len(s)
}

// SortedPrint ordena e imprime. O nome fala de imprimir; a reordenação do
// slice do chamador não aparece em lugar nenhum da assinatura — esse é o
// problema invisível.
//
// sort.Ints não é deprecado: desde o Go 1.22 ele só chama slices.Sort. Fica
// aqui porque é o que se escrevia em 2023 e porque a mutação é a mesma.
func SortedPrint(x []int) {
	sort.Ints(x)
	fmt.Println(x)
}

// SortedPrintSafe entrega a mesma saída visível sem tocar no slice de quem
// chamou. Repare que ela TAMBÉM muta: pureza não é "nunca mutar", é a mutação
// não escapar. O `c` nasce e morre aqui dentro.
func SortedPrintSafe(x []int) {
	c := slices.Clone(x)
	slices.Sort(c)
	fmt.Println(c)
}

func main() {
	fmt.Println("1) o ponteiro explícito — e desnecessário")
	nome := "Marco Ollivier"
	fmt.Printf("   Len(&nome)     = %d   <- pede permissão de escrita para só ler\n", Len(&nome))
	fmt.Printf("   LenValue(nome) = %d   <- mesmo número, sem pedir nada\n", LenValue(nome))
	fmt.Println()

	fmt.Println("2) o compartilhamento invisível — sem um único `*` à vista")

	// O mesmo slice de partida nos dois casos, para o contraste ser honesto.
	original := []int{3, 1, 2}

	seguro := slices.Clone(original)
	fmt.Printf("   antes:                 %v\n", seguro)
	fmt.Print("   SortedPrintSafe(x) ->  ")
	SortedPrintSafe(seguro)
	fmt.Printf("   depois:                %v   <- intacto\n", seguro)

	exposto := slices.Clone(original)
	fmt.Printf("   antes:                 %v\n", exposto)
	fmt.Print("   SortedPrint(x)     ->  ")
	SortedPrint(exposto)
	fmt.Printf("   depois:                %v   <- o chamador foi reescrito\n", exposto)
}
