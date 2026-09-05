# Lab 02 — estado compartilhado e ponteiros

> **O ponto:** o `*` que você **vê** não é o que machuca. `Len(&nome)` pede
> permissão de escrita e nem usa — mas pelo menos o `&` avisa. `SortedPrint(xs)`
> não tem um único `*` à vista e reordena o slice de quem chamou. A assinatura
> `func([]int)` não dá nenhuma pista.

**Go:** `1.21` — o Go de quando a palestra nasceu, e por isso pinado. O pacote
`slices` já é stdlib aqui (`slices.Clone`, `slices.Sort` entraram no 1.21), mas
o `slices.Sorted(slices.Values(x))` de hoje ainda não existe: isso é 1.23. A
variável de loop também ainda é uma só por `for`, compartilhada — o que muda no
1.22.

**Bloco da palestra:** 2

## Como rodar

```bash
make run       # a demo de palco
make version   # prova qual toolchain está rodando (deve dizer go1.21.13)
make test      # testes
make all       # vet + test
```

Ou, da raiz do repo:

```bash
make lab LAB=02-estado-global-ponteiros TARGET=run
```

## Notas

A saída tem dois blocos. No primeiro, `Len` e `LenValue` devolvem **14** — o
ponteiro não comprou nada: uma `string` já é passada como cabeçalho de 16 bytes,
os caracteres não são copiados. O que ele comprou foi um estado a mais, `nil`,
que vira panic (há um teste só para isso).

No segundo, `SortedPrintSafe` e `SortedPrint` imprimem **a mesma linha**,
`[1 2 3]`. A diferença aparece só na linha seguinte, quando o programa mostra o
slice do chamador: intacto num caso, reescrito no outro. É o contraste do
slide — a saída não denuncia, o rastro sim.

Detalhes propositais:

- **`Len` e `SortedPrint` são anti-exemplos.** Não os torne seguros: eles são o
  lado esquerdo do slide. O par com `LenValue` e `SortedPrintSafe` é o conteúdo.
- **`SortedPrintSafe` também muta** — só que um clone que nasce e morre dentro
  dela. Pureza não é "nunca mutar", é a mutação não escapar.
- **`sort.Ints` não é deprecado.** Do Go 1.22 em diante ele apenas chama
  `slices.Sort`. Está aqui porque é o que se escrevia em 2023, e porque a
  mutação é exatamente a mesma.
- **`GOTOOLCHAIN` está pinado no `Makefile`**, não em `auto` como nos outros
  labs: o `auto` só sobe de versão, então sem o pin o lab compilaria no 1.27 em
  modo de linguagem 1.21.

## Relação com `exemplos/00-basico-sobre-ponteiros`

O par `SortedPrint`/`SortedPrintSafe` é o mesmo daquele pacote — lá na versão
curada, testada por `Example`, para o slide. Aqui ele **roda**, no Go de 2023, e
vem acompanhado do caso do ponteiro explícito, que o exemplo curado não tem.
