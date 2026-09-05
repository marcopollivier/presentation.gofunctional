# Lab 030 — a variável de loop

> **O ponto:** o footgun de doze anos que sumiu sem ninguém reescrever código.
> O fonte é o mesmo; mudou o significado do `for`. Linguagem 1.21 imprime
> `3 3 3`, linguagem 1.22 imprime `0 1 2`.

**Go:** toolchain `1.27` (é o compilador novo que sabe compilar **as duas**
semânticas). A linguagem é o que importa aqui, e ela é escolhida por alvo — ver
abaixo.

**Bloco da palestra:** 3

## Como rodar

```bash
make comparar   # os dois lados na mesma tela, mesmo compilador
make antes      # linguagem 1.21 -> 3 3 3
make depois     # linguagem 1.22 -> 0 1 2
make run        # segue a linha `go` do go.mod (hoje: 1.21)
make test       # os dois comportamentos, provados de uma vez
```

Saída de `make comparar`:

```
linguagem 1.21 (antes):   3 3 3
linguagem 1.22 (depois):  0 1 2
linguagem do go.mod:      3 3 3
```

## O que decide o resultado (a parte que engana)

**Não é o toolchain.** É a **versão de linguagem**. Trocar o `GOTOOLCHAIN` não
muda nada — e um toolchain antigo nem compila um módulo que declare versão
maior que ele:

```
$ GOTOOLCHAIN=go1.21.13 go run .     # com go.mod dizendo go 1.27
go: go.mod requires go >= 1.27 (running go 1.21.13)
```

A versão de linguagem sai da linha `go` do `go.mod`, e pode ser trocada por
chamada — para cima ou para baixo — com `-gcflags=-lang`:

```bash
go run -gcflags=-lang=go1.21 .   # 3 3 3
go run -gcflags=-lang=go1.22 .   # 0 1 2
```

É o que os alvos `antes` e `depois` fazem, e é por isso que eles **não**
dependem do `go.mod`: você pode mexer naquela linha à vontade que a comparação
continua funcionando.

O `go.mod` está em `go 1.21` de propósito — assim `make run` mostra o "antes".
Aquela linha é o botão do experimento. Para rodar no toolchain antigo de
verdade (e não só na linguagem antiga):

```bash
make run GOTOOLCHAIN=go1.21.13
```

## Os dois comportamentos no mesmo pacote

O artefato mais direto deste lab está nos testes. Um arquivo começa com
`//go:build go1.21` e outro com `//go:build go1.22` — a build tag fixa a versão
de linguagem **por arquivo**, então os dois comportamentos coexistem no mesmo
pacote, na mesma compilação:

```
$ go test -v ./...
--- PASS: TestLinguagem121CompartilhaAVariavel   (3 3 3)
--- PASS: TestLinguagem122TemUmIPorIteracao      (0 1 2)
```

Um `go test`, um binário, as duas semânticas provadas lado a lado. A diferença
entre os dois arquivos é **uma linha** de build tag.

Rodando num toolchain anterior ao 1.22, o arquivo do "depois" é simplesmente
excluído do build e só o "antes" roda — nada quebra:

```bash
make test GOTOOLCHAIN=go1.21.13   # passa, com um teste a menos
```

## Por que isso é o primeiro item do bloco 3

Não houve migração. Nenhum `i := i` precisou ser escrito, nenhum código
existente quebrou — o compilador passou a fazer o certo por padrão, e o
`//go:build` de quem ficou para trás mantém o comportamento antigo. É o modelo
de evolução que a palestra usa como fio condutor: a linguagem andou na direção
do paradigma sem escrever "funcional" em release note nenhuma.

O contraexemplo curado, com o par `LoopVarShared`/`LoopVarPerIteration`, está em
[`exemplos/02-loopvar`](../../exemplos/02-loopvar/). Aqui o mesmo ponto **roda**,
e os dois lados aparecem na mesma tela.
