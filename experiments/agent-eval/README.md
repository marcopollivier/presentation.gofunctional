# Experimento: paradigma como especificação para agentes

Objetivo: transformar o Bloco 6 da palestra ("paradigma na era dos agentes") de
opinião em **dado verificável**. A mesma tarefa é pedida a um agente em duas
condições e as saídas são comparadas por critérios objetivos.

> **Não fabrique resultados.** Esta pasta é só o protocolo e o formato de
> registro. `results/` começa vazia; o autor roda de verdade e preenche.

## Condições

- **A — sem restrição de estilo.** Prompt cru (`task.md`), sem `CLAUDE.md` no
  contexto, sem apontar os testes existentes.
- **B — com as regras.** Mesmo prompt, com o `CLAUDE.md` da raiz no contexto e
  os testes do repo disponíveis como sinal automático.

## Passos

1. Rodar a Condição A num ambiente limpo; salvar a saída em
   `results/<data>-A/`.
2. Rodar a Condição B; salvar em `results/<data>-B/`.
3. Aplicar `criteria.md` às duas saídas; preencher a tabela em cada pasta.
4. Levar a comparação para o slide.

## Hipótese

A Condição B reduz efeitos colaterais, elimina a necessidade de mock nos testes,
aloca menos e faz o bug de mutação in-place **não aparecer**. Se os dados
refutarem isso, o slide diz isso — o experimento vale pela honestidade, não pela
confirmação.

## Formato de `results/<data>-<cond>/`

```
results/2026-11-XX-A/
  prompt.txt      # o prompt exato usado (inclui contexto, se houver)
  output/         # o código gerado pelo agente
  criteria.md     # a tabela de criteria.md preenchida para esta saída
  notes.md        # observações qualitativas
```
