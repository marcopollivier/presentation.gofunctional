# Paradigma na era dos agentes

O segundo eixo da palestra. Não é marketing de IA nem defesa nostálgica de
artesanato — é argumento técnico verificável. A pergunta plantada no começo da
palestra é a que qualquer pessoa faz em 2026:

> "Por que eu preciso pensar em paradigma, pureza e imutabilidade? O Claude Code
> escreve o código pra mim do jeito que ele achar melhor."

A resposta, e o arco que fecha a palestra: **as duas coisas não competem, elas
se potencializam — mas só se você fizer a parte que continua sendo sua.**

---

## Os seis pilares

### 1. O gargalo virou verificação, não geração

Escrever código ficou barato; revisar não. Função pura é desproporcionalmente
mais barata de revisar — você lê a assinatura e sabe o universo do que pode
acontecer. Função com efeito colateral implícito exige que você leia o resto do
sistema para saber se está certa. Quando o volume de código gerado sobe, o custo
de revisão vira o fator dominante, e pureza é redução direta desse custo.

_No repo:_ `exemplos/08-agentes/reviewability.go` — `MovingAveragePure` vs `MovingAverageStateful`,
a mesma lógica, uma revisável isoladamente e a outra não.

### 2. Testabilidade é o loop de feedback do agente

Agente rende quando tem sinal automático — compilador e teste. Função pura é
trivialmente testável: sem mock, sem fixture, table-driven direto. Isso é o que
permite o agente se autocorrigir. Código com estado compartilhado exige mock e
setup; o agente erra mais e você percebe menos.

_No repo:_ `reviewability_test.go` (table-driven sem setup) e
`exemplos/05-concorrencia/concurrency_synctest_test.go` (concorrência testada de forma determinística).

### 3. Tipos são o revisor que não cansa

Erro explícito no tipo de retorno (`(T, error)`, `Result[T]`) faz o agente não
conseguir "esquecer" o caminho de falha. É guardrail de máquina, não convenção
social.

_No repo:_ `exemplos/01-paradigma/composition.go` (erro no retorno), `exemplos/04-fpgo/fpgo_example.go` e
`exemplos/06-option/option_pre127.go` / `evolution-127/option.go` (o erro/ausência vive no tipo).

### 4. Imutabilidade limita o raio de dano

Se o agente altera uma função pura, o impacto é local. Se altera uma struct
mutável compartilhada, o impacto é o programa. Isso sempre foi verdade — só
passou a importar muito mais quando o número de mudanças por dia multiplicou.

_No repo:_ `exemplos/07-benchmark-sort/sort.go` (variantes que copiam vs a que muta) e o próprio
post-mortem do benchmark.

### 5. Composição dá tarefas do tamanho certo

Agente vai melhor com escopo pequeno e contrato claro. Arquitetura feita de
funções pequenas e componíveis produz isso naturalmente; um método de 300 linhas
com sete responsabilidades, não.

_No repo:_ `exemplos/03-iterators/iterators.go` (pipeline de funções pequenas), `exemplos/01-paradigma/composition.go`.

### 6. O contra-argumento honesto — o corpus de treino

Um modelo treinado em corpus de Go produz o Go idiomático **mediano**:
imperativo e mutável, porque é o que existe em quantidade. Ou seja, se você quer
as garantias acima, **você tem que escrever as regras**: `CLAUDE.md`, linters,
testes, revisão. A disciplina que supostamente "não é mais necessária" é
exatamente o que faz o agente render. O paradigma não perdeu função — mudou de
função. Deixou de ser estilo pessoal e virou especificação legível por máquina.

_No repo:_ o próprio `CLAUDE.md`, escrito como spec projetável.

---

## A prova mais forte está neste repositório

O bug do benchmark de 2023 (ver `benchmark-2023-vs-2026.md`) existe porque
`sort.Ints` muta in-place. É exatamente o tipo de erro que passa despercebido
numa revisão de código imperativo — e que **não teria como existir se a
ordenação fosse pura**. A classe de erro é inexpressável na versão pura. Não é o
agente que resolve isso; é o paradigma que a torna impossível de escrever.

---

## Protocolo do experimento reproduzível

O bloco não pode ser opinião no slide. O experimento abaixo gera dado
verificável. **Estrutura aqui; os resultados são preenchidos pelo autor rodando
de verdade — nunca fabricados.** Ver `experiments/agent-eval/`.

**A tarefa** (`experiments/agent-eval/task.md`): a mesma pedida a um agente em
duas condições:

- **Condição A — sem restrição de estilo.** Prompt cru, sem `CLAUDE.md`, sem
  apontar os testes.
- **Condição B — com as regras.** Mesmo prompt, com o `CLAUDE.md` deste repo no
  contexto e os testes existentes disponíveis como sinal.

**Critérios objetivos** (`experiments/agent-eval/criteria.md`), medidos nas duas
saídas:

| Critério | Como medir |
|--|--|
| Nº de efeitos colaterais | Contar estado global mutável, I/O escondido, mutação de argumento |
| Testável sem mock? | O teste gerado usa mock/fixture ou é table-driven direto? |
| Alocações | `-benchmem` na função equivalente |
| Aparece o bug de mutação in-place? | A ordenação/transformação muta a entrada do chamador? |
| Erro no tipo de retorno? | O caminho de falha está no tipo ou escondido? |

**Registro** (`experiments/agent-eval/results/`): uma pasta por execução, com as
duas saídas e a tabela preenchida. A hipótese a testar — e a ser confirmada ou
refutada honestamente — é que a Condição B reduz efeitos colaterais, elimina a
necessidade de mock e faz o bug de mutação não aparecer.
