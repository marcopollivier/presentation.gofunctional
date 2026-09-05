# Lab 02-estado-global-ponteiros — Makefile local.
#
# Este lab é um módulo Go isolado e roda no Go 1.21.
#
# O GOTOOLCHAIN aqui é PINADO, não `auto` como nos outros labs. O `auto` só
# SOBE de versão: com `go 1.21` no go.mod e o 1.27 instalado, o lab compilaria
# no 1.27 em modo de linguagem 1.21 — parecido, mas não é a mesma coisa, e a
# promessa deste lab é rodar no 1.21 de verdade. A versão exata abaixo faz o Go
# baixar o toolchain certo na primeira vez e usar o cache depois.
#
# `make version` mostra qual toolchain está de fato rodando.

GO          ?= go
GOTOOLCHAIN ?= go1.21.13
export GOTOOLCHAIN

.PHONY: run test vet tidy all version

## run: executa o lab
run:
	$(GO) run .

## test: roda os testes deste lab
test:
	$(GO) test ./...

## vet: análise estática deste lab
vet:
	$(GO) vet ./...

## tidy: acerta o go.mod/go.sum deste lab
tidy:
	$(GO) mod tidy

## version: mostra o toolchain que este lab usa (deve dizer go1.21.x)
version:
	@$(GO) version
	@$(GO) list -m -f 'módulo {{.Path}} — go {{.GoVersion}}'

## all: verificação completa deste lab
all: vet test
