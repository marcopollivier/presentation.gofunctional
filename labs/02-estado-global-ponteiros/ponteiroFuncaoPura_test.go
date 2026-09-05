package main

import (
	"fmt"
	"testing"
)

// Os testes fixam por escrito o que os três momentos mostram na tela — e um
// deles fixa justamente o desconforto do lab: a mesma chamada, com o mesmo
// argumento, respondendo duas coisas diferentes.

// Len é leitura pura: chamada duas vezes seguidas, mesma resposta, e o valor
// apontado continua onde estava.
func TestLenNaoAlteraNada(t *testing.T) {
	nome := "Marco Ollivier"

	primeira := Len(&nome)
	segunda := Len(&nome)

	if primeira != segunda {
		t.Errorf("Len(&nome) devolveu %d e depois %d — deveria ser leitura pura", primeira, segunda)
	}
	if nome != "Marco Ollivier" {
		t.Errorf("nome = %q depois de duas leituras, esperado intacto", nome)
	}
	if primeira != 14 {
		t.Errorf("Len(&nome) = %d, esperado 14", primeira)
	}
}

// O estado a mais que só a versão com ponteiro tem. Não existe string nil:
// `Len(nome string)` não conseguiria nem chegar neste caso.
func TestLenComNilEntraEmPanico(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Len(nil) deveria entrar em pânico — é o estado extra que o ponteiro traz")
		}
	}()

	Len(nil)
}

// Momento 3: Changer escreve na variável do chamador. Mesma assinatura de Len,
// verbo oposto.
func TestChangerReescreveOChamador(t *testing.T) {
	nome := "Marco Ollivier"

	Changer(&nome)

	if nome != "Marco Paulo Ollivier" {
		t.Errorf("nome = %q depois de Changer, esperado %q — a escrita É o ponto", nome, "Marco Paulo Ollivier")
	}
}

// A caixa é a mesma; o conteúdo é que mudou. É o endereço repetido nas três
// linhas da demo.
func TestOEnderecoNaoMudaDepoisDaEscrita(t *testing.T) {
	nome := "Marco Ollivier"
	antes := fmt.Sprintf("%p", &nome)

	Changer(&nome)
	depois := fmt.Sprintf("%p", &nome)

	if antes != depois {
		t.Errorf("endereço de nome mudou de %s para %s — passar ponteiro não move a variável", antes, depois)
	}
}

// O preço do ponteiro, em forma de teste: MESMO argumento, resposta diferente.
//
// Len continua pura — não escreveu nada, não guardou nada. O que se perdeu foi
// a transparência referencial: `Len(&nome)` não pode mais ser substituída pelo
// seu resultado, porque o resultado depende de QUANDO ela é chamada. Com
// `Len(nome)` por valor, o argumento seria o texto, e 14 seria 14 para sempre.
func TestMesmoArgumentoRespostaDiferente(t *testing.T) {
	nome := "Marco Ollivier"

	antes := Len(&nome)
	Changer(&nome)
	depois := Len(&nome)

	if antes != 14 || depois != 20 {
		t.Errorf("Len(&nome) = %d antes e %d depois, esperado 14 e 20", antes, depois)
	}
	if antes == depois {
		t.Error("as duas chamadas coincidiram — o exemplo perdeu o ponto")
	}
}
