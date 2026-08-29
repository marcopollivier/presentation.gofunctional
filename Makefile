# Makefile — material de apoio da palestra.
# Os alvos de benchmark geram a saída que vira slide.

BENCHTIME ?= 300ms
COUNT     ?= 6

# Pacote do benchmark de ordenação. Isolado numa variável porque os alvos de
# palco apontam para ele: rodar `./...` varreria os oito exemplos e encheria a
# tela de "ok ... 0.1s" antes da tabela que interessa.
SORTPKG   ?= ./exemplos/07-benchmark-sort/

.PHONY: test vet bench bench-sort benchstat evolution all

## test: roda todos os testes do módulo raiz
test:
	go test ./...

## vet: análise estática
vet:
	go vet ./...

## bench: roda todos os benchmarks com estatística de alocação
bench:
	go test -bench=. -benchmem -run='^$$' -benchtime=$(BENCHTIME) ./...

## bench-sort: só o benchmark de ordenação (o achado central da palestra)
bench-sort:
	go test -bench=Sort -benchmem -run='^$$' -benchtime=$(BENCHTIME) $(SORTPKG)

## benchstat: roda o bench de sort COUNT vezes e resume com benchstat.
## Requer: go install golang.org/x/perf/cmd/benchstat@latest
benchstat:
	go test -bench=Sort -benchmem -run='^$$' -benchtime=$(BENCHTIME) -count=$(COUNT) $(SORTPKG) | tee bench.txt
	benchstat bench.txt

## evolution: roda o submódulo Go 1.27 (GOTOOLCHAIN=auto cobre quem ainda
## estiver num toolchain anterior ao 1.27)
evolution:
	cd evolution-127 && GOTOOLCHAIN=auto go test ./...

## all: verificação completa do módulo raiz
all: vet test bench-sort
