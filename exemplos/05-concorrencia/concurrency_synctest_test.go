package concorrencia

import (
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

// Go 1.25 — testing/synctest. Teste DETERMINÍSTICO de código concorrente.
//
// Na palestra de 2023, o argumento "concorrência combina com o modelo
// funcional" ficava na base da percepção pessoal — difícil de provar porque
// teste de concorrência com tempo real é flaky. O synctest fecha essa lacuna:
// dentro do "bubble", o relógio é falso e só avança quando TODAS as goroutines
// estão bloqueadas. Um teste que dependeria de sleeps reais roda instantâneo e
// sem flakiness.
//
// A ponte com o resto da palestra: a goroutine abaixo não compartilha estado
// mutável — cada uma calcula e devolve seu resultado por um canal. É
// concorrência no estilo funcional (comunicar em vez de compartilhar), e é
// justamente esse estilo que fica trivial de testar de forma determinística.

// fanOutSquares dispara uma goroutine por número, cada uma "trabalhando" por
// um tempo antes de devolver o quadrado. Sem estado compartilhado: só canais.
func fanOutSquares(nums []int) []int {
	results := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Duration(n) * time.Second) // tempo "real" — falso no bubble
			results[i] = n * n
		}()
	}
	wg.Wait()
	return results
}

func TestFanOutSquares(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		got := fanOutSquares([]int{1, 2, 3, 4})
		want := []int{1, 4, 9, 16}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("fanOutSquares = %v; want %v", got, want)
			}
		}
	})
	// Apesar dos sleeps de até 4 segundos, o teste roda instantaneamente:
	// o relógio do bubble avança sozinho quando todas as goroutines bloqueiam.
}
