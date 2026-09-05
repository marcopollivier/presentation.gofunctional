# Lab 022 — a mutação não observável, versão moderna

> **O ponto:** o "depois" do [021](../021-mutacao-nao-observavel/). Guardar a
> promessa de não tocar no slice do chamador custava três linhas e uma variável
> temporária; hoje custa uma, e a stdlib faz o clone. O argumento ergonômico
> que sustentava a mutação em 2023 caiu.

**Go:** `1.27` — a corrente. O par com o 021, que roda no 1.21, existe para os
dois estarem na tela ao mesmo tempo.

**Bloco da palestra:** 2, com um pé no 3

## Como rodar

```bash
make run       # o antes/depois
make version   # prova qual toolchain está rodando (deve dizer go1.27.x)
make test      # testes
make all       # vet + test
```

Ou, da raiz do repo:

```bash
make lab LAB=022-mutacao-nao-observavel-moderna TARGET=run
```

## O que a saída mostra

```
x antes da chamada   [9 1 3]
x depois da chamada  [1 3 9]
x depois da chamada  [9 1 3]
```

A **segunda** linha sai de dentro de `SortedPrint` e mostra o resultado
ordenado. A **terceira** é o `main` mostrando o `x` do chamador — inalterado.
Compare com o 021, onde a linha do chamador vinha reescrita.

## As três linhas que viraram uma

```go
// 021 — Go 1.21
c := make([]int, len(x))
copy(c, x)
sort.Ints(c)

// 022 — Go 1.23+
c := slices.Sorted(slices.Values(x))
```

`slices.Values(x)` transforma o slice num iterador (`iter.Seq`) e
`slices.Sorted` consome esse iterador para um slice **novo**, já ordenado. O
original nunca é escrito — não porque alguém tomou cuidado, mas porque nesta
formulação não existe caminho para escrever nele.

É por aqui que o capítulo desmonta o argumento de 2023. Naquela época a defesa
da mutação era ergonômica: `sort.Ints(x)` era uma linha, o caminho imutável era
um parágrafo. O 1.21 trouxe `slices` para a stdlib, o 1.23 trouxe os
iteradores, e a versão que **não** muta virou a mais curta das duas. Sem
release note nenhuma escrever a palavra "funcional".

A mutação, aliás, continua existindo: `slices.Sorted` ordena o slice que ela
mesma alocou. Continua não sendo observável — que é o assunto do 021.

## O mesmo nome, o comportamento oposto

`SortedPrint` aqui tem o **mesmo nome** e a **mesma assinatura** —
`func([]int)` — da `SortedPrint` do 021, e faz o contrário: lá reescrevia o
slice do chamador, aqui não encosta nele.

Nem o nome nem a assinatura contam a verdade. Só o corpo. É a tese do capítulo
levada ao limite, e vale deixar o par aberto lado a lado no palco.

## O custo, que o bloco 4 cobra

`slices.Sorted` aloca um slice novo a cada chamada. Essa alocação é exatamente
o que o benchmark de [`exemplos/07-benchmark-sort`](../../exemplos/07-benchmark-sort/)
mede — e em 2023 ela foi lida como "o preço do paradigma". O post-mortem lá
mostra que a conta estava errada, e é o gancho para o bloco 4.

## Testes

`make test` prova que o argumento sai como entrou, para várias entradas; que o
resultado é o mesmo do caminho antigo de três linhas (mudou a ergonomia, não a
semântica); e que escrever no resultado **não** alcança a entrada — o oposto
exato do `TestAMutacaoAlcancaOArrayDeOrigem` do 021.

O `Example` fixa a saída de palco. As duas últimas linhas dela usam o mesmo
rótulo, `x depois da chamada`: a do meio é o resultado, a de baixo é o
chamador. Se o texto do `Printf` mudar, o `Example` falha e é só atualizar.
