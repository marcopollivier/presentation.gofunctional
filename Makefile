# Makefile — material de apoio da palestra.
# Os alvos de benchmark geram a saída que vira slide.

BENCHTIME ?= 300ms
COUNT     ?= 6

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
	go test -bench=Sort -benchmem -run='^$$' -benchtime=$(BENCHTIME) ./...

## benchstat: roda o bench de sort COUNT vezes e resume com benchstat.
## Requer: go install golang.org/x/perf/cmd/benchstat@latest
benchstat:
	go test -bench=Sort -benchmem -run='^$$' -benchtime=$(BENCHTIME) -count=$(COUNT) ./... | tee bench.txt
	benchstat bench.txt

## evolution: roda o submódulo Go 1.27 (baixa o toolchain via GOTOOLCHAIN=auto)
evolution:
	cd evolution-127 && GOTOOLCHAIN=auto go test ./...

## all: verificação completa do módulo raiz
all: vet test bench-sort
