# Lab 02 — ponteiro e função pura

> **O ponto:** `Len(&nome)` pede um `*string` para fazer uma coisa que não
> precisa de ponteiro nenhum — ler o tamanho. Pedir ponteiro é pedir permissão
> de escrita. Esta função pediu e não usa.

**Go:** `1.21` — o Go de quando a palestra nasceu. Pinado no `Makefile`, não só
declarado no `go.mod`: ver a nota no fim.

**Bloco da palestra:** 2

## Como rodar

```bash
make run       # imprime 14
make version   # prova qual toolchain está rodando (deve dizer go1.21.13)
make vet       # análise estática
```

Ou, da raiz do repo:

```bash
make lab LAB=02-estado-global-ponteiros TARGET=run
```

## Notas

O programa imprime `14` — o mesmo `14` que sairia de `len(nome)` direto. O
ponteiro não comprou nada: uma `string` em Go já é passada como cabeçalho
(ponteiro para os bytes + tamanho), 16 bytes, e os caracteres não são copiados.
O argumento de performance não se sustenta.

O que o `*string` comprou foi custo. A função ganhou permissão de escrever na
variável do chamador — permissão que ela não usa, mas que quem lê a assinatura
não tem como descartar. E ganhou um estado que a versão por valor não tem:
`s == nil`, que aqui vira panic. Não existe `string` nil.

`&nome` no ponto de chamada é o aviso. Guarde essa imagem: no exemplo seguinte,
com slice, a mesma permissão é entregue **sem nenhum `&` ou `*` à vista**.

### Por que o toolchain está pinado

`GOTOOLCHAIN=auto`, que os outros labs usam, só **sobe** de versão. Com
`go 1.21` no `go.mod` e o 1.27 instalado na máquina, o lab compilaria no 1.27 em
modo de linguagem 1.21 — parecido, mas não é a mesma coisa. O `Makefile` fixa
`go1.21.13` para o lab rodar no 1.21 de verdade.
