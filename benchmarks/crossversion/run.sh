#!/usr/bin/env bash
#
# Matriz de medição: {benchmark de 2023, medição corrigida} × {Go 1.21 … 1.27}.
#
# Rodar só o benchmark corrigido em várias versões mediria a soma de duas
# coisas — o que o Go melhorou e o viés do benchmark antigo — e o número não
# significaria nada. Rodando os dois em todas as versões, cada linha isola uma
# coisa só, e a diferença entre elas é ela mesma um resultado.
#
# Uso:
#   ./run.sh install    # pré-aquece os toolchains e o benchstat
#   ./run.sh run        # roda a matriz, um arquivo por versão
#   ./run.sh compare    # tabelas do benchstat
#   ./run.sh all        # os três acima, em ordem
#
# Variáveis: COUNT (default 10), BENCHTIME (default 1s), BENCH (default .)
#
#   COUNT=1 BENCHTIME=100ms ./run.sh run     # smoke, ~2 min
#   BENCH='/5$|/1k$' ./run.sh run            # só os tamanhos pequenos

set -euo pipefail

cd "$(dirname "$0")"

VERSIONS=(go1.21.13 go1.22.12 go1.23.12 go1.24.13 go1.25.14 go1.26.8 go1.27.1)

COUNT="${COUNT:-10}"
BENCHTIME="${BENCHTIME:-1s}"
BENCH="${BENCH:-.}"
RESULTS="${RESULTS:-results}"

# -cpu=1 remove a variação de escalonamento entre execuções.
# -run='^$' impede que os testes rodem junto e sujem a medição.
# -count é o que dá ao benchstat repetições para estimar variância; sem
# repetição não existe significância, e a comparação vira leitura de folha de
# chá.
FLAGS=(-run='^$' -bench="$BENCH" -benchmem -cpu=1 -count="$COUNT" -benchtime="$BENCHTIME")

# GOTOOLCHAIN é fixado na versão EXATA a cada rodada. É o ponto crítico do
# script: sem isso o Go pode resolver outro toolchain sozinho e eu mediria a
# mesma versão sete vezes sem perceber, com sete rótulos diferentes. O go.mod
# declara `go 1.21` justamente para que todo toolchain da lista aceite este
# módulo — então nada aqui força troca, e o pin explícito é a única coisa que
# escolhe a versão.
run_go() {
	local version="$1"; shift
	GOTOOLCHAIN="$version" go "$@"
}

cmd_install() {
	echo "==> pré-aquecendo toolchains (o primeiro uso de cada um baixa ~100 MB)"
	for v in "${VERSIONS[@]}"; do
		printf '  %-12s ' "$v"
		run_go "$v" version
	done

	echo "==> benchstat"
	if command -v benchstat >/dev/null 2>&1; then
		echo "  já instalado: $(command -v benchstat)"
	else
		go install golang.org/x/perf/cmd/benchstat@latest
	fi
}

# environment.txt é registro da medição, não enfeite: sem ele, um resultado
# guardado por seis meses não é auditável. Vai versionado junto.
write_environment() {
	local out="$RESULTS/environment.txt"
	{
		echo "# Registro da medição"
		echo "data:        $(date -u '+%Y-%m-%dT%H:%M:%SZ') (UTC)"
		echo "uname:       $(uname -a)"
		if [ "$(uname -s)" = "Darwin" ]; then
			echo "cpu:         $(sysctl -n machdep.cpu.brand_string)"
			echo "cores:       $(sysctl -n hw.ncpu) ($(sysctl -n hw.perflevel0.logicalcpu 2>/dev/null || echo '?') performance)"
			echo "mem:         $(( $(sysctl -n hw.memsize) / 1024 / 1024 / 1024 )) GB"
		fi
		echo "go (host):   $(go version)"
		echo
		echo "parâmetros:  count=$COUNT benchtime=$BENCHTIME bench=$BENCH cpu=1"
		echo "go.mod lang: $(awk '/^go /{print $2}' go.mod)"
		echo
		echo "toolchains:"
		for v in "${VERSIONS[@]}"; do
			printf '  %-12s %s\n' "$v" "$(run_go "$v" version 2>&1)"
		done
	} > "$out"
	echo "==> $out"
	cat "$out"
}

cmd_run() {
	mkdir -p "$RESULTS"
	write_environment

	for v in "${VERSIONS[@]}"; do
		# O arquivo leva o nome da versão porque o benchstat usa o nome do
		# arquivo como rótulo da coluna. Nome errado aqui = tabela mentirosa.
		local out="$RESULTS/$v.txt"
		echo "==> $v -> $out"
		run_go "$v" test "${FLAGS[@]}" . | tee "$out"
	done
}

cmd_compare() {
	local latest="${VERSIONS[${#VERSIONS[@]}-1]}"
	local files=()
	for v in "${VERSIONS[@]}"; do
		[ -f "$RESULTS/$v.txt" ] && files+=("$RESULTS/$v.txt")
	done

	if [ ${#files[@]} -eq 0 ]; then
		echo "sem resultados em $RESULTS/ — rode ./run.sh run antes" >&2
		exit 1
	fi

	echo "==> evolução entre versões (colunas = toolchains)"
	benchstat "${files[@]}" | tee "$RESULTS/compare-versions.txt"

	# O recorte 2023-vs-correto precisa de um truque: o benchstat compara
	# ARQUIVOS, não linhas dentro de um arquivo. Então separo os dois grupos em
	# dois arquivos e removo o prefixo "Legacy" dos nomes, para que
	# BenchmarkLegacyInPlace/1k e BenchmarkInPlace/1k virem a MESMA linha em
	# colunas diferentes — que é o pareamento que eu quero ver.
	echo
	echo "==> 2023 vs. corrigido, dentro do $latest"
	local src="$RESULTS/$latest.txt"
	if [ ! -f "$src" ]; then
		echo "  (sem $src)" >&2
		return 0
	fi

	grep -E '^(Benchmark(Legacy)?(InPlace|Immutable))' "$src" \
		| grep 'Legacy' | sed 's/BenchmarkLegacy/Benchmark/' > "$RESULTS/_legacy.txt"
	grep -E '^Benchmark(InPlace|Immutable)' "$src" \
		| grep -v 'InPlaceCopy' > "$RESULTS/_fixed.txt"

	# O cabeçalho (goos/goarch/pkg/cpu) some no grep; o benchstat aceita, mas
	# recolocá-lo mantém os arquivos legíveis fora dele.
	for f in "$RESULTS/_legacy.txt" "$RESULTS/_fixed.txt"; do
		{ grep -E '^(goos|goarch|pkg|cpu):' "$src"; cat "$f"; } > "$f.tmp" && mv "$f.tmp" "$f"
	done
	mv "$RESULTS/_legacy.txt" "$RESULTS/legacy.txt"
	mv "$RESULTS/_fixed.txt" "$RESULTS/fixed.txt"

	benchstat "$RESULTS/legacy.txt" "$RESULTS/fixed.txt" | tee "$RESULTS/compare-bias.txt"
}

case "${1:-all}" in
	install) cmd_install ;;
	run)     cmd_run ;;
	compare) cmd_compare ;;
	all)     cmd_install; cmd_run; cmd_compare ;;
	*)       echo "uso: $0 {install|run|compare|all}" >&2; exit 2 ;;
esac
