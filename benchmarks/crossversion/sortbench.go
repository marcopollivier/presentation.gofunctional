package crossversion

import "sort"

// Núcleo comparável do benchmark: o código que precisa compilar e rodar
// IDÊNTICO nos toolchains 1.21 a 1.27.
//
// Isso proíbe qualquer API posterior ao 1.21 daqui para baixo:
//
//	slices.Sorted / slices.Values  -> 1.23
//	iter.Seq                       -> 1.23
//	testing.B.Loop                 -> 1.24
//	math/rand/v2                   -> 1.22
//
// Sobram `sort.Ints`, `make` e `copy`, que é exatamente o vocabulário do
// benchmark de 2023 — e é isso que torna a comparação legítima.
//
// Sobre sort.Ints: ela NÃO está deprecada (não há marcador `Deprecated:` no
// doc, nem no 1.27). Desde o 1.22 ela apenas chama slices.Sort internamente.
// Essa mudança de implementação é parte do que a matriz mede, não um problema:
// o fonte é o mesmo, o que muda embaixo é o Go.

// sink existe para impedir dead code elimination.
//
// Sem `testing.B.Loop` (1.24+), que segura os valores por construção, o
// compilador pode provar que o resultado de um sort nunca é observado e apagar
// a chamada inteira. Escrever num global de pacote torna o resultado observável
// de fora e o corpo do loop indelével.
var sink []int

// SortInPlace é a versão mutável: ordena o slice do chamador e destrói a
// entrada dele.
//
// A assinatura não dá nenhum sinal disso — sem `*`, sem anotação, sem nada. É
// a tese da palestra em uma linha, e foi exatamente esse silêncio que corrompeu
// o benchmark de 2023 (ver legacy_bench_test.go).
func SortInPlace(s []int) {
	sort.Ints(s)
	sink = s
}

// SortInPlaceCopy copia para um buffer JÁ EXISTENTE e ordena a cópia.
//
// Preserva a entrada do chamador sem alocar nada: é a parcela "cópia + sort" do
// custo da imutabilidade, isolada da parcela "alocação". A diferença para
// SortImmutable é, portanto, o preço do allocator e do GC — que é onde se
// espera ver o efeito do Green Tea GC (1.26) e da alocação especializada por
// tamanho (1.27).
func SortInPlaceCopy(dst, src []int) {
	copy(dst, src)
	sort.Ints(dst)
	sink = dst
}

// SortImmutable é a versão pura: aloca, copia e ordena a cópia.
//
// O chamador não tem como perceber que algo foi ordenado — a mutação existe
// (sort.Ints reescreve `c` inteiro), mas morre aqui dentro. Pureza não é
// abstinência de mutação, é a mutação não escapar.
func SortImmutable(s []int) []int {
	c := make([]int, len(s))
	copy(c, s)
	sort.Ints(c)
	sink = c
	return c
}

// CopyOnly é a testemunha: só a cópia, sem ordenar.
//
// Existe para dar o custo puro do sort por subtração — InPlaceCopy menos
// CopyOnly — SEM manipular o cronômetro. É a via que continua válida em n=5,
// onde o overhead de StopTimer/StartTimer distorce o BenchmarkInPlace.
func CopyOnly(dst, src []int) {
	copy(dst, src)
	sink = dst
}

// xorshift64 é um PRNG de cinco linhas, com seed fixa.
//
// Não é math/rand de propósito: a stdlib trocou de implementação e de
// comportamento de seeding entre 1.20 e 1.22 (auto-seed, depois math/rand/v2).
// Usar a stdlib faria cada toolchain ordenar uma sequência DIFERENTE de
// números, e a comparação entre colunas perderia o sentido. Com xorshift64 o
// 1.21 e o 1.27 ordenam exatamente os mesmos dados.
type xorshift64 struct{ state uint64 }

func newRand(seed uint64) *xorshift64 { return &xorshift64{state: seed} }

func (r *xorshift64) next() uint64 {
	r.state ^= r.state << 13
	r.state ^= r.state >> 7
	r.state ^= r.state << 17
	return r.state
}

// seed fixa: a mesma sequência em toda versão, toda rodada, toda máquina.
const benchSeed = 0x9E3779B97F4A7C15

// makeData devolve n inteiros pseudoaleatórios determinísticos.
func makeData(n int) []int {
	r := newRand(benchSeed)
	s := make([]int, n)
	for i := range s {
		// >>32 deixa os valores em 0..4.29e9: dez dígitos, legíveis quando o
		// teste é projetado no telão, e faixa larga o bastante para as
		// duplicatas serem desprezíveis mesmo em n=1M (~116 colisões
		// esperadas). Faixa estreita mudaria o comportamento do pdqsort, que
		// tem caminho próprio para muitos elementos repetidos.
		s[i] = int(r.next() >> 32)
	}
	return s
}

// sizes são os quatro casos da palestra. Cada um roda como sub-benchmark
// próprio (b.Run), e não como quatro invocações do mesmo caso com rótulo
// trocado — que foi o que os slides de 2023 sugerem ter acontecido.
var sizes = []struct {
	name string
	n    int
}{
	{"5", 5},
	{"1k", 1000},
	{"10k", 10000},
	{"1M", 1000000},
}
