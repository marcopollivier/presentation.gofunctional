# Post-mortem: o benchmark de 2023

O benchmark que, nos slides de 2023, "provava" que mutabilidade era mais rápida
estava metodologicamente furado. Este documento registra os números originais,
as três causas do erro, a metodologia nova e os resultados de hoje. É o eixo do
arco antes/depois da palestra — e o gancho do bloco "paradigma na era dos
agentes".

## Os números de 2023

Comando registrado nos slides:
`go test -bench=. -bench=Sort5 -cpu=1 -benchmem -benchtime=5s`

| Benchmark | ops em 5s | ns/op | B/op | allocs/op |
|--|--|--|--|--|
| `BenchmarkSort5Native` | 2.097.556 | 2847 | 24 | 1 |
| `BenchmarkSort5Immutable` | 166.477 | 35100 | 8216 | 2 |
| `BenchmarkSort1000Native` | 2.123.732 | 2829 | 24 | 1 |
| `BenchmarkSort1000Immutable` | 168.615 | 34581 | 8216 | 2 |

**O que não fecha:** 5 elementos e 1000 elementos com praticamente o mesmo
tempo e **exatamente** a mesma memória. Ordenar 5 e ordenar 1000 não pode custar
o mesmo.

## As três causas (confirmadas no código de `v1-devpr24`)

1. **Rótulo trocado — confirmado no fonte, não só na CLI.**
   `BenchmarkSort5Native` e `BenchmarkSort5Immutable` chamavam, os dois,
   `generateSlice(1000)`. O caso "5 elementos" ordenava 1000 no próprio código.
   Por isso os números de 5 e 1000 batiam: **era o mesmo benchmark**. (A flag
   `-bench=Sort5` duplicada na linha de comando dos slides ainda reforçava o
   engano, mas a causa raiz estava no código.)

2. **Setup no loop — não se aplicava.** No código commitado o `generateSlice`
   estava fora do loop, com `b.ResetTimer()`. A causa clássica não era essa.

3. **A causa de fundo — mutação in-place.** `sort.Ints` ordena in-place. Da 2ª
   iteração em diante, o benchmark "nativo" reordenava um slice **já ordenado**
   (quase de graça, pdqsort é O(n) no melhor caso), enquanto o "imutável"
   clonava e ordenava dados aleatórios frescos toda vez. O viés era sistemático
   e sempre a favor do mutável.

**O achado que abre o bloco dos agentes:** o benchmark que "provava" que
mutabilidade é mais rápida foi corrompido pela própria mutabilidade que media.
E essa classe de erro é **inexpressável** numa versão pura: se a ordenação não
mutasse a entrada, não haveria como, sem querer, medir "ordenar o que já está
ordenado". Eu escrevi esse bug sozinho em 2023. Um agente escreveria o mesmo bug
em 2026, só que mais rápido e em mais lugares. Não é o agente que resolve isso;
é o paradigma que torna a classe de erro impossível de escrever.

## A metodologia nova (`exemplos/07-benchmark-sort/sort_bench_test.go`)

- Cada tamanho roda o **seu** tamanho: 5, 1k, 10k, 1M de verdade (subbenchmarks).
- Todos recebem input **desordenado e fresco a cada iteração**, e a restauração
  fica dentro da medição para os três — ninguém ordena o que já está ordenado.
- `testing.B.Loop()` (Go 1.24): setup uma vez, entradas/saídas vivas contra
  dead-code elimination.
- Evita-se `StopTimer/StartTimer` por iteração — para size 5 são dezenas de
  milhões de iterações e alternar o timer em cada uma é patologicamente lento.
  (Benchmark também é código que precisa de revisão.)
- Três abordagens: nativo in-place (reusa buffer), imutável via `slices.Clone`,
  lazy via `slices.Sorted(slices.Values(...))`.

## Os números de 2026

`make bench-sort` (Apple M3 Pro, go1.26.0, `-benchtime=300ms`):

| Benchmark | ns/op | B/op | allocs/op |
|--|--|--|--|
| `SortNative/5` | 11.8 | 0 | 0 |
| `SortNative/1k` | 11.201 | 0 | 0 |
| `SortNative/10k` | 297.051 | 0 | 0 |
| `SortNative/1M` | 69.020.225 | 0 | 0 |
| `SortImmutable/5` | 28.5 | 48 | 1 |
| `SortImmutable/1k` | 11.770 | 8.192 | 1 |
| `SortImmutable/10k` | 246.630 | 81.920 | 1 |
| `SortImmutable/1M` | 69.170.633 | 8.003.584 | 1 |
| `SortLazy/5` | 137.7 | 184 | 7 |
| `SortLazy/1k` | 15.401 | 25.272 | 15 |
| `SortLazy/10k` | 250.869 | 357.688 | 22 |
| `SortLazy/1M` | 72.582.483 | 41.678.203 | 41 |

> Números de uma execução; para a versão de slide, rode `make benchstat`
> (`-count=6`) e use o resumo com variância.

> ⚠️ **Medidos no go1.26.0, antes de o módulo subir para 1.27.** Uma execução de
> conferência no go1.27.0 deu números sensivelmente melhores no caso grande
> (`SortNative/1M` ~54,8M ns/op contra os 69,0M da tabela), mas as **conclusões
> não mudam**: nativo continua 0 alloc, imutável continua 1 alloc do tamanho do
> slice, lazy continua alocando mais. Antes do palco, re-rode `make benchstat` no
> 1.27 e substitua a tabela — não copie os números acima como se fossem atuais.

## A conclusão nova — honesta, não binária

- **O tempo agora escala com o tamanho.** 5 e 1k não são mais idênticos. O bug
  sumiu.
- **O custo da imutabilidade é uma alocação.** `SortImmutable` tem exatamente 1
  alloc (o clone, `8 bytes × n`); `SortNative` reusa buffer e tem 0.
- **Esse custo é ruído em escala, visível só no pequeno.** Em size 5 o imutável
  é ~2,4× mais lento (28,5 vs 11,8 ns). Em 1M o sort domina e os tempos
  convergem (69,0 vs 69,2 ms): a alocação vira detalhe.
- **A versão lazy aloca proporcionalmente ao tamanho** porque `slices.Sorted`
  faz `append` num slice que cresce — é o preço de coletar a sequência num
  slice. Ótimo para o palco: nem todo caminho "funcional" é grátis; é preciso
  medir.

E o contraponto que fecha o argumento de 2023 (ver `exemplos/03-iterators/iterators.go`): quando a
operação **não** precisa materializar (map/filter/reduce que colapsa num
escalar), o pipeline lazy aloca **zero** e é mais rápido:

| Benchmark | ns/op | B/op | allocs/op |
|--|--|--|--|
| `SumSquaresLazy` (iterators) | 5.223 | 0 | 0 |
| `SumSquaresSlices` (materializa) | 19.858 | 169.209 | 17 |

O custo que 2023 media não era "programação funcional". Era **copiar a
coleção**. Com iterators, esse custo pode simplesmente não existir.
