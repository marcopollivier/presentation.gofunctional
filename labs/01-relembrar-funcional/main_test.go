package main

import (
	"fmt"
	"testing"
)

// Os testes são a segunda metade do argumento do slide.
//
// Repare no custo: para testar f() não é preciso preparar nada nem limpar nada
// depois. Para testar badF() é preciso lembrar de zerar `value` — senão um
// teste contamina o próximo. Essa cerimônia extra não é chatice de linter; é a
// regra 6 do CLAUDE.md aparecendo em forma de trabalho. Uma função que precisa
// de arrumação em volta para ser testada já contou onde ela impura.

// f é pura: mesma entrada, mesma saída, sempre.
func TestFEhDeterministica(t *testing.T) {
	casos := []struct {
		x    float64
		quer float64
	}{
		{x: 0, quer: 3},
		{x: 2, quer: 7},  // o número do slide
		{x: -2, quer: 7}, // x² não liga para o sinal
		{x: 3, quer: 12},
	}

	for _, c := range casos {
		if got := f(c.x); got != c.quer {
			t.Errorf("f(%v) = %v, esperado %v", c.x, got, c.quer)
		}
		// De novo, o mesmo argumento: nada de "na segunda vez é diferente".
		if got := f(c.x); got != c.quer {
			t.Errorf("f(%v) na 2ª chamada = %v, esperado %v", c.x, got, c.quer)
		}
	}
}

// Chamar f não deixa rastro: o estado global continua onde estava.
func TestFNaoTocaNoGlobal(t *testing.T) {
	defer restauraValue(t)
	value = 42

	f(2)

	if value != 42 {
		t.Errorf("value = %v depois de f(2), esperado 42 intacto", value)
	}
}

// badF devolve o MESMO número que f — e é exatamente por isso que ela engana.
// O que a denuncia é o efeito, não o retorno.
func TestBadFRetornaIgualMasSujaOGlobal(t *testing.T) {
	defer restauraValue(t)
	value = 0

	got := badF(2)

	if got != f(2) {
		t.Errorf("badF(2) = %v, esperado o mesmo que f(2) = %v", got, f(2))
	}
	if value != 7 {
		t.Errorf("value = %v depois de badF(2), esperado 7 — o efeito É o ponto", value)
	}
}

// restauraValue é a cerimônia que só a função impura obriga a escrever.
func restauraValue(t *testing.T) {
	t.Helper()
	value = 0
}

// Example é a demo de palco, com saída verificada pelo próprio teste.
func Example() {
	// Transparência referencial: a expressão e o seu valor são intercambiáveis.
	fmt.Println(f(2) + f(2))
	fmt.Println(7.0 + 7.0)

	// Output:
	// 14
	// 14
}
