package main

import "fmt"

// Lab 01 — "vamos relembrar funcional": a função matemática e a função Go,
// lado a lado. Bloco 1.
//
// O slide é f(x) = x² + 3, e f(2) = 7. Na matemática isso não é "devolve 7
// desta vez": f(2) É 7. Em qualquer canto do papel, em qualquer ordem, quantas
// vezes você olhar. Dá para riscar f(2) e escrever 7 no lugar sem mudar o
// resto da conta — isso tem nome, transparência referencial, e é de onde vem
// toda a conversa deste bloco.
//
// f() abaixo mantém a promessa. badF() calcula exatamente a mesma coisa,
// devolve exatamente o mesmo número — e mesmo assim quebra a promessa, porque
// escreve em `value`. O detalhe que faz o slide: as duas têm a MESMA
// assinatura, func(float64) float64. A caixa do "Contrato" no diagrama não
// distingue uma da outra. Você tem que abrir o corpo para saber.
//
// badF é anti-exemplo, como manda o CLAUDE.md — é o lado direito do slide.
// Não a torne pura.
//
// Roda com: make run

// value é o estado global que badF suja. Existe só para ser o crime.
var value float64

// f é a tradução direta de f(x) = x² + 3.
//
// Entrada, processamento, saída — e nada mais. O `x` do parâmetro é o mesmo
// `x` do slide, de propósito: a assinatura em Go é o contrato que o diagrama
// desenha em volta de f(x).
func f(x float64) float64 {
	return x*x + 3
}

// badF devolve o mesmo número que f e, de brinde, escreve em `value`.
//
// É o tipo de função a evitar ao máximo: quem lê `func(float64) float64` não
// tem como saber que ela mexe em algo fora. A matemática não tem equivalente
// para isso — não existe "f(2) = 7 e, a propósito, o papel agora está
// diferente".
func badF(x float64) float64 {
	value = x*x + 3
	return value
}

func main() {
	fmt.Println("f(x) = x² + 3")
	fmt.Println()

	// Cada resultado é guardado ANTES de ler `value`. Não é preciosismo: a
	// spec do Go só garante a ordem entre chamadas de função, então ler o
	// global no meio da lista de argumentos de um Printf poderia acontecer
	// antes da chamada — e o slide mostraria o número errado.

	// A pura, duas vezes. Mesma entrada, mesma saída, e o mundo intacto.
	r1 := f(2)
	fmt.Printf("  f(2)    = %-4v value = %v\n", r1, value)
	r2 := f(2)
	fmt.Printf("  f(2)    = %-4v value = %v   <- nada mudou lá fora\n", r2, value)

	// A impura. Repare que o retorno é IDÊNTICO. O que muda é o que sobra.
	r3 := badF(2)
	fmt.Printf("  badF(2) = %-4v value = %v   <- o programa não é mais o mesmo\n", r3, value)
	fmt.Println()

	// O teste do papel: riscar a chamada e escrever o resultado no lugar.
	comChamada := f(2) + f(2)
	comValor := 7.0 + 7.0
	fmt.Printf("  f(2) + f(2) = %v   e   7 + 7 = %v   -> trocar uma pela outra é seguro\n", comChamada, comValor)
	fmt.Println("  com badF, essa troca apagaria a escrita em `value`: a expressão")
	fmt.Println("  e o seu resultado deixaram de ser a mesma coisa.")
}
