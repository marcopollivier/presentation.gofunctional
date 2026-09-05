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
- **Um pacote por ponto da palestra, numerado na ordem do palco.** `exemplos/`
  vai de `00-basico-sobre-ponteiros` a `08-agentes`; o número é a posição na narrativa, não
  uma hierarquia. Ao mexer na ordem dos slides, renumere junto.
- **Não apague o material de 2023.** O "antes" está na tag `v1-devpr24` e é a
  evidência do arco antes/depois. Onde o contraste é o conteúdo do slide, o
  "antes" aparece inline (ex.: `exemplos/02-loopvar/loopvar.go`).
- **Escreva em pt-BR.** Comentários, docs e mensagens de commit são em
  português — é o idioma da palestra. Identificadores em inglês, como em Go.

### Os contra-exemplos são conteúdo — não "conserte" nenhum

Boa parte deste repo é deliberadamente "errada": é o lado esquerdo do slide.
Um agente que aplica as regras de estilo abaixo sem ler isto vai refatorar o
material da palestra e destruir o contraste. **Nenhum destes deve ser
corrigido, unificado ou tornado puro:**

| Anti-exemplo | Por que existe |
|--|--|
| `08-agentes/reviewability.go` — `var accumulator` global | viola a regra 1 **de propósito**; é o slide de revisibilidade |
| `02-loopvar/loopvar.go` — `LoopVarShared` | reproduz o footgun pré-1.22 via ponteiro; deve devolver `[3 3 3]` |
| `07-benchmark-sort/sort.go` — `SortNative` muta in-place | a comparação entre as três variantes É o benchmark |
| `03-iterators/iterators.go` — `SumSquaresOfEvensSlices` | materializa slices para medir o custo que o bench de 2023 confundia com "o preço do paradigma" |
| `01-paradigma/composition.go` — `Compose` (eager) | a resposta de 2023; `ComposeFn` é a de hoje. Ambas ficam |
| `06-option/option_pre127.go` — funções de pacote | o "antes" do 1.27; a escada de variáveis em `HalfIfEven` é o ponto |
| `evolution-127/ceiling_hkt.go` | **não compila de propósito** (`//go:build ignore`) |

(Caminhos relativos a `exemplos/`, exceto o último.)

Em caso de dúvida, o comentário no topo de cada arquivo explica o papel dele.
Se um exemplo parece ingênuo ou "consertável", leia o comentário antes de mexer.

## Comandos

```bash
go test ./...                          # roda todos os testes (módulo raiz)
go test -run TestCompose ./...         # um teste por nome
go test -run 'TestSort.*' -v ./...     # um grupo por regex, verboso
go vet ./...                           # análise estática
make all                               # vet + test + bench-sort (verificação completa)
make bench                             # todos os benchmarks com alocação
make bench-sort                        # só o benchmark de ordenação (achado central)
make benchstat                         # roda N vezes e resume com benchstat
make evolution                         # exemplos do submódulo (métodos genéricos)

make labs                              # testes de todos os labs (exemplos reais)
make lab LAB=020-nome                  # um lab só (TARGET=run para executar)
make labs-vet                          # go vet em cada lab
make lab-new NAME=021-nome             # cria um lab (GOVERSION=1.21 para pinar)
```

Os alvos de bench aceitam `BENCHTIME` (default `300ms`) e `COUNT` (default `6`):
`make benchstat COUNT=10`. `make benchstat` exige
`go install golang.org/x/perf/cmd/benchstat@latest` e escreve `bench.txt` na
raiz — artefato local, não commite (não há `.gitignore` no repo).

Para demonstrar o teto de HKT no palco (o erro do compilador **é** o slide):

```bash
cd evolution-127 && go build ceiling_hkt.go
# interface method must have no type parameters
```

**Módulo raiz:** Go 1.27 (`go.mod`, módulo `gofunctional`). Depende de
`github.com/IBM/fp-go/v2`. **Submódulo `evolution-127/`:** também Go 1.27. Os
dois módulos são separados: `go test ./...` na raiz **não** alcança
`evolution-127/` — use `make evolution`.

**Labs (`labs/`):** cada lab é mais um módulo isolado, com go.mod e versão de Go
próprios. `go test ./...` e `make all` na raiz **não** alcançam nenhum deles —
use `make labs`. É de propósito: um lab pinado no 1.21, ou com uma dep quebrada,
não pode derrubar o material curado.

A separação **não** é mais por versão (o 1.27 saiu em ago/2026 e a raiz subiu
junto): `evolution-127/` redeclara `Option[T]`, `Some`, `None` e `HalfIfEven`
para mostrar o "depois" com métodos genéricos, e os mesmos nomes estão em
`exemplos/06-option/` como o "antes". O par só coexiste em pacotes separados —
a colisão é a demonstração. Não funda os dois.

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

