# Makefile — dispatcher dos labs.
#
# Cada lab é um módulo Go isolado (go.mod próprio). Este arquivo não conhece
# nenhum deles por nome: a lista sai do wildcard abaixo, então criar a pasta com
# um go.mod dentro já basta para o lab entrar em `make test` / `make vet`.
#
# `_template/` fica de fora de graça — lá o arquivo se chama `go.mod.tmpl`.

LABS     := $(sort $(patsubst %/go.mod,%,$(wildcard */go.mod)))
TEMPLATE := _template

# Versão de Go dos labs novos: por padrão, a mesma do módulo raiz. Passe
# GOVERSION=1.21 para pinar um lab numa versão anterior — é o caso de uso que
# justifica o go.mod isolado.
GOVERSION ?= $(shell awk '/^go /{print $$2; exit}' ../go.mod)

# Alvo que `make lab LAB=...` repassa ao Makefile do lab.
TARGET ?= test

.PHONY: test vet tidy all list lab lab-new

## test, vet, tidy, all: repassa o alvo a todos os labs
test vet tidy all:
	@if [ -z "$(LABS)" ]; then \
		echo "nenhum lab ainda — crie um com: make lab-new NAME=01-nome"; \
		exit 0; \
	fi
	@for lab in $(LABS); do \
		echo "==> $$lab ($@)"; \
		$(MAKE) --no-print-directory -C $$lab $@ || exit 1; \
	done

## list: lista os labs existentes
list:
	@if [ -z "$(LABS)" ]; then echo "(nenhum lab ainda)"; else printf '%s\n' $(LABS); fi

## lab: roda um lab só — make lab LAB=01-nome [TARGET=run]
lab:
	@test -n "$(LAB)" || { echo "uso: make lab LAB=<pasta> [TARGET=run]" >&2; exit 2; }
	@test -f "$(LAB)/go.mod" || { echo "lab '$(LAB)' não existe (veja: make list)" >&2; exit 2; }
	@$(MAKE) --no-print-directory -C $(LAB) $(TARGET)

## lab-new: cria um lab a partir do template — make lab-new NAME=01-nome [GOVERSION=1.21]
lab-new:
	@test -n "$(NAME)" || { echo "uso: make lab-new NAME=NN-slug [GOVERSION=1.21]" >&2; exit 2; }
	@test ! -e "$(NAME)" || { echo "'$(NAME)' já existe — não vou sobrescrever" >&2; exit 2; }
	@mkdir -p "$(NAME)"
	@for f in $(TEMPLATE)/*.tmpl; do \
		out="$(NAME)/$$(basename $$f .tmpl)"; \
		sed -e 's|__NAME__|$(NAME)|g' -e 's|__GOVERSION__|$(GOVERSION)|g' "$$f" > "$$out"; \
	done
	@echo "criado: labs/$(NAME) (go $(GOVERSION))"
	@echo "próximo: make lab LAB=$(NAME) TARGET=run  — e anote o lab no labs/README.md"
