package main

import "fmt"

// Lab 02 — ponteiro e função pura: o ponteiro não é o vilão; o que a função
// faz com ele é. Bloco 2.
//
// Len e Changer recebem exatamente a mesma coisa, um *string. A forma da
// assinatura é idêntica; o que muda é o verbo do corpo — uma lê, a outra
// escreve. Nada no `func(*string)` distingue as duas.
//
// Os três momentos existem para separar o que a gente costuma confundir num
// ponteiro: a CAIXA e o CONTEÚDO.
//
//	Momento 1  a variável nasce. O %p impresso é o endereço da caixa `nome`.
//	Momento 2  Len(&nome) lê e devolve 14. Caixa e conteúdo intactos.
//	Momento 3  Changer(&nome) escreve. O endereço continua O MESMO — a caixa
//	           é a mesma — e o conteúdo virou outro.
//
// O endereço estável nos três momentos é o ponto do slide: passar ponteiro não
// é mover a variável, é entregar a chave da caixa. Quem tem a chave pode só
// olhar ou pode trocar o que está dentro, e o chamador não tem como saber qual
// dos dois vai acontecer sem abrir o corpo da função.
//
// Roda com: make run — no Go 1.21, ver o README.

func main() {
	// Momento 1: Declaração da variável e alocação de memória
	var nome = "Marco Ollivier"
	fmt.Printf("Momento 1: Endereco [%p] Nome [%s]\n", &nome, nome)

	// Momento 2: uma referencia é passada
	var nomeLen = Len(&nome)
	fmt.Printf("Momento 2: Endereco [%p] Nome [%s] Tamanho [%d]\n", &nome, nome, nomeLen)

	// Moment 3: uma função que altera a referência é chamada
	Changer(&nome)
	nomeLen = Len(&nome)
	fmt.Printf("Momento 2: Endereco [%p] Nome [%s] Tamanho [%d]\n", &nome, nome, nomeLen)

}

// Changer recebe o mesmo *string que Len e faz o oposto: escreve na variável
// do chamador. É o outro uso da mesma chave — e o que torna o ponteiro do Len
// caro, mesmo que o Len não faça nada de errado.
func Changer(s *string) {
	var novoNome = "Marco Paulo Ollivier"
	*s = novoNome
}

// Aqui temos um ponteiro explicito
// Ainda assim temos uma funcao pura
//
// Pura no sentido que importa: não escreve em nada, não guarda estado, não faz
// I/O. Ler `*s` é leitura e mais nada — e o `Changer` acima mostra que a mesma
// assinatura permitiria escrever.
//
// O preço do ponteiro aparece no Momento 3: `Len(&nome)` devolve 14 antes e 20
// depois. MESMO argumento, resposta diferente. A função é pura, mas a expressão
// deixou de ser substituível pelo seu resultado — a transparência referencial
// do lab 01 se perde não pelo corpo da função, e sim pelo que o argumento
// passou a significar. Com `Len(nome)`, por valor, isso não teria como ocorrer.
func Len(s *string) int {
	return len(*s)
}
