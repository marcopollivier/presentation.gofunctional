# Lab 022-mutacao-nao-observavel-moderna — Makefile local.
#
# Este lab é um módulo Go isolado e roda no Go 1.27 — é o "depois" do par com o
# 021, que está pinado no 1.21.
#
# Aqui o pin não é necessário para a versão existir (o 1.27 é a corrente): é
# para o palco ser determinístico e para os dois labs do par exibirem, lado a
# lado, toolchains diferentes de verdade. `auto` também funcionaria.
#
# `make version` mostra qual toolchain está de fato rodando.

GO          ?= go
GOTOOLCHAIN ?= go1.27.0
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

## version: mostra o toolchain que este lab usa (deve dizer go1.27.x)
version:
	@$(GO) version
	@$(GO) list -m -f 'módulo {{.Path}} — go {{.GoVersion}}'

## all: verificação completa deste lab
all: vet test
