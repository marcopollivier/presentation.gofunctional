# Benchmark cross-version — mutável vs. imutável, do Go 1.21 ao 1.27

Refaz a medição que sustentou a conclusão de 2023 ("imutabilidade custa caro"),
corrige o erro metodológico daquele benchmark e mede como o resultado evoluiu
em sete releases do Go.

## Por que este módulo é isolado

É um módulo Go próprio, com `go.mod` separado do módulo raiz e dos `labs/`.
Duas razões:

1. **`go test ./...` na raiz não o alcança.** Uma medição que leva dezenas de
   minutos não pode entrar no caminho da verificação de palco.
2. **A diretiva `go` precisa ser 1.21**, e a raiz está em 1.27.

### Por que `go 1.21` no `go.mod`

Aquela linha **não** é "versão mínima do toolchain": é a versão da
**linguagem**. Ela decide o que o compilador aceita e com que semântica.

Fixando em 1.21, o mesmo fonte compila com semântica idêntica nos sete
toolchains, e a única variável entre uma coluna e outra passa a ser
compilador + runtime + GC. Se dissesse `go 1.27`, o toolchain 1.21 recusaria
compilar e a matriz morreria na primeira coluna.

Não há diretiva `toolchain` de propósito — quem escolhe a versão é o `run.sh`,
por rodada, via `GOTOOLCHAIN`.

Consequência prática: **nenhuma API posterior ao 1.21** no código comparável.
Nada de `slices.Sorted` (1.23), `iter` (1.23) ou `testing.B.Loop` (1.24). O
vocabulário é `sort.Ints`, `make`, `copy` e o laço clássico
`for i := 0; i < b.N; i++` — que é exatamente o vocabulário de 2023, e é o que
torna a comparação legítima.

> `sort.Ints` **não está deprecada** — não há marcador `Deprecated:` no doc, nem
> no 1.27. Desde o 1.22 ela apenas chama `slices.Sort` internamente. Essa troca
> de implementação é parte do que a matriz mede.

## A matriz: dois eixos

Trocar só a versão do Go mediria a soma de duas coisas — o que o Go melhorou
**e** o viés do benchmark antigo — e o número não significaria nada. Então:

**{código de 2023 com o erro preservado, medição corrigida} × {Go 1.21 … 1.27}**

- a linha *legacy* em todas as versões isola **o artefato**
- a linha corrigida em todas as versões isola **a evolução real**
- a diferença entre as duas é, ela mesma, um resultado

### O erro de 2023

`sort.Ints` ordena in-place. A partir da **segunda** iteração do loop o slice já
está ordenado — o lado mutável passou a competir com a melhor entrada possível,
enquanto o lado imutável recebia dados aleatórios sempre, porque nunca tocava no
original.

Não foi descuido de medição: foi **não-localidade**. `func SortInPlace(s []int)`
não tem `*`, não tem anotação, não tem nada que avise que o argumento será
destruído. O erro estava na assinatura, não no laço.

## Como rodar

```bash
./run.sh install    # pré-aquece os sete toolchains (~100 MB cada, uma vez)
./run.sh run        # a matriz — um arquivo por versão em results/
./run.sh compare    # as tabelas do benchstat
./run.sh all        # os três, em ordem
```

Sobrescrevíveis por ambiente:

```bash
COUNT=1 BENCHTIME=100ms ./run.sh run    # smoke, minutos
BENCH='InPlace' ./run.sh run            # um recorte
RESULTS=/tmp/x ./run.sh run             # noutro diretório
```

Os testes, que valem mais que os benchmarks, rodam sozinhos:

```bash
go test -v -run 'TestLegacy|TestImmutable' .
```

### Condições de execução — não é preciosismo

Sete toolchains em sequência levam dezenas de minutos, e um notebook esquentando
no meio faz atribuir ao Go o que é do ventilador.

- máquina **ociosa**: sem build, sem Docker, sem browser pesado
- **na tomada** (em bateria o macOS reduz frequência)
- atenção ao **throttling térmico**: se a última versão da lista parecer
  sistematicamente pior, desconfie da temperatura antes de desconfiar do Go —
  e rode a lista na ordem inversa para confirmar
- `-cpu=1` já remove a variação de escalonamento
- `-count=10` é o que dá ao benchstat repetições para estimar variância; com
  `count=1` ele imprime `± ∞` e **nenhuma comparação tem significância**

O `results/environment.txt` registra máquina, toolchains e parâmetros. Vai
versionado: sem ele, um resultado guardado por seis meses não é auditável.

## Como ler os resultados

### Os benchmarks

| benchmark | o que mede |
|---|---|
| `LegacyInPlace` | o benchmark de 2023, com o viés preservado |
| `LegacyImmutable` | o par dele — nunca foi enviesado, e é por isso que serve de controle |
| `InPlace` | custo puro do sort (restauração fora do cronômetro) |
| `InPlaceCopy` | sort + cópia, buffer reaproveitado, zero alocação |
| `Immutable` | sort + cópia + alocação nova |
| `CopyOnly` | só a cópia, sem sort — a testemunha |
| `TimerPair` | `sem`/`com` o par `StopTimer/StartTimer`, para quantificar o viés do instrumento |

### As subtrações

```
Immutable − InPlaceCopy    = custo isolado da ALOCAÇÃO
Immutable − InPlace        = custo total da IMUTABILIDADE
InPlaceCopy − CopyOnly     = custo puro do sort, SEM cronômetro manipulado
TimerPair(com) − (sem)     = resíduo do par que entra na conta, por iteração
```

