package main

import (
	"testing"
	_ "fmt"
)

func TestFiltrar(t *testing.T) {
	return
	game := new(adivinha)
	game.Carregar(5)
	game.palavras = []string{"cursa", "farsa", "morsa", "saram", "sarar", "sarca", "sarda", "sarna", "sarou", "sorta", "sorva", "surda", "surja", "surra", "surta"}
	game.Filtrar("SAROU", "22200")

	wanted := false
	for _, p := range game.palavras {
		if p == "sarar" {
			wanted = true
		}
	}
	
	if !wanted {
		t.Fatalf("O filtro não permaneceu com a palavra sarar")
	}
}

func TestProximoChute(t *testing.T) {
	game := new(adivinha)
	game.Carregar(5)
	game.palavras = []string{"saram", "sarar", "sarca", "sarda", "sarna"}

	meuChute := game.ProximoChute(true)
	if meuChute == "" {
		t.Fatalf("Ficou sem ter o que chutar")
	}
}
