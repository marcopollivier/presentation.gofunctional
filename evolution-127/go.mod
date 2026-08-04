// Submódulo separado de propósito: os exemplos aqui exigem Go 1.27 (métodos
// genéricos), que na época desta escrita ainda está em RC. Manter isto fora do
// módulo raiz deixa a raiz compilando na estável e transforma o 1.27 num
// "estudo da evolução" opt-in. Rode com: GOTOOLCHAIN=auto go test ./...
module gofunctional/evolution-127

go 1.27

toolchain go1.27rc1
