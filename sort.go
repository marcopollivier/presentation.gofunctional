package gofunctional

import "slices"

// Este arquivo é o "primeiro problema" da palestra: coleções em Go.
// Ele mostra três formas de ordenar — o mesmo resultado, três contratos
// diferentes — e é a base do benchmark que corrige o erro de 2023.

// SortNative ordena in-place: MUTA o slice do chamador e não devolve nada.
// É a versão imperativa idiomática. Em 2023 usávamos sort.Ints aqui; desde
// o Go 1.22 sort.Ints apenas delega para slices.Sort — então já não há
// motivo para o boxing de interface do pacote sort.
func SortNative(slice []int) {
	slices.Sort(slice)
}

// SortImmutable é a versão funcional escrita à mão: copia primeiro
// (slices.Clone substitui o make+copy manual de 2023), ordena a cópia e
// devolve um slice novo. O input do chamador continua intacto.
func SortImmutable(slice []int) []int {
	sorted := slices.Clone(slice)
	slices.Sort(sorted)
	return sorted
}

// SortLazy é a revelação da palestra nova: o "sort imutável" que eu
// implementei na mão em 2023 hoje é uma linha da stdlib. slices.Values
// transforma o slice num iterator (iter.Seq[int], Go 1.23) e slices.Sorted
// coleta os valores já ordenados num slice novo — sem tocar no original.
// Mesma garantia de imutabilidade de SortImmutable, expressa como pipeline.
func SortLazy(slice []int) []int {
	return slices.Sorted(slices.Values(slice))
}
