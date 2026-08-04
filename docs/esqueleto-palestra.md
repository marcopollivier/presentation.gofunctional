# Esqueleto da palestra — versão 2026 (~40 min)

Outline slide a slide. Para cada bloco: **[mantém]** o que vem de 2023,
**[muda]** o que foi revisto, **[inédito]** o que é novo, e o arquivo/benchmark
deste repo que sustenta o slide.

**Tese:** em 2023 eu disse que FP em Go era forçação de barra. Seis releases
depois o Go andou na direção do paradigma — sem nunca dizer "funcional". Aqui
está o que ele me deu, o que ainda não vai dar (higher-kinded types), e por que
isso importa na era dos agentes.

---

## Bloco 0 — Abertura e gancho (2 min)

- Quem sou eu; a palestra de 2023 na GopherCon BR e a conclusão de então.
- **[inédito] O gancho plantado:** "em 2026 alguém vai perguntar: pra que pensar
  em paradigma se o agente escreve o código pra mim?" — pergunta feita, resposta
  guardada para o Bloco 6.

## Bloco 1 — Alinhando os patos (4 min)

- **[mantém]** Funcional não resolve todos os seus problemas; é só mais um
  paradigma; toda escolha implica abrir mão de outra coisa; não sigam hype.
- _Sustenta:_ o tom honesto que o resto da palestra cobra de si mesma.

## Bloco 2 — O paradigma (6 min)

- **[mantém]** Base matemática: `f(x) = x² + 3`. Contrato / entrada /
  processamento / saída. Evitar mudança de estado.
- **[mantém]** Funções puras e imutabilidade como definição.
- _Sustenta:_ `composition.go` (`AddOne`, `Double`, `ComposeFn` como `g∘f`).

## Bloco 3 — O que o Go me deu desde 2023 (12 min)

O coração da versão nova. Um exemplo rodando por release.

- **[inédito] Loop var (1.22).** O footgun de 12 anos que sumiu. Antes/depois.
  - _Sustenta:_ `loopvar.go` — `[0 1 2]` vs `[3 3 3]`.
- **[muda] Coleções (1.23).** Em 2023 esse era "o primeiro problema": map/filter
  custava copiar a coleção. Iterators lazy mudam a conclusão.
  - _Sustenta:_ `iterators.go` + benchmark `SumSquares*` (lazy: 0 alloc; slices:
    169 KB / 17 allocs).
- **[inédito] fp-go v2 (1.24).** Generic type aliases destravaram a lib.
  `Result[T] = Either[error, T]`. Estilo `Pipe`.
  - _Sustenta:_ `fpgo_example.go`.
- **[muda] Concorrência (1.25).** Em 2023 "concorrência combina com funcional"
  era percepção pessoal. `synctest` torna testável e determinístico.
  - _Sustenta:_ `concurrency_synctest_test.go`.
- **[inédito] Métodos genéricos (1.27).** `Option` à mão vira encadeável.
  - _Sustenta:_ `option_pre127.go` (antes, `Pipe`) vs
    `evolution-127/option.go` (depois, `x.Chain(f).Map(g)`).

## Bloco 4 — O benchmark que me traiu (8 min)

- **[muda] Post-mortem ao vivo.** Os números de 2023 (5 == 1000) e as três
  causas. A causa raiz: `sort.Ints` muta in-place, o "nativo" reordenava o que
  já estava ordenado.
- **[inédito] O reframe:** o benchmark que "provava" mutabilidade foi corrompido
  por mutabilidade. Classe de erro inexpressável na versão pura. **Ponte para o
  Bloco 6.**
- _Sustenta:_ `sort_bench_test.go`, `docs/benchmark-2023-vs-2026.md`,
  `git checkout v1-devpr24` para mostrar o código antigo.

## Bloco 5 — O teto (4 min)

- **[inédito]** Go ganhou monads concretos e encadeáveis (1.27), mas continua
  sem higher-kinded types. Dá pra ter `Option.Map`; não dá pra abstrair sobre
  "todo Functor".
- _Sustenta:_ `evolution-127/ceiling_hkt.go` — o erro do compilador ao vivo:
  `interface method must have no type parameters`.

## Bloco 6 — Paradigma na era dos agentes (7 min)

- **[inédito] A resposta ao gancho do Bloco 0.** As duas coisas se potencializam
  — se você fizer a sua parte. Os seis pilares, condensados:
  1. verificação virou o gargalo; pureza barateia a revisão;
  2. testabilidade é o loop de feedback do agente;
  3. tipos são o revisor que não cansa;
  4. imutabilidade limita o raio de dano;
  5. composição dá tarefas do tamanho certo;
  6. o contra-argumento honesto: o corpus de treino é imperativo, então **você**
     tem que escrever as regras (`CLAUDE.md`, linters, testes).
- **[inédito]** O `CLAUDE.md` deste repo projetado: o paradigma virou spec de
  máquina.
- _Sustenta:_ `reviewability.go` (contraste no palco), `CLAUDE.md`,
  `docs/paradigma-e-agentes.md`, `experiments/agent-eval/`.

## Bloco 7 — Conclusão (3 min)

- **[muda]** É viável usar funcional com Go — **mais** do que em 2023 —, mas a
  linguagem não foi pensada pra isso. A forçação de barra não sumiu: saiu da
  sintaxe e foi para a capacidade de abstração (HKT).
- **[inédito]** O paradigma não perdeu função na era dos agentes; mudou de
  função. Deixou de ser estilo e virou especificação.
- Referência de fim: `docs/o-que-mudou.md`.

---

### Distribuição de tempo

| Bloco | Min | Acumulado |
|--|--|--|
| 0 — Abertura e gancho | 2 | 2 |
| 1 — Alinhando os patos | 4 | 6 |
| 2 — O paradigma | 6 | 12 |
| 3 — O que o Go me deu | 12 | 24 |
| 4 — O benchmark que me traiu | 8 | 32 |
| 5 — O teto | 4 | 36 |
| 6 — Paradigma e agentes | 7 | 43 |
| 7 — Conclusão | 3 | 46 |

> ~46 min com folga; para uma faixa de 40 min, comprimir o Bloco 3 (menos tempo
> por release) e enxugar o Bloco 1.
