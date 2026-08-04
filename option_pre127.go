package gofunctional

// Option[T] à mão — a versão PRÉ-1.27, o "antes" do divisor de águas.
//
// O tipo é trivial: um valor e um bit de presença. O problema está nas
// operações. Map leva Option[A] para Option[B] — precisa introduzir um novo
// parâmetro de tipo B. Até o Go 1.27, um MÉTODO não podia declarar seus
// próprios parâmetros de tipo: `func (o Option[A]) Map[B](f func(A) B) Option[B]`
// era ilegal. Por isso Map/Chain têm que ser FUNÇÕES DE PACOTE.
//
// A consequência é ergonômica: sem métodos, não há encadeamento. Você aninha
// chamadas ou cria variáveis intermediárias. É por isso que as libs de FP em
// Go pré-1.27 obrigam `Pipe3(x, Map(f), Chain(g))` em vez de `x.Map(f).Chain(g)`.
// O submódulo evolution-127/ mostra o mesmo tipo com métodos genéricos.

type Option[T any] struct {
	value   T
	present bool
}

func Some[T any](v T) Option[T] { return Option[T]{value: v, present: true} }
func None[T any]() Option[T]     { return Option[T]{} }

// MapOption: função livre porque a saída Option[B] introduz o tipo B.
func MapOption[A, B any](o Option[A], f func(A) B) Option[B] {
	if !o.present {
		return None[B]()
	}
	return Some(f(o.value))
}

// ChainOption (flatMap): idem — f devolve Option[B].
func ChainOption[A, B any](o Option[A], f func(A) Option[B]) Option[B] {
	if !o.present {
		return None[B]()
	}
	return f(o.value)
}

func GetOrElse[T any](o Option[T], def T) T {
	if !o.present {
		return def
	}
	return o.value
}

// HalfIfEven demonstra a dor ergonômica: valida (Chain) que n é par, divide
// por 2 (Map), extrai com default (GetOrElse). Sem encadeamento, o fluxo vira
// uma escada de variáveis intermediárias. Compare com evolution-127/option.go.
func HalfIfEven(n int) int {
	step1 := Some(n)
	step2 := ChainOption(step1, func(n int) Option[int] {
		if n%2 != 0 {
			return None[int]()
		}
		return Some(n)
	})
	step3 := MapOption(step2, func(n int) int { return n / 2 })
	return GetOrElse(step3, -1)
}
