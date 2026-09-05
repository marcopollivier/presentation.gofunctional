# Lab 02 — ponteiro e função pura

> **O ponto:** `Len` e `Changer` recebem exatamente a mesma coisa, um
> `*string`. Uma lê, a outra reescreve a variável de quem chamou. A assinatura
> `func(*string)` não distingue as duas — o ponteiro não é o vilão, o que a
> função faz com ele é.

**Go:** `1.21` — o Go de quando a palestra nasceu. Pinado no `Makefile`, não só
declarado no `go.mod`: ver a nota no fim.

**Bloco da palestra:** 2

## Como rodar

```bash
make run       # os três momentos
make version   # prova qual toolchain está rodando (deve dizer go1.21.13)
make test      # testes
make all       # vet + test
```

Ou, da raiz do repo:

```bash
make lab LAB=02-estado-global-ponteiros TARGET=run
```

## O que a saída mostra

```
Momento 1: Endereco [0x14000122010] Nome [Marco Ollivier]
Momento 2: Endereco [0x14000122010] Nome [Marco Ollivier] Tamanho [14]
Momento 3: Endereco [0x14000122010] Nome [Marco Paulo Ollivier] Tamanho [20]
```

O endereço é **o mesmo nas três linhas**. É a coisa mais importante da tela, e a
mais fácil de passar batido: passar um ponteiro não move a variável nem cria
outra — entrega a chave da caixa onde ela mora. A caixa nunca mudou; o que
mudou foi o conteúdo, entre o Momento 2 e o Momento 3.

E o conteúdo mudou porque alguém com a chave decidiu escrever. `Changer` tem a
mesma forma de assinatura que `Len` e faz o oposto.

## O incômodo que este lab planta

`Len` é pura no sentido que importa: não escreve, não guarda estado, não faz
I/O. A anotação no código está certa.

Mesmo assim, olhe o que o Momento 3 faz com ela:

```go
Len(&nome)   // 14
Changer(&nome)
Len(&nome)   // 20  — mesmo argumento
```

Mesma função pura, mesmo argumento, resposta diferente. O que se perdeu não foi
a pureza — foi a **transparência referencial** do [lab 01](../01-relembrar-funcional/):
`Len(&nome)` deixou de poder ser trocada pelo seu resultado, porque o resultado
agora depende de *quando* a chamada acontece. Com `Len(nome)`, por valor, o
argumento seria o texto, e 14 seria 14 para sempre.

É por isso que "pedir ponteiro" custa caro mesmo quando a função se comporta:
o argumento deixa de ser um valor e passa a ser um endereço que qualquer um
pode reescrever no meio do caminho. Fora o estado a mais que vem junto — `nil`,
que aqui vira panic; não existe `string` nil (há um teste só para isso).

E o de sempre: uma `string` já é passada como cabeçalho de 16 bytes, sem copiar
os caracteres. O argumento de performance não sustenta o `*`.

`&nome` no ponto de chamada é, pelo menos, um aviso. Guarde a imagem: com
slice, a mesma permissão é entregue **sem nenhum `&` ou `*` à vista** — é o par
`SortedPrint`/`SafeSortedPrint` em
[`exemplos/00-basico-sobre-ponteiros`](../../exemplos/00-basico-sobre-ponteiros/).

## Testes

`make test` fixa por escrito o que a tela mostra: que `Len` não altera nada, que
`Changer` reescreve o chamador, que o endereço **não** muda depois da escrita, e
o teste que carrega o incômodo — `TestMesmoArgumentoRespostaDiferente`, 14 antes
e 20 depois.

## Por que o toolchain está pinado

`GOTOOLCHAIN=auto`, que os outros labs usam, só **sobe** de versão. Com
`go 1.21` no `go.mod` e o 1.27 instalado na máquina, o lab compilaria no 1.27 em
modo de linguagem 1.21 — parecido, mas não é a mesma coisa. O `Makefile` fixa
`go1.21.13` para o lab rodar no 1.21 de verdade.
