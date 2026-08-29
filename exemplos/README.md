# Exemplos — roteiro de estudo

Um pacote por ponto da palestra, **na ordem em que aparecem no palco**. Cada
pasta é autocontida: dá para abrir uma, ler o comentário do topo do arquivo e
rodar os testes sem conhecer as outras.

```bash
go test ./...                      # tudo
go test ./exemplos/03-iterators/   # só um exemplo
```

## O caminho

| # | Pasta | O ponto | Release | Bloco |
|--|--|--|--|--|
| 01 | [`01-paradigma`](01-paradigma/) | Funções puras e composição. `ComposeFn` é `g∘f` de verdade — lazy; `Compose` é a resposta de 2023, eager, preservada para contraste. | — | 2 |
| 02 | [`02-loopvar`](02-loopvar/) | O footgun de 12 anos que sumiu. `[0 1 2]` hoje, `[3 3 3]` antes. | **1.22** | 3 |
| 03 | [`03-iterators`](03-iterators/) | `iter.Seq` e pipeline lazy map/filter/reduce, sem materializar slices. Ataca a conclusão do benchmark de coleções de 2023. | **1.23** | 3 |
| 04 | [`04-fpgo`](04-fpgo/) | `Result`/`Either` com `IBM/fp-go` v2, estilo `Pipe`. Generic type aliases destravaram a lib. | **1.24** | 3 |
| 05 | [`05-concorrencia`](05-concorrencia/) | `testing/synctest`: concorrência testada de forma determinística, sem flakiness. | **1.25** | 3 |
| 06 | [`06-option`](06-option/) | `Option[T]` à mão com operações como **funções de pacote** — o "antes" do 1.27, sem encadeamento. | — | 3 |
| 07 | [`07-benchmark-sort`](07-benchmark-sort/) | Mutável vs imutável vs lazy, e o benchmark refeito com metodologia correta. **O post-mortem.** | 1.24 (`B.Loop`) | 4 |
| 08 | [`08-agentes`](08-agentes/) | Puro vs estado compartilhado: o custo de revisão na era dos agentes. | — | 6 |

E fora daqui, fechando o arco:

| Pasta | O ponto | Release | Bloco |
|--|--|--|--|
| [`../evolution-127`](../evolution-127/) | O **"depois"**: o mesmo `Option`, agora com métodos genéricos e encadeável. E `ceiling_hkt.go`, o teto — por que Go tem monads concretos mas não tem higher-kinded types. | **1.27** | 3, 5 |

## Como ler

Os arquivos `.go` **são** o material. O comentário no topo de cada um explica o
papel dele na palestra — leia antes do código. Os testes são a segunda camada de
documentação: mostram o contrato exercitado, sem mock nem fixture.

Boa parte do que está aqui é **deliberadamente "errada"** — é o lado esquerdo do
slide. `SortNative` muta de propósito, `LoopVarShared` reproduz o footgun antigo,
`reviewability.go` tem um global que não deveria existir. Não são descuidos; são
o contraste. O `CLAUDE.md` da raiz lista todos e explica cada um.

## Por que `06-option` e `evolution-127` estão separados

Os dois declaram os **mesmos nomes** — `Option[T]`, `Some`, `None`,
`HalfIfEven`. É intencional: é o mesmo tipo antes e depois dos métodos genéricos
do Go 1.27. Ficam em pacotes distintos porque no mesmo pacote não compilariam.
Abra os dois lado a lado — a diferença cabe numa linha:

```go
// 06-option (o antes) — Map precisa ser função de pacote
step := MapOption(ChainOption(Some(n), f), g)

// evolution-127 (o depois) — Map é método genérico, então encadeia
Some(n).Chain(f).Map(g).GetOrElse(-1)
```
