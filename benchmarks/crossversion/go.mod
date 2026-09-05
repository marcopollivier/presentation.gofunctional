module gofunctional/benchmarks/crossversion

// A diretiva `go` está fixada em 1.21 DE PROPÓSITO, e não é a versão mínima do
// toolchain: é a versão da LINGUAGEM. Ela decide quais construções o compilador
// aceita e com que semântica — a variável de loop por iteração do 1.22, por
// exemplo, é gated por esta linha, não pelo compilador instalado.
//
// Fixando em 1.21, o mesmo fonte compila com semântica idêntica nos sete
// toolchains da matriz, e a única variável que sobra entre uma coluna e outra é
// compilador/runtime/GC. Se aqui dissesse `go 1.27`, o toolchain 1.21 recusaria
// compilar e a matriz morreria na primeira coluna.
//
// Não há diretiva `toolchain` de propósito: quem escolhe a versão é o run.sh,
// por rodada, via GOTOOLCHAIN.
go 1.21
