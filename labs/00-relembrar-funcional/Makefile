# Lab 00-relembrar-funcional — Makefile local.
#
# Este lab é um módulo Go isolado: a versão de Go que vale aqui está no go.mod
# desta pasta, não na da raiz. GOTOOLCHAIN=auto faz o Go baixar o toolchain
# certo sozinho se o instalado for outro — é o que permite um lab pinado numa
# versão antiga rodar na mesma máquina que a raiz em 1.27.

GO          ?= go
GOTOOLCHAIN ?= auto
export GOTOOLCHAIN

.PHONY: run test vet tidy all

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

## all: verificação completa deste lab
all: vet test
