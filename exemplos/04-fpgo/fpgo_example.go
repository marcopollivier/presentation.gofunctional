package fpgo

import (
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	R "github.com/IBM/fp-go/v2/result"
)

// Bloco "exemplos e opções de implementação": a fp-go como solução para o
// Either. Em 2023 a lib era work-in-progress; hoje está na v2, exige Go 1.24+
// e usa GENERIC TYPE ALIASES — o próprio Result é `type Result[T] = Either[error, T]`,
// que só passou a ser legal com o 1.24. A evolução da linguagem é o que
// destravou a biblioteca.
//
// O ponto didático (e a ponte para o submódulo evolution-127): repare que a
// composição aqui é `Pipe3(x, Chain(...), Map(...), GetOrElse(...))` — funções
// de pacote costuradas por Pipe, porque ANTES do Go 1.27 um método não podia
// ter seus próprios parâmetros de tipo. Com métodos genéricos isso vira
// `x.Chain(...).Map(...)`. Mesmo Either, ergonomia diferente.

// parsePositive: string -> Result[int]. Falha (Left) se não for número ou se
// não for positivo. O erro vive no tipo de retorno — não há panic, não há
// segundo canal escondido.
func parsePositive(s string) R.Result[int] {
	n, err := strconv.Atoi(s)
	if err != nil {
		return R.Left[int](err)
	}
	if n <= 0 {
		return R.Left[int](fmt.Errorf("%q não é positivo", s))
	}
	return R.Of(n)
}

// atLeastTen: int -> Result[int]. Uma função Kleisli (A -> Result[B]),
// candidata natural a Chain.
func atLeastTen(n int) R.Result[int] {
	if n < 10 {
		return R.Left[int](fmt.Errorf("%d é menor que 10", n))
	}
	return R.Of(n)
}

// ProcessFpGo costura o pipeline: parse (pode falhar) -> Chain valida o
// mínimo -> Map dobra -> GetOrElse extrai com um default. Qualquer Left no
// caminho curto-circuita e cai no default (-1), sem if de erro espalhado.
func ProcessFpGo(s string) int {
	return F.Pipe3(
		parsePositive(s),
		R.Chain(atLeastTen),
		R.Map(func(n int) int { return n * 2 }),
		R.GetOrElse(func(error) int { return -1 }),
	)
}
