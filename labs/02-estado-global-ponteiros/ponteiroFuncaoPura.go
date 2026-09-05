package main

import "fmt"

func main() {
	var nome = "Marco Ollivier"
	var nomeLen = Len(&nome)

	fmt.Println(nomeLen)
}

func Len(s *string) int {
	return len(*s)
}
