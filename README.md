# Programação Funcional com Go — Possibilidade ou forçação de barra?

Material de apoio da palestra de **Marco Ollivier**, apresentada na GopherCon
Brasil 2023 e reapresentada, em versão atualizada, em 2026.

Me siga no [![Instagram](https://img.shields.io/badge/Instagram-%23E4405F.svg?style=for-the-badge&logo=Instagram&logoColor=white)](https://instagram.com/marcopollivier)

## O que é este repositório

Uma coleção de exemplos pequenos e autocontidos — não uma biblioteca. Cada
arquivo sustenta um ponto da palestra: funções puras, imutabilidade,
first-class/higher-order functions, iterators, monads à mão, e um conjunto de
benchmarks que compara ordenação in-place vs. imutável.

**A tese de 2026:** em 2023 eu disse que programação funcional em Go era
forçação de barra. Em seis releases o Go corrigiu escopo de variável de loop
(1.22), ganhou iterators lazy (1.23), generic type aliases e `testing.B.Loop`
(1.24), `testing/synctest` (1.25) e métodos genéricos (1.27) — sem escrever a
palavra "funcional" em release note nenhum. A linguagem andou na direção do
paradigma. O que ela ainda não dá é higher-kinded types — e é aí que a
forçação de barra foi morar.

## Requisitos

- **Módulo raiz:** Go **1.27**. Dependência: `github.com/IBM/fp-go/v2`.
- **Submódulo `evolution-127/`:** também Go **1.27**. Continua separado não por
  versão, mas porque **redeclara** `Option[T]`/`Some`/`None`/`HalfIfEven` para
  mostrar o "depois" — os mesmos nomes que `exemplos/06-option/` usa para o
  "antes". O par só coexiste em pacotes distintos.

## Como rodar

```bash
go test ./...                      # todos os testes do módulo raiz
go test ./exemplos/03-iterators/   # só um exemplo
go vet ./...                       # análise estática
make bench-sort                    # o benchmark central (ordenação)
make benchstat                     # roda N vezes e resume (requer benchstat)
make evolution                     # exemplos do submódulo (métodos genéricos, 1.27)
```

Para ver o benchmark de 2023 (com o bug) e comparar com o de hoje:

```bash
git checkout v1-devpr24  # o código de 2023, preservado na tag
go test -bench=Sort ./...
git checkout main
```

## Índice — um pacote por ponto da palestra

Os exemplos ficam em [`exemplos/`](exemplos/), numerados **na ordem do palco**.
Cada pasta é autocontida: abra uma, leia o comentário do topo e rode os testes.
O roteiro de estudo detalhado está em [`exemplos/README.md`](exemplos/README.md).

| Pasta | Conceito | Release |
|--|--|--|
| [`01-paradigma`](exemplos/01-paradigma/) | Funções puras e function composition (`Compose`, `ComposeFn`) | — |
| [`02-loopvar`](exemplos/02-loopvar/) | Variável de loop por iteração (antes/depois) | 1.22 |
| [`03-iterators`](exemplos/03-iterators/) | `iter.Seq`, pipeline lazy map/filter/reduce | 1.23 |
| [`04-fpgo`](exemplos/04-fpgo/) | `Result`/`Either` com `IBM/fp-go` v2 | 1.24 (aliases) |
| [`05-concorrencia`](exemplos/05-concorrencia/) | Teste determinístico de concorrência | 1.25 |
| [`06-option`](exemplos/06-option/) | `Option[T]` à mão, operações como funções (o "antes") | — |
| [`07-benchmark-sort`](exemplos/07-benchmark-sort/) | Mutável vs imutável vs lazy + o post-mortem do bug de 2023 | 1.24 (`B.Loop`) |
| [`08-agentes`](exemplos/08-agentes/) | Puro vs estado compartilhado (custo de revisão) | — |
| [`evolution-127/option.go`](evolution-127/option.go) | `Option[T]` com **métodos genéricos** encadeáveis | 1.27 |
| [`evolution-127/ceiling_hkt.go`](evolution-127/ceiling_hkt.go) | O teto: por que não há higher-kinded types | 1.27 |

## Documentação

- [`docs/esqueleto-palestra.md`](docs/esqueleto-palestra.md) — o roteiro slide a slide (~40 min).
- [`docs/paradigma-e-agentes.md`](docs/paradigma-e-agentes.md) — paradigma na era dos agentes: os seis pilares.
- [`docs/o-que-mudou.md`](docs/o-que-mudou.md) — o registro técnico, release a release (1.22 → 1.27).
- [`docs/benchmark-2023-vs-2026.md`](docs/benchmark-2023-vs-2026.md) — o post-mortem do benchmark antigo.
- [`experiments/agent-eval/`](experiments/agent-eval/) — protocolo do experimento com agente.

---

## Submissão / Call for Papers

> Texto pronto para submeter em eventos. Versionado aqui para reuso.

### Título

**Programação Funcional com Go: possibilidade ou forçação de barra?**
_Três anos depois — o que a linguagem me deu, e o que ainda não vai dar._

### Abstract curto (~300 caracteres)

Em 2023 concluí que programação funcional em Go era forçação de barra. Seis
releases depois, o Go andou na direção do paradigma — loop var, iterators,
métodos genéricos — sem nunca dizer "funcional". Volto para mostrar, com código
e benchmarks, o que mudou, o que ainda trava, e por que isso importa na era dos
agentes.

### Abstract médio (~600 caracteres)

Em 2023 eu subi no palco da GopherCon BR e disse que programação funcional em
Go era, no fundo, forçação de barra. Seis releases depois, o Go corrigiu o
escopo da variável de loop, ganhou iterators lazy, generic type aliases e
métodos genéricos — andando na direção do paradigma sem escrever "funcional" em
release note nenhum. Nesta versão atualizada, reviso minhas próprias conclusões
com código que compila e benchmarks refeitos (incluindo o post-mortem de um
benchmark meu de 2023 que estava furado), mostro onde a linguagem ainda traça o
limite (higher-kinded types) e por que pureza e imutabilidade deixaram de ser
estilo pessoal para virar especificação legível por máquina na era dos agentes.

### Descrição para os organizadores

Esta é a continuação de uma palestra que apresentei na **GopherCon Brasil 2023**
e reapresentei em 2024. A versão de 2026 não repete o conteúdo: ela revisa as
conclusões à luz de tudo que mudou no Go entre a 1.21 e a 1.27.

O fio condutor é honesto e verificável: eu mostro um benchmark **meu**, de 2023,
que "provava" que mutabilidade era mais rápida — e demonstro, com o código na
tela, que ele estava metodologicamente furado (o benchmark foi corrompido pela
própria mutação in-place que ele media). É o gancho para o segundo eixo: na era
dos agentes, o gargalo virou **revisão**, não geração de código, e função pura é
desproporcionalmente mais barata de revisar e testar.

**Minhas credenciais:** staff engineer; palestrei sobre FP em Go na GopherCon
BR. Trabalhei com **Clojure em produção no Nubank** (arquitetura hexagonal — foi
onde o paradigma funcional virou a chave pra mim), com **Go e Node na Flash**, e
passei por OLX. Ajudo a organizar os meetups **GopheRIO** e **GoLangSP**.

**Outline (~40 min):**

- **(4 min) Alinhando os patos.** Funcional não resolve tudo; é mais um
  paradigma; toda escolha abre mão de algo. Não sigam hype. _Gancho:_ "e o
  agente não escreve isso pra você de qualquer jeito?" — pergunta plantada, sem
  resposta ainda.
- **(6 min) O paradigma.** Base matemática (`f(x) = x² + 3`), contrato/entrada/
  processamento/saída, evitar mudança de estado. Pureza e funções puras.
- **(12 min) O que o Go me deu desde 2023.** Loop var (1.22), iterators lazy
  (1.23), generic type aliases + `B.Loop` (1.24), `synctest` (1.25), métodos
  genéricos (1.27) — cada um com o exemplo rodando.
- **(8 min) O benchmark que me traiu.** Post-mortem ao vivo: o de 2023 vs o
  refeito. Iterators lazy alocando zero contra a versão que copia coleção.
- **(7 min) Paradigma na era dos agentes.** Verificação virou o gargalo;
  pureza, tipos e imutabilidade como spec de máquina. O contra-argumento do
  corpus de treino, sem esconder.
- **(3 min) Conclusão.** É viável usar funcional com Go — mais do que em 2023 —,
  mas a linguagem não foi pensada pra isso, e o teto agora é higher-kinded
  types: dá pra ter `Option.Map`, não dá pra abstrair sobre "todo Functor".

**A audiência sai sabendo:** usar iterators e `slices`/`maps` para escrever
map/filter/reduce lazy sem materializar coleções; escrever um `Option`/`Result`
encadeável com métodos genéricos (1.27); reconhecer quando imutabilidade custa
(e quando é ruído); escrever um benchmark honesto com `testing.B.Loop`; e
justificar, com argumento técnico, por que pureza e tipos explícitos aumentam o
rendimento de um agente.

### Nível e pré-requisitos

**Intermediário.** Espera-se Go básico (structs, slices, goroutines) e alguma
familiaridade com generics. Não é preciso ter visto a palestra de 2023 nem
conhecer programação funcional de antemão — os conceitos são introduzidos.