`InPlace` e `InPlaceCopy − CopyOnly` medem a mesma coisa por caminhos
independentes. **Em n ≥ 1k elas têm que convergir**; se não convergirem, uma
das duas está errada — e isso é conteúdo.

### Números que precisam de ressalva

**`InPlace` em n=5.** `b.StopTimer()`/`b.StartTimer()` custam tempo, e esse
custo cai em cima de cada iteração. Com o sort inteiro na casa de dezenas de
nanossegundos, o overhead é da mesma ordem do que se quer medir. **A partir de
1k é desprezível.** Use `InPlaceCopy − CopyOnly` para o n=5.

**Sobre o instrumento, que rendeu duas descobertas:**

1. Um benchmark cujo corpo é só o par `StopTimer/StartTimer` **não termina**. O
   tempo cronometrado por iteração é ~zero — o cronômetro está parado
   justamente durante o que se quer medir — então o framework sobe `b.N`
   perseguindo o `-benchtime` e o relógio de parede explode.
2. O motivo do custo é maior que "chamada de função": em
   `$GOROOT/src/testing/benchmark.go`, **`StopTimer` e `StartTimer` chamam
   `runtime.ReadMemStats`, que é stop-the-world**. Cada par para o mundo duas
   vezes, e esse custo fica *fora* da conta (em `StopTimer` o `ReadMemStats` vem
   depois de fechar a duração; em `StartTimer`, antes de reabrir).

Por isso `BenchmarkTimerPair` é diferencial e usa carga de 1k: é a única forma
que converge.

## Resultados

### Achados provisórios (smoke: `COUNT=1`, `BENCHTIME=50ms`)

Sem significância estatística — `count=1`. Servem para dizer o que procurar na
rodada boa, não para ir ao slide.

**1. O viés do benchmark de 2023, isolado (Go 1.27.1):**

```
             │ legacy  │ corrigido │
InPlace/5    │  6.01n  │   37.67n  │
InPlace/1k   │ 578.3n  │  9874.0n  │   ← 17× de subestimação
Immutable/5  │ 22.55n  │   22.30n  │
Immutable/1k │  9.18µ  │    9.38µ  │   ← praticamente igual, como tem que ser
```

O lado imutável não muda entre as duas versões do benchmark — ele nunca foi
enviesado. **Todo o erro está do lado mutável**, que é justamente o lado cuja
assinatura não avisava nada.

**2. O sintoma esperado NÃO se reproduziu — e isso é resultado.**

A expectativa era `ns/op` quase idêntico entre n=5 e n=1000 no `LegacyInPlace`.
Não é o que acontece:

```
LegacyInPlace   n=5  6.01n   n=1k  578.3n   →  96×
InPlace         n=5 37.67n   n=1k  9874n    →  262×
```

O viés não *achata* a curva: ele a torna **linear** em vez de O(n log n) —
ordenar dados já ordenados é O(n) no pdqsort. Noventa e seis vezes não é "quase
idêntico". Ou seja: **o `ns/op` plano dos slides de 2023 não é explicado por
este erro**, e a hipótese alternativa que você levantou fica mais forte — a de
que os quatro slides mostram o mesmo caso com rótulo trocado, dado o
`-bench=. -bench=Sort5` duplicado no comando. São dois erros distintos, não um.

**3. O sort ficou 2,7× mais rápido do 1.21 para o 1.22:**

```
InPlaceCopy/1k    1.21: 24000 ns   1.22: 8963 ns   …   1.27: 8919 ns
```

O salto está exatamente na release em que `sort.Ints` passou a chamar
`slices.Sort` — pdqsort genérico, sem dispatch de interface por comparação. O
benchmark de 2023 mediu um sort 2,7× mais lento que o de hoje, e a melhoria veio
de **generics na stdlib**.

**4. Em 1.21, o caminho "mutável" alocava:**

```
InPlace/5    1.21: 24 B/op, 1 alloc     1.22+: 0 B/op, 0 allocs
```

`sort.Ints` no 1.21 passava por `sort.Sort(IntSlice(x))`, e a conversão para a
interface escapava. O contraste "mutável não aloca, imutável aloca" de 2023 era,
em parte, falso.

**5. `Immutable − InPlaceCopy` (custo da alocação):** com `count=1` o ruído é
maior que o efeito. Fica para a rodada boa — é onde se espera ver Green Tea GC
(1.26) e alocação especializada por tamanho (1.27).

### Tabela final para slide

<!-- Cole aqui a tabela formatada depois de rodar com COUNT=10 BENCHTIME=1s.
     Sugestão de recorte, por tamanho:
       custo puro do sort        InPlaceCopy − CopyOnly
       custo da alocação         Immutable − InPlaceCopy
       custo da imutabilidade    Immutable − InPlace
     e a linha do artefato:      LegacyInPlace vs InPlace -->

_(a preencher)_

## Arquivos

| arquivo | papel |
|---|---|
| `sortbench.go` | núcleo comparável: as quatro operações, o PRNG, os tamanhos |
| `legacy_bench_test.go` | o benchmark de 2023, **com o erro preservado** — não conserte |
| `fixed_bench_test.go` | a medição corrigida e as testemunhas |
| `sortbench_test.go` | os dois testes que demonstram o erro sem cronômetro |
| `run.sh` | `install` / `run` / `compare` / `all` |
| `results/` | saída por versão, `environment.txt` e as comparações |
