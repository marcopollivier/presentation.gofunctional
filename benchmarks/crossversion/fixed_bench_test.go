package crossversion

import "testing"

// A medição corrigida.
//
// A regra que o benchmark de 2023 violou: toda iteração tem que receber a
// MESMA qualidade de entrada. Aqui os dados são restaurados a cada volta, e o
// custo dessa restauração fica fora do cronômetro.
//
// As três medidas e o que se subtrai delas:
//
//	InPlace       custo puro do sort (restauração fora do cronômetro)
//	InPlaceCopy   sort + cópia, buffer reaproveitado, zero alocação
//	Immutable     sort + cópia + alocação nova
//
//	Immutable - InPlaceCopy  = custo isolado da ALOCAÇÃO
//	Immutable - InPlace      = custo total da IMUTABILIDADE
//
// Mais duas testemunhas, que não estavam no desenho original e existem para
// não deixar o n=5 sem resposta (ver BenchmarkInPlace abaixo):
//
//	CopyOnly      só a cópia, sem sort e sem mexer no cronômetro
//	TimerOverhead o par StopTimer/StartTimer vazio
//
//	InPlaceCopy - CopyOnly   = custo puro do sort, por uma via que NÃO
//	                           manipula o cronômetro — válida inclusive em n=5

// BenchmarkInPlace mede o custo puro do sort.
//
// ATENÇÃO — a ressalva que precisa estar no slide: b.StopTimer()/StartTimer()
// custam tempo, e esse custo cai em cima de CADA iteração. Para n=5, onde o
// sort inteiro é da ordem de dezenas de nanossegundos, o overhead é da mesma
// ordem de grandeza do que se quer medir, e o número resultante diz mais sobre
// o cronômetro do que sobre a ordenação. A partir de 1k é desprezível.
//
// BenchmarkTimerOverhead quantifica esse viés, e a subtração
// InPlaceCopy - CopyOnly dá o mesmo número por um caminho que não sofre dele.
// Em n >= 1k as duas vias têm que convergir; se não convergirem, alguma das
// duas está errada e isso é conteúdo.
func BenchmarkInPlace(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			original := makeData(size.n)
			data := make([]int, size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Restauração fora do cronômetro: é isto que garante que toda
				// iteração receba dados desordenados, e é exatamente o que
				// faltava em 2023.
				b.StopTimer()
				copy(data, original)
				b.StartTimer()

				SortInPlace(data)
			}
		})
	}
}

// BenchmarkInPlaceCopy mede sort + cópia sem alocar: o buffer é criado uma vez
// e reaproveitado.
//
// Não usa StopTimer — a cópia É parte do que se mede — e por isso serve de
// contraprova ao número de n=5 do BenchmarkInPlace.
func BenchmarkInPlaceCopy(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			original := makeData(size.n)
			buf := make([]int, size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SortInPlaceCopy(buf, original)
			}
		})
	}
}

// BenchmarkImmutable mede o caminho puro completo: alocar, copiar, ordenar.
// A diferença para InPlaceCopy é o preço do allocator e do GC.
func BenchmarkImmutable(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			original := makeData(size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SortImmutable(original)
			}
		})
	}
}

// BenchmarkCopyOnly é a testemunha da cópia: mesmo laço do InPlaceCopy, sem o
// sort. Serve para obter o custo puro do sort por subtração, sem cronômetro
// manipulado.
func BenchmarkCopyOnly(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			original := makeData(size.n)
			buf := make([]int, size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				CopyOnly(buf, original)
			}
		})
	}
}

// BenchmarkTimerPair mede quanto custa o par StopTimer/StartTimer — o viés
// embutido no BenchmarkInPlace.
//
// Duas armadilhas foram encontradas medindo isto, e as duas são conteúdo:
//
//  1. Um benchmark cujo corpo é só `b.StopTimer(); b.StartTimer()` NÃO TERMINA.
//     O tempo cronometrado por iteração é ~zero (o cronômetro está parado
//     justamente durante o que se quer medir), então o framework sobe b.N
//     tentando alcançar o -benchtime e o relógio de parede explode. Testado nos
//     dois formatos ingênuos; nenhum converge.
//
//  2. O motivo do custo é maior do que "chamada de função". Em
//     $GOROOT/src/testing/benchmark.go, StopTimer e StartTimer chamam
//     runtime.ReadMemStats — que é STOP-THE-WORLD. Cada par para o mundo duas
//     vezes. Esse custo é de microssegundos e, pior para quem tenta medi-lo,
//     fica FORA da conta: em StopTimer o ReadMemStats vem depois de fechar a
//     duração, e em StartTimer vem antes de reabrir.
//
// O que sobra dentro da conta é o resíduo — leitura de relógio, chamadas, e o
// rastro do STW nos caches. É pequeno em absoluto e devastador em n=5.
//
// A medição que funciona é diferencial e com carga grande o bastante para
// segurar o b.N: os dois sub-benchmarks fazem o MESMO trabalho de 1k
// elementos, e um carrega o par vazio a mais.
//
//	com - sem = resíduo do par que ENTRA na conta, por iteração
//
// Esse é o número a subtrair do BenchmarkInPlace em n=5. A via alternativa,
// InPlaceCopy - CopyOnly, chega ao custo do sort sem passar por aqui, e as
// duas têm que fechar.
func BenchmarkTimerPair(b *testing.B) {
	const n = 1000 // grande o bastante para b.N não explodir; ver comentário acima
	original := makeData(n)
	buf := make([]int, n)

	b.Run("sem", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			SortInPlaceCopy(buf, original)
		}
	})

	b.Run("com", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			b.StartTimer()
			SortInPlaceCopy(buf, original)
		}
	})
}
