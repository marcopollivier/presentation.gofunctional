package ordenacao

import (
	"math/rand/v2"
	"testing"
)

// POR QUE O BENCHMARK DE 2023 ESTAVA ERRADO
// =========================================
// Os slides de 2023 mostravam 5 elementos e 1.000 elementos com o MESMO
// tempo e EXATAMENTE a mesma memória. Isso não fecha. Três causas, todas
// confirmadas no código original (recuperável com `git checkout v1-devpr24`):
//
//  1. Rótulo trocado: BenchmarkSort5* e BenchmarkSort1000* chamavam os dois
//     generateSlice(1000). O caso "5 elementos" ordenava 1.000 no fonte —
//     por isso os números batiam: era o mesmo benchmark rodando duas vezes.
//
//  2. (No código commitado o setup estava fora do loop, com ResetTimer, então
//     essa causa clássica não se aplicava — mas vale registrar o cuidado.)
//
//  3. A CAUSA DE FUNDO — e o melhor material da palestra nova:
//     sort.Ints ordena IN-PLACE. Da 2ª iteração em diante, o benchmark
//     "nativo" reordenava um slice JÁ ORDENADO (quase de graça), enquanto o
//     "imutável" clonava e ordenava dados aleatórios frescos toda vez. O viés
//     era sistemático e sempre a favor do mutável.
//
//     O benchmark que "provava" que mutabilidade é mais rápida foi corrompido
//     por mutabilidade. Essa classe de erro é INEXPRESSÁVEL numa versão pura:
//     se a ordenação não mutasse a entrada, não haveria como medir "ordenar o
//     que já está ordenado" por acidente.
//
// A CORREÇÃO
// ==========
//   - Cada tamanho roda o SEU tamanho (5, 1k, 10k, 1M de verdade).
//   - Todo mundo recebe input DESORDENADO E FRESCO a cada iteração, e essa
//     restauração fica DENTRO da medição para os três — assim ninguém ordena
//     "o que já está ordenado" (a causa 3) e o custo de obter dados frescos é
//     pago por igual. O que sobra da diferença é o que a palestra quer isolar:
//     o nativo reusa um buffer (0 alloc), o imutável ALOCA um slice novo.
//   - testing.B.Loop() (Go 1.24) roda o setup uma vez e mantém entradas/saídas
//     vivas contra dead-code elimination. (Evitamos StopTimer/StartTimer por
//     iteração de propósito: para size 5 são dezenas de milhões de iterações,
//     e alternar o timer em cada uma é patologicamente lento — um ótimo
//     lembrete de que benchmark também é código que precisa de revisão.)
//
// Rode com:  make bench        (ou: go test -bench=. -benchmem ./...)

// benchRNG usa um gerador determinístico (semente fixa) para que a massa de
// dados seja reprodutível entre execuções — o oposto do rand.Seed(time.Now())
// de 2023, que tornava cada run incomparável com o anterior.
func makeUnsorted(size int) []int {
	r := rand.New(rand.NewPCG(1, 2))
	s := make([]int, size)
	for i := range s {
		s[i] = r.IntN(size)
	}
	return s
}

var sizes = []struct {
	name string
	n    int
}{
	{"5", 5},
	{"1k", 1_000},
	{"10k", 10_000},
	{"1M", 1_000_000},
}

// sink evita que o compilador elimine resultados não usados nas versões que
// retornam um slice novo.
var sink []int

// Nativo (mutável, in-place). Reusa um buffer pré-alocado: copia o input
// desordenado fresco para dentro dele e ordena in-place. A cópia é contada
// (é o preço de dar dados frescos a cada iteração, igual para todos), mas o
// buffer é reaproveitado — por isso 0 alloc. É o oposto do bug de 2023, onde
// o nativo reordenava um slice já ordenado e parecia mágico.
func BenchmarkSortNative(b *testing.B) {
	for _, sz := range sizes {
		base := makeUnsorted(sz.n)
		b.Run(sz.name, func(b *testing.B) {
			buf := make([]int, sz.n)
			for b.Loop() {
				copy(buf, base)
				SortNative(buf)
			}
		})
	}
}

// Imutável via clone. Aqui o clone É a abordagem — entra na medição de
// propósito. É o custo real de "não mutar a entrada do chamador".
func BenchmarkSortImmutable(b *testing.B) {
	for _, sz := range sizes {
		base := makeUnsorted(sz.n)
		b.Run(sz.name, func(b *testing.B) {
			for b.Loop() {
				sink = SortImmutable(base)
			}
		})
	}
}

// Imutável via iterators (slices.Sorted(slices.Values(...))). Mesma garantia
// que Immutable, escrita como pipeline lazy. Mede o custo do caminho da
// stdlib de 1.23 contra o clone manual.
func BenchmarkSortLazy(b *testing.B) {
	for _, sz := range sizes {
		base := makeUnsorted(sz.n)
		b.Run(sz.name, func(b *testing.B) {
			for b.Loop() {
				sink = SortLazy(base)
			}
		})
	}
}
