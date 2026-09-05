package main

import "fmt"

// Lab 030 — a variável de loop. Bloco 3, e o primeiro item da lista do
// "o que mudou".
//
// Este é o footgun de doze anos que sumiu sem ninguém precisar reescrever
// código. Até o Go 1.21, o `i` do `for` era UMA variável, reusada a cada volta:
// as três closures capturavam a mesma caixa, e quando alguém finalmente as
// chamava, a caixa já valia 3. Do 1.22 em diante cada iteração tem o seu `i`.
//
//	linguagem 1.21  ->  3 3 3
//	linguagem 1.22  ->  0 1 2
//
// O fonte abaixo é o MESMO nos dois casos. Não há flag no código, não há
// `i := i`, não há nada a refatorar: mudou o significado do `for`.
//
// O que decide qual dos dois você vê NÃO é o toolchain instalado — é a versão
// de LINGUAGEM, que sai da linha `go` do go.mod e pode ser trocada por chamada
// com `-gcflags=-lang`. Trocar só o GOTOOLCHAIN não muda nada, e um toolchain
// antigo nem compila um módulo que declare versão maior que ele. `make
// comparar` põe os dois lados na tela com o mesmo compilador.
//
// O go.mod deste lab declara `go 1.21` de propósito: `make run` mostra o
// "antes". Aquela linha é o botão do experimento.
//
// Roda com: make comparar — ver o README.

func main() {
	// Antes: as três closures capturam a MESMA variável → imprime 3 3 3
	// As closures são só GUARDADAS aqui. Nada é impresso ainda — e é esse
	// adiamento que revela o problema: quando elas finalmente rodam, lá
	// embaixo, o laço já terminou. Com uma variável só, ela vale 3.
	var fs []func()
	for i := 0; i < 3; i++ {
		fs = append(fs, func() { fmt.Println(i) })
	}
	// Go 1.22+: cada iteração tem seu i → imprime 0 1 2
	for _, f := range fs {
		f()
	}
}
