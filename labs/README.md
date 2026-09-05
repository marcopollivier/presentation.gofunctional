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
| 01 | [`01-relembrar-funcional`](01-relembrar-funcional/) | `f(x) = x² + 3` no papel vs em Go. A pura e a que suja um global devolvem o mesmo 7 e têm a **mesma assinatura** — o contrato não distingue as duas. | 1.27 | 1 |
| 02 | [`02-estado-global-ponteiros`](02-estado-global-ponteiros/) | `Len` e `Changer` recebem o mesmo `*string`: uma lê, a outra reescreve o chamador. O endereço não muda nos três momentos — e `Len(&nome)` responde 14 e depois 20. | **1.21** | 2 |

## Como rodar

Da raiz do repositório:

```bash
make labs                          # testes de todos os labs
make lab LAB=01-nome               # um lab só
make lab LAB=01-nome TARGET=run    # executa em vez de testar
make labs-vet                      # go vet em cada lab
```

Ou de dentro de um lab, que é o que dá menos atrito no palco:

```bash
cd labs/01-nome
make run
```

## Como criar um lab

```bash
make lab-new NAME=01-nome                  # usa a versão de Go do módulo raiz
make lab-new NAME=02-antes GOVERSION=1.21  # pina numa versão anterior
```

Isso copia `_template/`, substitui os placeholders e cria `go.mod`, `Makefile`,
`README.md` e `main.go` já com o comentário de topo esboçado. Ele **não**
sobrescreve pasta existente. Depois: preencha o comentário de topo e acrescente a
linha do lab na tabela acima.

`GOTOOLCHAIN=auto` está no Makefile de cada lab, então um lab pinado numa versão
que você não tem instalada baixa o toolchain sozinho na primeira execução.

## Convenções

- **`NN-slug`** — o número é a ordem de aparição no palco, não hierarquia. Como
  os labs são módulos isolados e **nenhum importa o outro**, renumerar é `git mv`
  + editar a linha `module`. Nada quebra.
- **O diretório tem hífen; o pacote, não** — Go não aceita hífen em nome de
  pacote. Mesma divergência intencional de `exemplos/`.
- **Comentário de topo obrigatório** no arquivo principal: qual bloco da
  palestra, o que o lab demonstra e — se ele for deliberadamente "errado" — por
  quê. Como em `exemplos/`, o lado esquerdo do slide é conteúdo: esse comentário
  é o que impede que alguém (ou algum agente) "conserte" a demonstração.
- **`_template/` não é um lab.** O `_` faz o toolchain do Go ignorar a pasta, e os
  arquivos `.tmpl` a mantêm fora do wildcard do Makefile.
