# Lab 01 — relembrar funcional

> **O ponto:** `f(x) = x² + 3` no papel é uma promessa: `f(2)` **é** 7, sempre,
> e dá para riscar `f(2)` e escrever `7` no lugar sem mudar a conta. `f()` em Go
> mantém a promessa; `badF()` devolve o mesmo 7 e mesmo assim a quebra, porque
> escreve num global. **As duas têm a mesma assinatura** — o "Contrato" do
> diagrama não distingue uma da outra.

**Go:** `1.27` — a corrente. Nada aqui depende da versão; é o capítulo de
abertura, feito só de stdlib.

**Bloco da palestra:** 1

## Como rodar

```bash
make run     # a demo de palco
make test    # testes
make all     # vet + test
```

Ou, da raiz do repo:

```bash
make lab LAB=01-relembrar-funcional TARGET=run
```

## Notas

A saída põe as três chamadas em três linhas para que o contraste caiba num
olhar: as duas primeiras (`f`) repetem o mesmo resultado com `value = 0`; a
terceira (`badF`) devolve **o mesmo número** e deixa `value = 7`. O retorno não
denuncia nada — o rastro sim.

Duas coisas propositais no código:

- **`badF` é anti-exemplo.** Não a torne pura: ela é o lado direito do slide.
- **Os resultados são guardados em variáveis antes de ler `value`.** A spec do
  Go só garante ordem de avaliação entre chamadas de função; ler o global
  direto na lista de argumentos do `Printf` poderia, legalmente, acontecer
  antes da chamada — e o palco veria o número errado.

E o custo aparece de novo no `main_test.go`: testar `f` não exige preparo nem
limpeza; testar `badF` exige zerar `value` para um teste não contaminar o
outro. Essa cerimônia é a impureza cobrando.
