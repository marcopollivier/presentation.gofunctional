package crossversion

import "testing"

// O benchmark de 2023, com o erro PRESERVADO.
//
// Não conserte nada aqui. Este arquivo é uma peça de museu: ele existe para
// medir o artefato, não o Go. Rodando-o nos sete toolchains eu isolo o quanto
// do "custo da imutabilidade" que apresentei em 2023 era viés de medição — e
// esse número não tem nada a ver com evolução da linguagem.
//
// O erro: sort.Ints ordena IN-PLACE. A partir da segunda iteração do loop, o
// slice já está ordenado, e o lado mutável passa a competir com a melhor
// entrada possível (pdqsort detecta sequência ordenada e sai quase de graça).
// O lado imutável, que nunca toca no original, continua recebendo dados
// aleatórios em toda iteração. Não é uma comparação entre dois algoritmos; é
// uma comparação entre "ordenar dados ordenados" e "ordenar dados aleatórios".
//
// O sintoma que se denuncia sozinho: ns/op praticamente igual para n=5 e para
// n=1000. Nenhum sort de verdade se comporta assim.
//
// E a causa não foi descuido: foi NÃO-LOCALIDADE. `func SortInPlace(s []int)`
// não avisa que destrói o argumento. Ver TestLegacyBenchmarkIsBiased.

func BenchmarkLegacyInPlace(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			// A entrada é criada UMA vez, fora do loop — e é por isso que ela
			// chega ordenada da segunda iteração em diante.
			data := makeData(size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SortInPlace(data)
			}
		})
	}
}

func BenchmarkLegacyImmutable(b *testing.B) {
	for _, size := range sizes {
		size := size
		b.Run(size.name, func(b *testing.B) {
			// Mesmo código do lado mutável — mas como SortImmutable não escreve
			// em `data`, aqui a entrada continua aleatória em toda iteração.
			// A assimetria está inteira nesta diferença invisível.
			data := makeData(size.n)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SortImmutable(data)
			}
		})
	}
}
