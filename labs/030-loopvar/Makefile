# Lab 030-loopvar — Makefile local.
#
# Este lab é um módulo Go isolado. O toolchain padrão é o 1.27 — e aqui isso
# não é detalhe: é o compilador NOVO que sabe compilar as duas semânticas, e
# só com ele os dois lados cabem num comando só. Um toolchain 1.21 não conhece
# `-lang=go1.22` e recusa a metade "depois".
#
# Para rodar no toolchain antigo de verdade (exige `go 1.21` no go.mod):
#
#	make run GOTOOLCHAIN=go1.21.13
#
# ATENÇÃO ao que decide o comportamento aqui: NÃO é o toolchain, é a VERSÃO DE
# LINGUAGEM. Trocar o GOTOOLCHAIN não muda o resultado — e um toolchain antigo
# nem compila um módulo que declare uma versão maior que ele. Quem manda é a
# linha `go` do go.mod, e ela pode ser trocada por chamada com -gcflags=-lang,
# para cima ou para baixo.
#
# Por isso `antes` e `depois` fixam a linguagem na flag em vez de depender do
# go.mod: mesmo fonte, mesmo compilador, só a versão de linguagem muda. Já o
# `run` segue o go.mod — é o alvo para experimentar trocando aquela linha.
#
# `make version` mostra qual toolchain está de fato rodando.

GO          ?= go
GOTOOLCHAIN ?= go1.27.0
export GOTOOLCHAIN

.PHONY: run test vet tidy all version antes depois comparar

## run: executa o lab
run:
	$(GO) run .

## antes: a semântica pré-1.22, fixada na flag (imprime 3 3 3)
antes:
	$(GO) run -gcflags=-lang=go1.21 .

## depois: a semântica 1.22+, fixada na flag (imprime 0 1 2)
depois:
	$(GO) run -gcflags=-lang=go1.22 .

## comparar: os dois lados na mesma tela — mesmo fonte, mesmo compilador
comparar:
	@printf 'linguagem 1.21 (antes):   '; $(GO) run -gcflags=-lang=go1.21 . | tr '\n' ' '; echo
	@printf 'linguagem 1.22 (depois):  '; $(GO) run -gcflags=-lang=go1.22 . | tr '\n' ' '; echo
	@printf 'linguagem do go.mod:      '; $(GO) run . | tr '\n' ' '; echo

## test: roda os testes deste lab
test:
	$(GO) test ./...

## vet: análise estática deste lab
vet:
	$(GO) vet ./...

## tidy: acerta o go.mod/go.sum deste lab
tidy:
	$(GO) mod tidy

## version: mostra o toolchain que este lab usa (deve dizer go1.27.x)
version:
	@$(GO) version
	@$(GO) list -m -f 'módulo {{.Path}} — go {{.GoVersion}}'

## all: verificação completa deste lab
all: vet test
