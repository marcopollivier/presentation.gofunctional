# evolution-127 — o "depois" do divisor de águas (Go 1.27)

Submódulo **separado de propósito**. Os exemplos aqui exigem **Go 1.27**
(métodos genéricos). A razão original da separação era de versão — o 1.27 estava
em RC e a raiz tinha que compilar na estável. Com o **1.27 lançado (ago/2026)** a
raiz subiu junto, e a separação passou a valer por um motivo melhor: este pacote
**redeclara** `Option[T]`, `Some`, `None` e `HalfIfEven` para mostrar o "depois",
e os mesmos nomes vivem em `../exemplos/06-option/option_pre127.go` como o "antes". O par só
coexiste em pacotes separados — **a colisão é a demonstração**.

## Como rodar

```bash
# 1.27 é estável; o toolchain resolve sozinho
go test ./...

# ou, da raiz do repo:
make evolution
```

## O que tem aqui

| Arquivo | Ponto |
|--|--|
| `option.go` | O mesmo `Option[T]` de `../exemplos/06-option/option_pre127.go`, agora com **métodos genéricos** — `Map[B]`/`Chain[B]` declaram seus próprios parâmetros de tipo. O uso vira encadeável: `Some(n).Chain(f).Map(g).GetOrElse(d)`. |
| `option_test.go` | Prova que o encadeamento funciona, inclusive quando `Map` **muda o tipo** (`int → string`). |
| `ceiling_hkt.go` | **Não compila de propósito** (`//go:build ignore`). Demonstra o teto: interface não pode exigir um `Map` genérico → Go tem monads concretos, mas **não tem higher-kinded types**. |

## A conclusão

Go ganhou monads concretos e encadeáveis. O que ele **não** dá é abstração
sobre "qualquer Functor". A forçação de barra saiu da **sintaxe** (resolvida
pelos métodos genéricos) e foi para a **capacidade de abstração**.

Veja o erro do teto com:

```bash
go build ceiling_hkt.go
# ./ceiling_hkt.go:35:5: interface method must have no type parameters
```
