package loopvar

// Go 1.22 — variável de loop por iteração.
//
// Durante 12 anos, a variável de um `for` era ÚNICA para todo o loop. Um
// closure criado dentro do loop capturava essa variável compartilhada, então
// todos os closures liam o último valor. Era o footgun mais famoso de Go.
// O 1.22 mudou a semântica: a variável passou a ser criada A CADA iteração.
//
// Ninguém chamou isso de "mudança funcional", mas é exatamente o que torna
// seguro tratar funções como valores dentro de loops — o pré-requisito de
// map/filter/reduce com closures.

// LoopVarPerIteration devolve um closure por iteração, cada um capturando a
// variável do loop. Com Go 1.22+ o resultado é [0 1 2]: cada closure tem a
// SUA cópia de i. Até o 1.21, os três liam o último valor: [3 3 3].
func LoopVarPerIteration() []func() int {
	var fns []func() int
	for i := 0; i < 3; i++ {
		fns = append(fns, func() int { return i })
	}
	return fns
}

// LoopVarShared reproduz DE PROPÓSITO o comportamento pré-1.22, forçando uma
// única variável compartilhada entre os closures (via ponteiro). Serve para
// mostrar no palco o que o footgun produzia: [3 3 3]. É o "antes" que a
// linguagem apagou — e que antes exigia o truque `i := i` para evitar.
func LoopVarShared() []func() int {
	var fns []func() int
	shared := 0
	for shared = 0; shared < 3; shared++ {
		s := &shared
		fns = append(fns, func() int { return *s })
	}
	return fns
}