**Módulo raiz (Go 1.27)** — um pacote por ponto da palestra, em `exemplos/`,
numerados na ordem do palco. O índice de estudo é `exemplos/README.md`.

- `00-basico-sobre-ponteiros/` (`package ponteiros`) — o slice como cabeçalho
  passado por valor sobre um array compartilhado. `SortedPrint` (muta o
  argumento) e `SafeSortedPrint` (`slices.Clone` antes de ordenar) têm a mesma
  assinatura e imprimem o mesmo — só o que sobra no chamador difere. A demo de
  palco é um `Example` com saída verificada, não um `main`.
- `01-paradigma/` (`package paradigma`) — funções puras e function composition.
  `Compose` (versão eager de 2023, preservada como didática) e `ComposeFn` (a
  composição matemática de verdade, `g∘f`, lazy).
- `02-loopvar/` (`package loopvar`) — Go 1.22, variável de loop por iteração
  (antes/depois inline).
- `03-iterators/` (`package iterators`) — Go 1.23, `iter.Seq` e pipeline lazy
  map/filter/reduce; a versão lazy vs a que materializa slices (ataca a
  conclusão do bench de coleções de 2023).
- `04-fpgo/` (`package fpgo`) — `Result`/`Either` com `IBM/fp-go` v2 (estilo
  `Pipe`).
- `05-concorrencia/` (`package concorrencia`) — Go 1.25, teste determinístico
  de código concorrente (`testing/synctest`). Só tem arquivo `_test.go`, de
  propósito: o exemplo **é** o teste.
- `06-option/` (`package option`) — `Option[T]` à mão, operações como **funções
  de pacote** (o "antes" do 1.27, sem encadeamento).
- `07-benchmark-sort/` (`package ordenacao`) — contraste de mutabilidade e o
  benchmark. `SortNative` (muta in-place, via `slices.Sort`), `SortImmutable`
  (`slices.Clone` + sort), `SortLazy` (`slices.Sorted(slices.Values(...))` — o
  "sort imutável" de 2023 agora numa linha da stdlib). **Mantenha as três
  variantes** — a comparação é o conteúdo. O `sort_bench_test.go` traz a
  metodologia correta (`testing.B.Loop`, tamanhos reais, input fresco por
  iteração) e o comentário no topo documenta por que o de 2023 estava furado.
- `08-agentes/` (`package agentes`) — duas implementações da mesma lógica, pura
  vs estado compartilhado, para o contraste de revisibilidade no palco.

O nome do diretório tem número e hífen; o do pacote, não (Go não permite). Essa
divergência é intencional — o número serve à narrativa, não ao compilador.

**Submódulo `evolution-127/` (Go 1.27):**

- `option.go` — o mesmo `Option[T]`, agora com **métodos genéricos**
  encadeáveis (`Some(n).Chain(f).Map(g).GetOrElse(d)`).
- `ceiling_hkt.go` — não compila de propósito (`//go:build ignore`): por que
  Go tem monads concretos mas não tem higher-kinded types.

**Labs (`labs/`) — exemplos reais, um módulo isolado por lab:**

Onde `exemplos/` é o trecho curado que cabe num slide, `labs/` é o exemplo que
**roda**: pode ter `main`, dependência externa e uma versão de Go diferente da
raiz — porque muitas vezes é justamente isso que ele demonstra. Cada lab tem
`go.mod`, `Makefile` e `README.md` próprios, e o `Makefile` de `labs/` descobre
os labs sozinho (`wildcard */go.mod`): criar a pasta basta, não há lista para
manter. `make lab-new NAME=CEE-slug [GOVERSION=1.21]` gera a partir de
`labs/_template/` (que não é um lab — o `_` faz o Go ignorar a pasta).

A numeração é `CEE-slug`: capítulo e exemplo dentro do capítulo — `020` é o
primeiro exemplo do capítulo 2, `021` o segundo. O dígito extra existe para um
capítulo crescer sem renumerar os outros. Fora isso, duas regras herdadas de
`exemplos/`: o diretório tem número e hífen, o pacote não;
e **o comentário de topo é obrigatório**. Labs são material de palco pelas mesmas
razões — se um parecer ingênuo, errado ou "consertável", leia o comentário antes
de mexer; provavelmente ele é o lado esquerdo de um slide. `labs/` **não**
substitui `exemplos/`: as duas pastas coexistem e servem a coisas diferentes.

**Docs:** `docs/` (esqueleto da palestra, paradigma-e-agentes, o-que-mudou,
post-mortem do benchmark) e `experiments/agent-eval/` (protocolo do experimento
com agente — resultados são preenchidos pelo autor, nunca fabricados).
