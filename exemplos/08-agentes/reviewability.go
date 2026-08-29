package agentes

// Bloco "paradigma na era dos agentes": o contraste de revisibilidade.
//
// As duas funções abaixo calculam a MESMA coisa — a média móvel simples de
// uma janela sobre uma série. Uma é pura; a outra depende de estado
// compartilhado. O argumento da palestra: quando o volume de código gerado
// sobe, o gargalo vira a REVISÃO, e função pura é desproporcionalmente mais
// barata de revisar — você lê a assinatura e sabe o universo do que pode
// acontecer.

// MovingAveragePure é pura: entrada -> saída, sem estado externo. Para revisar
// basta esta função. Sem mock, sem setup: o teste é table-driven direto. Se um
// agente a altera, o raio de dano é local — nada fora daqui pode quebrar.
func MovingAveragePure(series []float64, window int) []float64 {
	if window <= 0 || len(series) < window {
		return nil
	}
	out := make([]float64, 0, len(series)-window+1)
	var sum float64
	for i := 0; i < window; i++ {
		sum += series[i]
	}
	out = append(out, sum/float64(window))
	for i := window; i < len(series); i++ {
		sum += series[i] - series[i-window]
		out = append(out, sum/float64(window))
	}
	return out
}

// accumulator é o estado compartilhado que torna a versão abaixo difícil de
// revisar: para saber se MovingAverageStateful está certa, você precisa ler
// TAMBÉM quem mais toca neste pacote-global e em que ordem.
var accumulator float64

// MovingAverageStateful calcula o mesmo, mas apoiado num acumulador global.
// A assinatura MENTE: ela promete depender só de series/window, mas o
// resultado depende do valor herdado de `accumulator` e de quem chamou antes.
// É a classe de código que um agente altera "corretamente" no local e quebra
// à distância — e que passa despercebida numa revisão apressada.
func MovingAverageStateful(series []float64, window int) []float64 {
	if window <= 0 || len(series) < window {
		return nil
	}
	out := make([]float64, 0, len(series)-window+1)
	for i := 0; i+window <= len(series); i++ {
		accumulator = 0 // fácil de esquecer ao editar — e aí o bug aparece longe
		for j := i; j < i+window; j++ {
			accumulator += series[j]
		}
		out = append(out, accumulator/float64(window))
	}
	return out
}
