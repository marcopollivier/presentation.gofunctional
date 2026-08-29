package paradigma

import "errors"

// Bloco "aplicando conceitos": funções puras e function composition.
//
// AddOne e Double são puras: mesma entrada, mesma saída, sem efeito colateral.
// O erro no tipo de retorno é o contrato explícito — você lê a assinatura e
// já sabe o universo do que pode acontecer. É esse barateamento da revisão
// que a palestra nova liga ao trabalho com agentes.

func AddOne(x int) (int, error) {
	if x == 0 {
		return 0, errors.New("x cannot be 0")
	}
	return x + 1, nil
}

func Double(x int) (int, error) {
	if x == 0 {
		return 0, errors.New("x cannot be 0")
	}
	return x * 2, nil
}

type functionAsType func() (int, error)

// Compose (versão 2023, preservada como didática): recebe dois thunks já
// "aplicados", avalia os dois de forma EAGER e soma. Em qualquer erro devolve
// um closure que produz 0. É a resposta à pergunta "dá pra fazer composição
// em Go?" — dá, mas de um jeito que não é composição de verdade: os valores
// já foram calculados antes de compor.
func Compose(f functionAsType, g functionAsType) func() int {
	first, err1 := f()
	if err1 != nil {
		return func() int { return 0 }
	}

	second, err2 := g()
	if err2 != nil {
		return func() int { return 0 }
	}

	return func() int { return first + second }
}

// ComposeFn é a composição matemática de verdade: (g∘f)(x) = g(f(x)).
// Aqui compomos FUNÇÕES, não resultados já calculados — nada é avaliado até
// a função devolvida ser chamada com um x. O erro é propagado pelo caminho,
// exatamente o que os monads (Either/Result) fazem de forma encadeável.
func ComposeFn(f, g func(int) (int, error)) func(int) (int, error) {
	return func(x int) (int, error) {
		y, err := f(x)
		if err != nil {
			return 0, err
		}
		return g(y)
	}
}
