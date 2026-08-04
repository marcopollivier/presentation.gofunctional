# CLAUDE.md

Guia para o Claude Code (claude.ai/code) trabalhar neste repositório.

> Este arquivo tem papel duplo. Além de orientar o trabalho no repo, ele é
> **conteúdo de palestra**: a seção "Regras de estilo" é projetada no palco
> como exemplo de "o paradigma virou especificação legível por máquina".
> Escreva pensando que ele vai para o telão.

## Propósito

Material de apoio da palestra **"Programação Funcional com Go — Possibilidade
ou forçação de barra?"** (GopherCon BR 2023/2024; reapresentada em 2026).
Cada arquivo é um exemplo pequeno e autocontido explorando se idiomas de
programação funcional são viáveis em Go. É **material didático, não uma
biblioteca nem uma aplicação**. Não há `main`; tudo é exercitado por testes e
benchmarks.

Isso muda os critérios de qualidade:

- **Clareza didática > abstração.** Um exemplo que cabe num slide e é lido em
  voz alta vale mais que um genérico e "elegante".
- **Cada exemplo é autocontido, curto e projetável.** Nada de over-engineering.
- **Nomes de arquivo mapeiam seções da palestra** (`loopvar.go`, `iterators.go`,
  `sort.go`, `option_pre127.go`, `fpgo_example.go`, `reviewability.go`).
- **Não apague o material de 2023.** O "antes" está na tag `v1-devpr24` e é a
  evidência do arco antes/depois. Onde o contraste é o conteúdo do slide, o
  "antes" aparece inline (ex.: `loopvar.go`).

## Comandos

```bash
go test ./...                          # roda todos os testes (módulo raiz)
go test -run TestCompose ./...         # um teste por nome
go vet ./...                           # análise estática
make bench                             # todos os benchmarks com alocação
make bench-sort                        # só o benchmark de ordenação (achado central)
make benchstat                         # roda N vezes e resume com benchstat
make evolution                         # exemplos Go 1.27 (baixa o toolchain RC)
```

**Módulo raiz:** Go 1.26 (`go.mod`, módulo `gofunctional`). Depende de
`github.com/IBM/fp-go/v2`. **Submódulo `evolution-127/`:** Go 1.27 (RC),
isolado de propósito — rode com `GOTOOLCHAIN=auto`.

## Regras de estilo (spec do projeto — e slide)

O paradigma, aqui, não é preferência pessoal; é especificação. Estas regras
tornam o código barato de revisar e trivial de testar — que é o que faz um
agente (ou um humano sob volume) render:

1. **Pureza por padrão.** Uma função deve depender só dos seus argumentos e
   devolver só pelo retorno. Sem estado global, sem I/O escondido. Você lê a
   assinatura e sabe o universo do que pode acontecer.
2. **Imutabilidade em coleções.** Não mute o slice/map do chamador. Copie
   (`slices.Clone`) ou produza um novo (`slices.Sorted(slices.Values(...))`).
   A exceção é explícita no nome (ex.: `SortNative` diz que muta).
3. **Efeito colateral só na borda.** I/O, relógio, aleatoriedade e mutação
   ficam nas pontas; o miolo é puro e determinístico.
4. **Erro explícito no tipo de retorno.** `(T, error)` ou um `Result[T]`. Nunca
   um panic para controle de fluxo. O tipo é o revisor que não cansa: força o
   caminho de falha a aparecer.
5. **Composição em unidades pequenas.** Funções pequenas e componíveis, com
   contrato claro, em vez de um método de 300 linhas com sete responsabilidades.
6. **Teste sem cerimônia.** Se um exemplo precisa de mock ou fixture para ser
   testado, ele provavelmente violou (1)–(3). Prefira table-driven direto.

Contra-argumento honesto que a palestra não esconde: um modelo treinado em
corpus de Go produz o Go idiomático mediano — imperativo e mutável, porque é o
que existe em quantidade. Para ter as garantias acima com um agente, **você tem
que escrever as regras** — este arquivo, os testes, o `go vet`. A disciplina
não deixou de ser necessária; ela mudou de função.

## Estrutura

**Módulo raiz (Go 1.26):**

- `composition.go` — funções puras e function composition. `Compose` (versão
  eager de 2023, preservada como didática) e `ComposeFn` (a composição
  matemática de verdade, `g∘f`, lazy).
- `sort.go` — contraste de mutabilidade, base do benchmark. `SortNative`
  (muta in-place, via `slices.Sort`), `SortImmutable` (`slices.Clone` + sort),
  `SortLazy` (`slices.Sorted(slices.Values(...))` — o "sort imutável" de 2023
  agora numa linha da stdlib). **Mantenha as três variantes** — a comparação é
  o conteúdo.
- `sort_bench_test.go` — benchmark refeito com metodologia correta
  (`testing.B.Loop`, tamanhos reais, input fresco por iteração). O comentário
  no topo documenta por que o benchmark de 2023 estava furado.
- `loopvar.go` — Go 1.22, variável de loop por iteração (antes/depois inline).
- `iterators.go` — Go 1.23, `iter.Seq` e pipeline lazy map/filter/reduce; a
  versão lazy vs a que materializa slices (ataca a conclusão do bench de
  coleções de 2023).
- `option_pre127.go` — `Option[T]` à mão, operações como **funções de pacote**
  (o "antes" do 1.27, sem encadeamento).
- `fpgo_example.go` — `Result`/`Either` com `IBM/fp-go` v2 (estilo `Pipe`).
- `concurrency_synctest_test.go` — Go 1.25, teste determinístico de código
  concorrente (`testing/synctest`).
- `reviewability.go` — duas implementações da mesma lógica, pura vs estado
  compartilhado, para o contraste de revisibilidade no palco.

**Submódulo `evolution-127/` (Go 1.27):**

- `option.go` — o mesmo `Option[T]`, agora com **métodos genéricos**
  encadeáveis (`Some(n).Chain(f).Map(g).GetOrElse(d)`).
- `ceiling_hkt.go` — não compila de propósito (`//go:build ignore`): por que
  Go tem monads concretos mas não tem higher-kinded types.

**Docs:** `docs/` (esqueleto da palestra, paradigma-e-agentes, o-que-mudou,
post-mortem do benchmark) e `experiments/agent-eval/` (protocolo do experimento
com agente — resultados são preenchidos pelo autor, nunca fabricados).
