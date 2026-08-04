# evolution-127 — o "depois" do divisor de águas (Go 1.27)

Submódulo **separado de propósito**. Os exemplos aqui exigem **Go 1.27**
(métodos genéricos), que na escrita deste material ainda está em RC. Mantê-los
fora do módulo raiz deixa a raiz compilando na estável (Go 1.26) e transforma o
1.27 num estudo da evolução **opt-in**.

## Como rodar

```bash
# baixa o toolchain 1.27 automaticamente e roda os testes
GOTOOLCHAIN=auto go test ./...

# ou, da raiz do repo:
make evolution
```

## O que tem aqui

| Arquivo | Ponto |
|--|--|
| `option.go` | O mesmo `Option[T]` de `../option_pre127.go`, agora com **métodos genéricos** — `Map[B]`/`Chain[B]` declaram seus próprios parâmetros de tipo. O uso vira encadeável: `Some(n).Chain(f).Map(g).GetOrElse(d)`. |
| `option_test.go` | Prova que o encadeamento funciona, inclusive quando `Map` **muda o tipo** (`int → string`). |
| `ceiling_hkt.go` | **Não compila de propósito** (`//go:build ignore`). Demonstra o teto: interface não pode exigir um `Map` genérico → Go tem monads concretos, mas **não tem higher-kinded types**. |

## A conclusão

Go ganhou monads concretos e encadeáveis. O que ele **não** dá é abstração
sobre "qualquer Functor". A forçação de barra saiu da **sintaxe** (resolvida
pelos métodos genéricos) e foi para a **capacidade de abstração**.

Veja o erro do teto com:

```bash
GOTOOLCHAIN=auto go build ceiling_hkt.go
# ./ceiling_hkt.go:34:5: interface method must have no type parameters
```
