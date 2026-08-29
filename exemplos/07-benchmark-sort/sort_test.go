package ordenacao

import (
	"slices"
	"testing"
)

func TestSortNative(t *testing.T) {
	data := []int{3, 1, 4, 1, 5, 9}
	SortNative(data)
	expected := []int{1, 1, 3, 4, 5, 9}
	if !slices.Equal(data, expected) {
		t.Errorf("SortNative() = %v; want %v", data, expected)
	}
}

func TestSortImmutable(t *testing.T) {
	original := []int{3, 1, 4, 1, 5, 9}
	data := slices.Clone(original)

	sorted := SortImmutable(data)

	if want := []int{1, 1, 3, 4, 5, 9}; !slices.Equal(sorted, want) {
		t.Errorf("SortImmutable() = %v; want %v", sorted, want)
	}
	// A coleção original não pode ter sido alterada — é o contrato funcional.
	if !slices.Equal(data, original) {
		t.Errorf("SortImmutable() mutou o slice original; got %v, want %v", data, original)
	}
}

func TestSortLazy(t *testing.T) {
	original := []int{3, 1, 4, 1, 5, 9}
	data := slices.Clone(original)

	sorted := SortLazy(data)

	if want := []int{1, 1, 3, 4, 5, 9}; !slices.Equal(sorted, want) {
		t.Errorf("SortLazy() = %v; want %v", sorted, want)
	}
	if !slices.Equal(data, original) {
		t.Errorf("SortLazy() mutou o slice original; got %v, want %v", data, original)
	}
}
