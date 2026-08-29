package paradigma

import (
	"errors"
	"testing"
)

func TestAddOne(t *testing.T) {
	tests := []struct {
		input    int
		expected int
		err      error
	}{
		{1, 2, nil},
		{2, 3, nil},
		{0, 0, errors.New("x cannot be 0")},
	}

	for _, tt := range tests {
		result, err := AddOne(tt.input)
		if result != tt.expected || (err != nil && err.Error() != tt.err.Error()) {
			t.Errorf("AddOne(%d) = %d, %v; want %d, %v", tt.input, result, err, tt.expected, tt.err)
		}
	}
}

func TestDouble(t *testing.T) {
	tests := []struct {
		input    int
		expected int
		err      error
	}{
		{1, 2, nil},
		{2, 4, nil},
		{0, 0, errors.New("x cannot be 0")},
	}

	for _, tt := range tests {
		result, err := Double(tt.input)
		if result != tt.expected || (err != nil && err.Error() != tt.err.Error()) {
			t.Errorf("Double(%d) = %d, %v; want %d, %v", tt.input, result, err, tt.expected, tt.err)
		}
	}
}

func TestCompose(t *testing.T) {
	tests := []struct {
		f        functionAsType
		g        functionAsType
		expected int
	}{
		{
			func() (int, error) { return AddOne(1) }, // 2
			func() (int, error) { return Double(2) }, // 4
			6,                                         // 2 + 4 = 6 (o teste de 2023 esperava 3, e por isso falhava)
		},
		{
			func() (int, error) { return AddOne(1) },
			func() (int, error) { return Double(0) },
			0, // erro na segunda função
		},
		{
			func() (int, error) { return AddOne(0) },
			func() (int, error) { return Double(2) },
			0, // erro na primeira função
		},
	}

	for _, tt := range tests {
		resultFunc := Compose(tt.f, tt.g)
		if result := resultFunc(); result != tt.expected {
			t.Errorf("Compose() = %d; want %d", result, tt.expected)
		}
	}
}

func TestComposeFn(t *testing.T) {
	// (Double ∘ AddOne)(3) = Double(AddOne(3)) = Double(4) = 8
	pipeline := ComposeFn(AddOne, Double)

	got, err := pipeline(3)
	if err != nil || got != 8 {
		t.Errorf("ComposeFn(AddOne, Double)(3) = %d, %v; want 8, nil", got, err)
	}

	// O erro do primeiro estágio corta o pipeline (curto-circuito).
	if _, err := pipeline(0); err == nil {
		t.Error("ComposeFn deveria propagar o erro de AddOne(0)")
	}
}
