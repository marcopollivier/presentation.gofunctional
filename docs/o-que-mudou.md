# O que mudou no Go entre 1.21 e 1.27

Registro técnico release a release, com foco no que afeta a viabilidade de
programação funcional. Fonte de consulta e referência de fim de palestra.

> Nenhuma dessas mudanças foi anunciada como "funcional". O argumento da
> palestra é justamente esse: a linguagem andou na direção do paradigma sem
> nomeá-lo.

---

## Go 1.22 — fev/2024 — variável de loop por iteração

A variável de um `for` passou a ser criada **a cada iteração**, não uma vez por
loop. O footgun clássico (closure capturando a variável compartilhada e todos
lendo o último valor) deixou de existir. Uma mudança de semântica de 12 anos,
feita para acomodar closures dentro de loops — o pré-requisito de map/filter/
reduce com funções.

- Exemplo no repo: `loopvar.go`
- Release notes: https://go.dev/doc/go1.22
- Discussão de design: https://go.dev/blog/loopvar-preview

## Go 1.23 — ago/2024 — iterators e range-over-func

**A mudança mais importante para esta palestra.** Chegaram `iter.Seq[V]` e
`iter.Seq2[K,V]`, o `range` sobre funções, e o suporte na stdlib: `slices.Values`,
`slices.Collect`, `slices.Sorted`, `slices.SortedFunc`, `slices.All`, `maps.Keys`,
`maps.Values`, `maps.All`. Isso dá map/filter/reduce **lazy e componível**, sem
materializar slices intermediários — atacando diretamente a conclusão do
benchmark de coleções de 2023, cujo custo medido era o de **copiar a coleção**.
Também entrou o pacote `unique` (interning), que conversa com semântica de valor
e imutabilidade.

- Exemplo no repo: `iterators.go`, `sort.go` (`SortLazy`)
- Release notes: https://go.dev/doc/go1.23
- Iterators: https://go.dev/blog/range-functions

## Go 1.24 — fev/2025 — generic type aliases e `testing.B.Loop`

- **Generic type aliases** destravaram bibliotecas de FP — é o que permitiu o
  `fp-go` v2, cujo `type Result[T] = Either[error, T]` só é legal a partir daqui.
- **`testing.B.Loop()`** mudou como se escreve benchmark: roda o setup uma vez,
  mantém parâmetros e resultados vivos (evita dead-code elimination) e conta as
  iterações sozinho.

- Exemplo no repo: `fpgo_example.go`, todos os benchmarks (`for b.Loop()`)
- Release notes: https://go.dev/doc/go1.24
- `B.Loop`: https://go.dev/blog/testing-b-loop

## Go 1.25 — ago/2025 — `testing/synctest`

Teste **determinístico** de código concorrente. Dentro de um "bubble" o relógio
é falso e só avança quando todas as goroutines estão bloqueadas — testes que
dependeriam de sleeps reais rodam instantâneos e sem flakiness. Reforça o
argumento "concorrência + testabilidade" que em 2023 estava só na percepção
pessoal. (O pacote foi experimental no 1.24, sob `GOEXPERIMENT=synctest`, e
estabilizou no 1.25 como `testing/synctest`, com a API `synctest.Test`.)

- Exemplo no repo: `concurrency_synctest_test.go`
- Release notes: https://go.dev/doc/go1.25

## Go 1.26 — fev/2026 — (o release que o brief pulou)

Não trouxe recurso de linguagem central para FP, mas vale registrar por causa
do custo relativo da abordagem imutável baseada em cópia:

- **Green Tea GC** habilitado por padrão.
- `new(expr)` — `new` passou a aceitar uma expressão como valor inicial.
- Tipos genéricos podem se referir a si mesmos na própria lista de parâmetros.
- Overhead de cgo ~30% menor; mais casos de alocação de backing array de slice
  na **stack**.

- Release notes: https://go.dev/doc/go1.26

## Go 1.27 — ago/2026 (lançado) — métodos genéricos

**O divisor de águas para a conclusão da palestra.** Um método pode declarar
seus próprios parâmetros de tipo. Antes, `func (o Option[A]) Map[B](f func(A) B) Option[B]`
era ilegal — e é por isso que as libs de FP obrigam `Pipe3(x, Map(f), Chain(g))`
em vez de `x.Map(f).Chain(g)`.

As duas limitações que **permanecem** e definem o teto:

- Métodos de interface **não podem** declarar parâmetros de tipo, nem métodos
  genéricos podem implementar métodos de interface.
- Métodos genéricos **não aparecem via reflection**.

Consequência conceitual — a nova conclusão: **Go ganhou monads concretos e
encadeáveis, mas continua sem higher-kinded types.** Dá para ter `Option.Map`;
não dá para escrever código que abstrai sobre "qualquer coisa que seja um
Functor". A forçação de barra saiu da sintaxe e foi para a capacidade de
abstração.

Também entrou **alocação especializada por tamanho**: algumas alocações abaixo
de 80 bytes ficam até ~30% mais baratas — mas as próprias release notes estimam
o impacto em programas reais em **~1%**. Nada de exagerar isso no palco.

- Exemplo no repo: `evolution-127/option.go`, `evolution-127/ceiling_hkt.go`
- Release notes: https://go.dev/doc/go1.27
- Proposta de métodos genéricos: https://go.dev/issue/49085

---

## Ecossistema

- **`IBM/fp-go`** saiu de work-in-progress para a **v2** (exige Go 1.24+, usa
  generic type aliases): `Option`, `Either`, `Result`, `IO`, `IOResult`,
  `Reader`, `ReaderIOResult`, optics (Lens, Prism, Traversal) e do-notation.
- **`fp-ts`** (citada na versão original por ser TypeScript) entrou em modo de
  manutenção. O sucessor de facto é o **Effect-TS** — Giulio Canti entrou no
  time em 2023 e o Effect 3.0 estabilizou a API em 2024.
- **`sort.Ints` foi superado por `slices.Sort`** (não é formalmente
  depreciado): desde o Go 1.22 `sort.Ints` apenas delega para `slices.Sort`.
  O caminho novo é genérico e sem boxing de interface.
