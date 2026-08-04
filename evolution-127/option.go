package evolution

// Option[T] — a versão Go 1.27, o "depois" do divisor de águas.
//
// Compare com option_pre127.go no módulo raiz. O tipo é idêntico; a diferença
// é uma só, e é toda a conclusão da palestra: aqui Map e Chain declaram seus
// PRÓPRIOS parâmetros de tipo NO MÉTODO (`Map[B any]`). Até o Go 1.27 isso era
// ilegal — um método não podia introduzir novos type params —, e por isso as
// libs de FP obrigavam `Pipe3(x, Map(f), Chain(g))`.
//
// Com métodos genéricos, o mesmo Option vira ENCADEÁVEL:
//     Some(n).Chain(f).Map(g).GetOrElse(d)
// Go ganhou monads concretos e encadeáveis. O que ele ainda NÃO tem está em
// ceiling_hkt.go.

type Option[T any] struct {
	value   T
	present bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, present: true} }
func None[T any]() Option[T]    { return Option[T]{} }

// Map é um MÉTODO GENÉRICO: declara seu próprio parâmetro de tipo B. Esta
// linha é a fronteira entre 2023 e 2026 — não compilava antes do 1.27.
func (o Option[A]) Map[B any](f func(A) B) Option[B] {
	if !o.present {
		return None[B]()
	}
	return Some(f(o.value))
}

// Chain (flatMap) — idem, com seu próprio B.
func (o Option[A]) Chain[B any](f func(A) Option[B]) Option[B] {
	if !o.present {
		return None[B]()
	}
	return f(o.value)
}

// GetOrElse não precisa de novo type param — usa o T do próprio Option.
func (o Option[T]) GetOrElse(def T) T {
	if !o.present {
		return def
	}
	return o.value
}

// HalfIfEven: o MESMO fluxo de option_pre127.go, agora como uma cadeia fluente
// em vez de uma escada de variáveis intermediárias. É a ergonomia que o 1.27
// destravou.
func HalfIfEven(n int) int {
	return Some(n).
		Chain(func(n int) Option[int] {
			if n%2 != 0 {
				return None[int]()
			}
			return Some(n)
		}).
		Map(func(n int) int { return n / 2 }).
		GetOrElse(-1)
}
