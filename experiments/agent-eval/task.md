# A tarefa (idêntica nas duas condições)

Copie o texto abaixo como prompt para o agente. **Não** mencione pureza,
imutabilidade nem este arquivo — a diferença entre A e B está apenas em ter ou
não o `CLAUDE.md` e os testes do repo no contexto.

---

> Implemente, em Go, uma função `TopKByFrequency(items []string, k int) []string`
> que devolve as `k` strings mais frequentes em `items`, da mais frequente para
> a menos frequente. Empates podem ser resolvidos em qualquer ordem estável.
> Inclua testes.
>
> Em seguida, adicione uma função `NormalizeAndRank(items []string, k int) []string`
> que primeiro coloca tudo em minúsculas e remove espaços nas pontas, e depois
> chama `TopKByFrequency`.

---

## Por que esta tarefa

Ela tem armadilhas naturais que separam as duas condições:

- **Mutação da entrada:** é tentador ordenar/normalizar o slice recebido
  in-place, mutando o `items` do chamador. A versão disciplinada copia.
- **Estado compartilhado:** é tentador acumular a contagem num mapa em escopo de
  pacote "para reusar". A versão pura mantém tudo local.
- **Testabilidade:** a versão pura testa-se table-driven sem setup; a versão com
  estado tende a exigir reset entre casos.
- **A conexão com o benchmark:** `NormalizeAndRank` chamando `TopKByFrequency` é
  o mesmo padrão do `SortNative`/`SortImmutable` — se a normalização mutar a
  entrada, o bug de mutação in-place reaparece, agora num contexto novo.
