package gofunctional

import (
	"iter"
	"slices"
)

// Go 1.23 — iterators e range-over-func. A mudança mais importante para esta
// palestra. `iter.Seq[V]` é literalmente `func(yield func(V) bool)`: uma
// sequência que produz valores sob demanda. Com isso dá para escrever
// map/filter/reduce LAZY e COMPONÍVEL, sem materializar slices intermediários.
//
// Isso ataca de frente a conclusão do benchmark de coleções de 2023: o custo
// que eu media era o de COPIAR a coleção a cada transformação. Com iterators,
// os valores fluem um a um e não há cópia intermediária nenhuma.

// MapSeq aplica f a cada elemento, preguiçosamente. Nada roda até alguém
// consumir a sequência de saída.
func MapSeq[A, B any](s iter.Seq[A], f func(A) B) iter.Seq[B] {
	return func(yield func(B) bool) {
		for a := range s {
			if !yield(f(a)) {
				return
			}
		}
	}
}

// FilterSeq deixa passar só os elementos que satisfazem keep.
func FilterSeq[A any](s iter.Seq[A], keep func(A) bool) iter.Seq[A] {
	return func(yield func(A) bool) {
		for a := range s {
			if keep(a) && !yield(a) {
				return
			}
		}
	}
}

// Reduce colapsa a sequência num único valor. É o consumidor: é aqui que o
// pipeline lazy finalmente roda.
func Reduce[A, B any](s iter.Seq[A], init B, f func(B, A) B) B {
	acc := init
	for a := range s {
		acc = f(acc, a)
	}
	return acc
}

// SumSquaresOfEvens compõe Filter -> Map -> Reduce sobre um pipeline LAZY.
// slices.Values transforma o slice em iter.Seq; daí em diante nenhum slice
// intermediário é alocado. Os valores atravessam os três estágios um a um.
func SumSquaresOfEvens(nums []int) int {
	seq := slices.Values(nums)
	evens := FilterSeq(seq, func(n int) bool { return n%2 == 0 })
	squares := MapSeq(evens, func(n int) int { return n * n })
	return Reduce(squares, 0, func(acc, n int) int { return acc + n })
}

// SumSquaresOfEvensSlices faz o MESMO cálculo materializando cada etapa num
// slice novo — o jeito pré-1.23. Cada estágio ALOCA um backing array. É
// exatamente o custo que o benchmark de 2023 media e chamava de "o preço da
// abordagem funcional". O benchmark comparando as duas (ver o teste) mostra
// que o custo era da CÓPIA, não do paradigma.
func SumSquaresOfEvensSlices(nums []int) int {
	var evens []int
	for _, n := range nums {
		if n%2 == 0 {
			evens = append(evens, n)
		}
	}
	squares := make([]int, len(evens))
	for i, n := range evens {
		squares[i] = n * n
	}
	sum := 0
	for _, n := range squares {
		sum += n
	}
	return sum
}
