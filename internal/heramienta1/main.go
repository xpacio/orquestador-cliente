package main

import (
	"fmt"
	"herramienta1" // Importas la lógica de la herramienta
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Se necesita un argumento")
		os.Exit(1)
	}

	argumento := os.Args[1]
	resultado := herramienta1.Procesar(argumento)
	fmt.Println(resultado)
}
