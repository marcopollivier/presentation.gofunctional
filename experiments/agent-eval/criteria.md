# Critérios objetivos

Aplicar a cada saída (A e B) e preencher. Números e sim/não — nada subjetivo na
tabela; o subjetivo vai em `notes.md`.

| Critério | Como medir | Condição A | Condição B |
|--|--|--|--|
| Efeitos colaterais | Contar: estado global mutável + I/O escondido + mutação de argumento | | |
| Muta a entrada? | `items` do chamador é alterado após a chamada? (teste que clona e compara) | | |
| Testável sem mock? | O teste gerado usa mock/fixture, ou é table-driven direto? | | |
| Alocações (`TopKByFrequency`) | `go test -bench -benchmem` numa entrada fixa; anotar allocs/op | | |
| Erro/边界 no tipo? | `k` inválido, entrada vazia: tratado no retorno ou com panic/estado? | | |
| Bug de mutação em `NormalizeAndRank` | A normalização muta o slice recebido? | | |
| Linhas na maior função | `gocyclo`/contagem manual da maior função | | |

## Como medir "muta a entrada"

Teste padrão a rodar sobre as duas saídas:

```go
func TestNoMutation(t *testing.T) {
    in := []string{"B", "a", "A", " b "}
    snapshot := slices.Clone(in)
    _ = NormalizeAndRank(in, 2)
    if !slices.Equal(in, snapshot) {
        t.Errorf("mutou a entrada: %v -> %v", snapshot, in)
    }
}
```

Se este teste **falha** na saída, o bug de mutação in-place apareceu — o mesmo
que corrompeu o benchmark de 2023, agora reproduzido por um agente.

## Registro

Copie esta tabela para `results/<data>-A/criteria.md` e
`results/<data>-B/criteria.md`, uma coluna preenchida em cada. A comparação
lado a lado é o slide.
