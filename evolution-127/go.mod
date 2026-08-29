// Submódulo separado de propósito. A razão original era de versão: o 1.27
// estava em RC e a raiz precisava compilar na estável. Com o 1.27 lançado
// (ago/2026) a raiz subiu junto, mas a separação FICA — e por um motivo melhor,
// que é o conteúdo da palestra: este pacote redeclara Option[T], Some, None e
// HalfIfEven para mostrar o "depois" com métodos genéricos. Os mesmos nomes
// existem em ../option_pre127.go como o "antes". O par só coexiste em pacotes
// separados — a colisão É a demonstração.
module gofunctional/evolution-127

go 1.27
