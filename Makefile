# Makefile — material de apoio da palestra.
# Os alvos de benchmark geram a saída que vira slide.

BENCHTIME ?= 300ms
COUNT     ?= 6

# Pacote do benchmark de ordenação. Isolado numa variável porque os alvos de
# palco apontam para ele: rodar `./...` varreria os oito exemplos e encheria a
# tela de "ok ... 0.1s" antes da tabela que interessa.
SORTPKG   ?= ./exemplos/07-benchmark-sort/

# Pasta dos exemplos reais. Cada lab lá dentro é um módulo isolado; este
# Makefile só delega para o dispatcher de labs/, que descobre os labs sozinho.
LABS_DIR  ?= labs

.PHONY: test vet bench bench-sort benchstat evolution all \
        labs labs-vet lab lab-new

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

## labs: roda os testes de todos os labs (exemplos reais, módulos isolados)
labs:
	@$(MAKE) --no-print-directory -C $(LABS_DIR) test

## labs-vet: análise estática em cada lab
labs-vet:
	@$(MAKE) --no-print-directory -C $(LABS_DIR) vet

## lab: roda um lab só — make lab LAB=020-nome [TARGET=run]
lab:
	@$(MAKE) --no-print-directory -C $(LABS_DIR) lab LAB=$(LAB) $(if $(TARGET),TARGET=$(TARGET))

## lab-new: cria um lab a partir do template —
## make lab-new NAME=021-nome [GOVERSION=1.21]
lab-new:
	@$(MAKE) --no-print-directory -C $(LABS_DIR) lab-new NAME=$(NAME) $(if $(GOVERSION),GOVERSION=$(GOVERSION))

## all: verificação completa do módulo raiz.
## Não inclui `labs` de propósito: um lab pinado noutra versão pode ter que
## baixar toolchain, e isso não pode travar a verificação de palco. Para tudo:
## make all labs
all: vet test bench-sort
