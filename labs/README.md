# Labs — exemplos reais

`exemplos/` é o material curado: trechos pequenos, sem `main`, que cabem num
slide. **Aqui é o outro lado.** Um lab é um exemplo *real*: um programa que roda,
que pode ter dependência externa e que pode estar pinado numa versão de Go
diferente da raiz — porque é isso que ele demonstra.

Cada lab é um **módulo Go isolado**: `go.mod` próprio, versão própria,
dependências próprias. Consequência prática que importa: nenhum lab é alcançado
pelo `go test ./...` nem pelo `make all` da raiz. Um lab pinado no Go 1.21, ou com
uma dep que quebrou, **não derruba** o material da palestra.

## O caminho

| # | Pasta | O ponto | Go | Bloco |
|--|--|--|--|--|
| 010 | [`010-relembrar-funcional`](010-relembrar-funcional/) | `f(x) = x² + 3` no papel vs em Go. A pura e a que suja um global devolvem o mesmo 7 e têm a **mesma assinatura** — o contrato não distingue as duas. | 1.27 | 1 |
| 020 | [`020-estado-global-ponteiros`](020-estado-global-ponteiros/) | `Len` e `Changer` recebem o mesmo `*string`: uma lê, a outra reescreve o chamador. O endereço não muda nos três momentos — e `Len(&nome)` responde 14 e depois 20. | **1.21** | 2 |

## Como rodar

Da raiz do repositório:

```bash
make labs                          # testes de todos os labs
make lab LAB=020-nome              # um lab só
make lab LAB=020-nome TARGET=run   # executa em vez de testar
make labs-vet                      # go vet em cada lab
```

Ou de dentro de um lab, que é o que dá menos atrito no palco:

```bash
cd labs/01-nome
make run
```

## Como criar um lab

```bash
make lab-new NAME=021-nome                 # usa a versão de Go do módulo raiz
make lab-new NAME=022-antes GOVERSION=1.21 # pina numa versão anterior
```

Isso copia `_template/`, substitui os placeholders e cria `go.mod`, `Makefile`,
`README.md` e `main.go` já com o comentário de topo esboçado. Ele **não**
sobrescreve pasta existente. Depois: preencha o comentário de topo e acrescente a
linha do lab na tabela acima.

O template traz `GOTOOLCHAIN=auto` no Makefile do lab, então uma versão mais
nova que a instalada é baixada sozinha na primeira execução. Para uma versão
mais **antiga**, o `auto` não serve — ele só sobe. Nesse caso troque o `auto`
pela versão exata (`GOTOOLCHAIN ?= go1.21.13`), como faz o
[`020-estado-global-ponteiros`](020-estado-global-ponteiros/).

## Convenções

- **`CEE-slug`** — o número diz **capítulo** e **exemplo dentro do capítulo**:
  `020` é o primeiro exemplo do capítulo 2, `021` o segundo, `022` o terceiro.
  Um capítulo cabe em quantos labs precisar sem renumerar o resto, que é o
  motivo de existir o dígito extra. A ordem é a do palco, não uma hierarquia.
  E como os labs são módulos isolados e **nenhum importa o outro**, renumerar
  quando for preciso é `git mv` + editar a linha `module`. Nada quebra.
- **O diretório tem hífen; o pacote, não** — Go não aceita hífen em nome de
  pacote. Mesma divergência intencional de `exemplos/`.
- **Comentário de topo obrigatório** no arquivo principal: qual bloco da
  palestra, o que o lab demonstra e — se ele for deliberadamente "errado" — por
  quê. Como em `exemplos/`, o lado esquerdo do slide é conteúdo: esse comentário
  é o que impede que alguém (ou algum agente) "conserte" a demonstração.
- **`_template/` não é um lab.** O `_` faz o toolchain do Go ignorar a pasta, e os
  arquivos `.tmpl` a mantêm fora do wildcard do Makefile.
