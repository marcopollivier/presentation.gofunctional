//go:build ignore

// ceiling_hkt.go — O TETO. Este arquivo NÃO COMPILA de propósito, e é a nova
// conclusão da palestra.
//
// Go 1.27 deu métodos genéricos: Option.Map existe (ver option.go). O passo
// seguinte — higher-kinded types — seria abstrair sobre "qualquer coisa que é
// um Functor", isto é, qualquer container que TENHA um Map. Em linguagens com
// HKT você escreve código genérico sobre `F[_]`. Em Go, não.
//
// A barreira está numa linha só: uma interface não consegue sequer EXIGIR um
// Map genérico, porque método de interface não pode declarar parâmetro de tipo.
// (E o espelho disso, dito nas release notes: um método genérico também não
// pode SATISFAZER um método de interface. As duas metades se fecham.)
//
// Está sob `//go:build ignore` para não quebrar `go build ./...`. Para ver o
// erro no palco:
//     go build ceiling_hkt.go
//
// Erro do compilador (go1.27rc1), registrado após rodar de verdade:
//     ./ceiling_hkt.go:34:5: interface method must have no type parameters
//     ./ceiling_hkt.go:34:25: undefined: B   (cascata do erro acima)
//
// Consequência: dá para ter Option.Map. NÃO dá para escrever código que
// abstrai sobre "todo Functor". A forçação de barra não sumiu — ela saiu da
// SINTAXE (resolvida pelos métodos genéricos) e foi para a CAPACIDADE DE
// ABSTRAÇÃO. É onde Go traça o limite: monads concretos sim, HKT não.

package evolution

// Functor tentaria dizer "tem um Map que leva A para qualquer B". A intenção
// é HKT: abstrair sobre o container. Mas a linha abaixo é rejeitada pelo
// compilador — método de interface não pode ter parâmetro de tipo próprio.
type Functor interface {
	Map[B any](f func(any) B) any
}
