package main

import (
	"log"
)

func main() {
	dir := "." // Directorio actual

	log.Println("Iniciando watcher en:", dir)

	err := WatchDirectory(dir)
	if err != nil {
		log.Fatalf("Error iniciando watcher: %v", err)
	}
}
