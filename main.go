package main

import (
	"fmt"
	"strings"
	"bytes"
)

var (
	PalavraDoDia string
)

func main() {
	// De acordo com as estatisticas por letra as palavras SAREM e CULTO
	// contemplam 10 das fontes mais usadas (entre as 11 primeiras)
	// excluindo o I porque seria a unica vogal facilitando a descoberta na
	// palavras de 5 letras

	TestarXVezes := 1

	fmt.Println("Quer me contar a palavra do dia e duvida da minha honestidade? [ENTER] para não informar")
	fmt.Scanf("%s", &PalavraDoDia)

	fmt.Println("Primeiro vamos jogar tentando acertar sempre usando as estatisticas")
	mapAcertouEmModalidade1 := make(map[int]int)
	for i := 0; i < TestarXVezes; i++ {
		acertouEm, morte := AlgoritmoJoga(true)
		if !morte {
			mapAcertouEmModalidade1[acertouEm]++
		}
	}

	// Esse algoritmo não faz sentido testar porque pos chutes em todas as vezes serão os mesmo
	fmt.Println("Agora vamos jogar eliminando a maior quantidade de letras possiveis nos dois primeiros chutes")
	mapAcertouEmModalidade2 := make(map[int]int)
	acertouEm, morte := AlgoritmoJogaDoisChutesEliminatorios()
	if !morte {
		mapAcertouEmModalidade2[acertouEm] = TestarXVezes
	}

	fmt.Println("Agora jogar chutando sempre a primeira palavra")
	mapAcertouEmModalidade3 := make(map[int]int)
	for i := 0; i < TestarXVezes; i++ {
		acertouEm, morte := AlgoritmoJoga(false)
		if !morte {
			mapAcertouEmModalidade3[acertouEm]++
		}
	}

	if TestarXVezes > 1 {
		ImprimirAcertoEm(mapAcertouEmModalidade1, TestarXVezes)
		ImprimirAcertoEm(mapAcertouEmModalidade2, TestarXVezes)
		ImprimirAcertoEm(mapAcertouEmModalidade3, TestarXVezes)
	}
}

func AlgoritmoJoga(ComEstatistica bool) (int, bool) {
	game := new(adivinha)
	game.Carregar(5)

	var rodadas []string
	morte := true
	for turn := 1; turn <= 6; turn++ {
		meuChute := strings.ToUpper(game.ProximoChute(ComEstatistica))
		if meuChute == "" {
			fmt.Printf("Possibilidades: %+v\n", game.palavras)
			fmt.Println("Alguma coisa deu muito errado e eu não tenho mais paupites")
			break
		}
		fmt.Printf("Rodada (%d) Universo de possibilidades (%d)\n", turn, len(game.palavras))
		if len(game.palavras) <= 60 {
			fmt.Printf("Possibilidades: %+v\n", game.palavras)
		}
		fmt.Printf("Chute: %s\n", meuChute)

		var acertos string
		for len(acertos) != len(meuChute) {
			fmt.Printf("Resul: ")
			if PalavraDoDia == "" {
				fmt.Scanf("%s", &acertos)
			} else {
				acertos = Acertos(meuChute, PalavraDoDia)
				//fmt.Println(acertos)
				PrintChute(meuChute, acertos)
			}
		}
		rodadas = append(rodadas, acertos)
		if acertos == "22222" {
			morte = false
			break
		}
		game.Filtrar(meuChute, acertos)
	}
	ImprimeRodadas(rodadas)
	return len(rodadas), morte
}

