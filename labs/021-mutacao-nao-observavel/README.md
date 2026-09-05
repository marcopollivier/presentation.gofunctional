# Lab 021 — a mutação que não se observa

> **O ponto:** `SortedPrint` e `SafeSortedPrint` têm a **mesma assinatura**,
> `func([]int)`, e não há um `&` nem um `*` em lugar nenhum — nem na
> declaração, nem na chamada. Mesmo assim uma reescreve o slice de quem chamou
> e a outra não. E o detalhe que dá nome ao lab: **as duas mutam**. A diferença
> é a mutação escapar ou não.

**Go:** `1.21` — o Go de quando a palestra nasceu. Pinado no `Makefile`, não só
declarado no `go.mod`: ver a nota no fim.

**Bloco da palestra:** 2

## Como rodar

```bash
make run       # o antes/depois das duas versões
make version   # prova qual toolchain está rodando (deve dizer go1.21.13)
make test      # testes
make all       # vet + test
```

Ou, da raiz do repo:

```bash
make lab LAB=021-mutacao-nao-observavel TARGET=run
```

## O que a saída mostra

```
Imutavel: x antes da chamada  [9 1 3]
Imutavel: x depois da chamada [9 1 3]
Mutavel: x antes da chamada   [9 1 3]
Mutavel: x depois da chamada  [1 3 9]
```

Mesma variável `x`, mesma forma de chamada, resultados opostos na linha
"depois". Nada no ponto de chamada distingue os dois casos — é preciso abrir o
corpo da função.

## Por que acontece

Um slice é um **cabeçalho**: ponteiro para o array, `len` e `cap`. Ele é
passado **por valor** — a função recebe uma cópia do cabeçalho, barata e
inofensiva. Só que essa cópia aponta para o **mesmo array**.

Quem recebe um `[]int` recebe, de graça, permissão de escrita na memória do
chamador. E a assinatura não tem um único caractere que faça você desconfiar.

É o contraste com o [lab 020](../020-estado-global-ponteiros/): lá o `&nome` no
ponto de chamada pelo menos **avisava**. Aqui não há aviso nenhum, e é por isso
que este é o caso perigoso dos dois.

## A mutação não observável

`SafeSortedPrint` também muta — `sort.Ints(c)` reescreve `c` inteiro. A
diferença é que `c` foi alocado dentro dela, não sai de lá e morre no `return`.
Ninguém no universo consegue notar que aquela mutação existiu.

É a definição operacional que a palestra usa: **pureza não é abstinência de
mutação, é a mutação não escapar**. A regra 2 do `CLAUDE.md` diz a mesma coisa
por outro lado — "não mute o slice do chamador", não "não mute".

O `make` + `copy` é a forma explícita, que mostra o que está acontecendo. No
1.21 já dava para escrever `slices.Clone(x)` numa linha; do 1.23 em diante o
trabalho inteiro cabe em `slices.Sorted(slices.Values(x))`. Essa evolução é
assunto do bloco 3.

## Testes

`make test` fixa os dois lados — que `SortedPrint` muta o chamador e que
`SafeSortedPrint` não — e mais dois que a tela não mostra:

- **`TestAMutacaoAlcancaOArrayDeOrigem`** — passar `original[:3]` entrega o
  array de trás inteiro. A escrita atravessa a fatia; o chamador não ofereceu
  aquilo, e perdeu mesmo assim.
- **`TestOContrasteSoApareceComEntradaDesordenada`** — com entrada já ordenada
  as duas funções deixam o mesmo resultado. É por isso que, no `main`, a demo
  imutável vem **primeiro**: depois de `SortedPrint` o `x` já está ordenado, e
  invertendo a ordem as quatro linhas sairiam iguais. O teste existe para
  ninguém "arrumar" o `main` sem perceber.

## Por que o toolchain está pinado

`GOTOOLCHAIN=auto`, que o template usa, só **sobe** de versão. Com `go 1.21` no
`go.mod` e o 1.27 instalado na máquina, o lab compilaria no 1.27 em modo de
linguagem 1.21 — parecido, mas não é a mesma coisa. O `Makefile` fixa
`go1.21.13` para o lab rodar no 1.21 de verdade.
