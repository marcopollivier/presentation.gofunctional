package main

import (
	"fmt"
	"sort"
)

// Lab 021 — a mutação que não se observa. Bloco 2.
//
// O par do 020, agora com slice no lugar do ponteiro — e a diferença é a tese
// do capítulo: aqui não existe `&` nem `*` em lugar nenhum. Nem na assinatura,
// nem na chamada. E mesmo assim SortedPrint reescreve o slice de quem chamou.
//
// Um slice é um cabeçalho (ponteiro para o array, len, cap) passado POR VALOR.
// A cópia do cabeçalho é barata e inofensiva, mas aponta para o MESMO array.
// Quem recebe um []int recebe, de graça, permissão de escrita na memória do
// chamador. No 020 o `&nome` no ponto de chamada pelo menos avisava; aqui não
// há aviso nenhum, e é por isso que este é o caso perigoso.
//
// O nome do arquivo é a lição: SafeSortedPrint TAMBÉM muta — ela ordena `c`.
// A diferença não é mutar ou não mutar, é a mutação ser observável de fora.
// `c` nasce e morre dentro da função; ninguém no universo consegue notar que
// ela existiu. Pureza não é abstinência de mutação, é a mutação não escapar.
//
// As duas funções têm a MESMA assinatura, func([]int). A diferença só aparece
// na linha "depois".
//
// SortedPrint é anti-exemplo — o lado esquerdo do slide. Não a torne segura.
//
// Roda com: make run — no Go 1.21, ver o README.

func main() {
	var x = []int{9, 1, 3}

	// A ordem importa: a demo imutável vem primeiro justamente porque a segunda
	// destrói o cenário. Depois de SortedPrint, `x` já está ordenado: invertendo
	// a ordem, as quatro linhas sairiam iguais e o contraste sumiria.
	fmt.Printf("Imutavel: x antes da chamada  %v\n", x)
	SafeSortedPrint(x)
	fmt.Printf("Imutavel: x depois da chamada %v\n", x)

	fmt.Printf("Mutavel: x antes da chamada   %v\n", x)
	SortedPrint(x)
	fmt.Printf("Mutavel: x depois da chamada  %v\n", x)

}

// SortedPrint ordena o slice recebido — e o slice de quem chamou junto, porque
// são o mesmo array. Nada na assinatura `func([]int)` avisa isso: nem um
// caractere. É a mutação observável, e é o problema invisível do slide.
func SortedPrint(x []int) {
	sort.Ints(x)
}

// SafeSortedPrint faz o mesmo trabalho sem tocar no slice do chamador: copia
// para um array novo e ordena a cópia.
//
// Ela muta — `sort.Ints(c)` reescreve `c` inteiro. Mas `c` foi alocado aqui,
// não sai daqui e morre no return. A mutação existe e não é observável, que é
// o ponto do lab.
//
// O `make` + `copy` é a versão explícita, que mostra o que está acontecendo.
// No 1.21 dava para escrever `slices.Clone(x)`, uma linha; a partir do 1.23 o
// trabalho inteiro cabe em `slices.Sorted(slices.Values(x))`. A evolução dessa
// mesma função ao longo das releases é assunto do bloco 3.
func SafeSortedPrint(x []int) {
	c := make([]int, len(x))
	copy(c, x)

	sort.Ints(c)
}