/* Esse metodo não faz sentido sem estatisticas */
func AlgoritmoJogaDoisChutesEliminatorios() (int, bool) {
	ComEstatistica := true
	game := new(adivinha)
	game.Carregar(5)

	var rodadas []string
	morte := true
	meuChuteAnterior := ""
	for turn := 1; turn <= 6; turn++ {
		meuChute := strings.ToUpper(game.ProximoChute(ComEstatistica))
		if turn == 2 {
			gameExclude := new(adivinha)
			gameExclude.Carregar(game.tamanho)
			gameExclude.Filtrar(meuChuteAnterior, strings.Repeat("0", len(meuChuteAnterior)))
			meuChute = strings.ToUpper(gameExclude.ProximoChute(ComEstatistica))
		}

		if meuChute == "" {
			fmt.Printf("Possibilidades: %+v\n", game.palavras)
			fmt.Println("Alguma coisa deu muito errado e eu não tenho mais paupites")
			break
		}

		fmt.Printf("Rodada (%d) Universo de possibilidades (%d)\n", turn, len(game.palavras))
		if len(game.palavras) <= 60 {
			fmt.Printf("Possibilidades: %+v\n", game.palavras)
		}
		fmt.Printf("Chute: %s\n", meuChute)

		var acertos string
		for len(acertos) != len(meuChute) {
			fmt.Printf("Resul: ")
			if PalavraDoDia == "" {
				fmt.Scanf("%s", &acertos)
			} else {
				acertos = Acertos(meuChute, PalavraDoDia)
				//fmt.Println(acertos)
				PrintChute(meuChute, acertos)
			}
		}
		rodadas = append(rodadas, acertos)
		if acertos == "22222" {
			morte = false
			break
		}
		game.Filtrar(meuChute, acertos)
		meuChuteAnterior = meuChute
	}
	ImprimeRodadas(rodadas)
	return len(rodadas), morte
}

func Acertos(chute, palavraDoDia string) (string) {
	acertos := bytes.Repeat([]byte("0"), len(chute))
	for i, letra := range chute {
		if chute[i] == palavraDoDia[i] {
			acertos[i] = '2'
			continue
		}
		if strings.ContainsRune(palavraDoDia, letra) {
			acertos[i] = '1'
			continue
		}
	}
	return string(acertos)
}

func ImprimeRodadas(rodadas []string) {
	fmt.Println("===============")
	fmt.Printf("joguei %d/6\n", len(rodadas))

	for _, acertos := range rodadas {
		for _, acerto := range acertos {
			switch acerto {
			case '0':
				fmt.Printf("⬛️")
			case '1':
			 	fmt.Printf("🟨 ")
			case '2':
				fmt.Printf("🟩 ")
			}

		}
		fmt.Printf("\n")
	}
}

func PrintChute(palavra, acerto string) {
	BackgroundAmarelo := "\033[0;43m\033[1;30m"
	BackgroundVerde := "\033[0;42m\033[1;30m"
	Reset := "\033[0m"

	for k, letra := range palavra {
		Color := ""
		switch acerto[k] {
		case '1':
			Color = BackgroundAmarelo
		case '2':
			Color = BackgroundVerde
		}
		fmt.Printf(Color + string(letra) + Reset)
	}

	fmt.Printf("\n")
}

func ImprimirAcertoEm(m map[int]int, tentativas int) {
	contador := 0
	fmt.Println("Acertei em:")
	for i := 1; i <= 6; i++ {
		v, ok := m[i]
		if !ok {
			continue
		} 
		fmt.Printf("%02d: %03d\n", i, v)
		contador += v
	}
	if contador < tentativas {
		fmt.Printf("XX: %03d\n", tentativas - contador)
	}
}

func JogoDemo() {
	game := new(adivinha)
	game.Carregar(5)

	game.Filtrar("SAREM", "01110")
	game.Filtrar("CULTO", "00000")
	game.Filtrar("PEGAR", "01111")
	fmt.Printf("%d %+v\n", len(game.palavras), game.palavras)	
	game.Filtrar("AGARE", "01212")
	fmt.Printf("%d %+v\n", len(game.palavras), game.palavras)	
	game.Filtrar("GRADE", "22222")

//	fmt.Printf("%d %+v\n", len(game.palavras), game.palavras)	
}

func JogoMauricioDia20() {
	game := new(adivinha)
	game.Carregar(5)

	game.Filtrar("SAREI", "11010")
	game.Filtrar("CESTA", "02112")
	fmt.Printf("Chute: %+v\n", game.ProximoChute(true))

	fmt.Printf("%+v\n", game.palavras)	
}

func JogoMauricioDia19() {
	game := new(adivinha)
	game.Carregar(5)

	game.Filtrar("LETRA", "00000")
	game.Filtrar("FUNGO", "02022")

	fmt.Printf("%+v\n", game.palavras)	
}

func JogoMauricioDia18() {
	game := new(adivinha)
	game.Carregar(5)

	game.Filtrar("LETRA", "00211")
	game.Filtrar("FUNGO", "00101")

	fmt.Printf("%+v\n", game.palavras)
}

func JogoCesarDia18() {
	game := new(adivinha)
	game.Carregar(5)

	game.Filtrar("CISNE", "00010")
	game.Filtrar("PARTO", "01111")
	
	fmt.Printf("%+v\n", game.palavras)
}
