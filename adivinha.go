package main

import (
	"strings"
	"os"
	"bufio"
	"io"
	"sort"
	"math/bits"
	"math/rand"
	"time"
	_ "fmt"
)

type adivinha struct {
	tamanho  uint
	palavras []string
}

type TiposFiltro int

const (
	NAOCONTEM TiposFiltro = iota
	CONTEM
	NAOPERTENCE
	PERTENCE
)

func (a *adivinha) Carregar(tamanho uint) error {
	a.tamanho = tamanho
	a.palavras = []string{}

	file, err := os.Open("palavras5letras.txt")

	if err != nil {
		return err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	
	for {
		palavra, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		a.palavras = append(a.palavras, strings.ToLower(string(palavra)))
	}
	return nil
}

func (a *adivinha) ProximoChute(consideraEstatisticas bool) string {
	if len(a.palavras) == 0 {
		return ""
	}
	if len(a.palavras) == 1 {
		return a.palavras[0]
	}

	if consideraEstatisticas {
		letras := a.ocorrenciaLetras()
		for universoLetras := int(a.tamanho); universoLetras <= len(letras); universoLetras++ {
			chutes := a.combinacoes(letras[0:universoLetras])
		
			mapChutes := make(map[string]bool, len(chutes))
			for _, chute := range chutes {
				mapChutes[chute] = true
			}

			for _, palavra := range a.palavras {
				if _, ok := mapChutes[palavra]; ok {
					return palavra
				}
			}
		}
	}

	// Vamos fazer um loop para evitar chutes com letras repetidas
	var palavrasParaChute []string
	for _, palavra := range a.palavras {
		repeteLetras := false
		for _, letra := range palavra {
			if strings.Count(palavra, string(letra)) > 1 {
				repeteLetras = true
				break
			}
		}
		if !repeteLetras {
			palavrasParaChute = append(palavrasParaChute, palavra)
		}	
	}

	if len(palavrasParaChute) == 0 {
		palavrasParaChute = make([]string, len(a.palavras))
		copy(palavrasParaChute, a.palavras)
	}

	rand.Seed(time.Now().UnixNano())
	return palavrasParaChute[rand.Intn(len(palavrasParaChute))]
}

func (a *adivinha) ocorrenciaLetras() []byte {
	type contador struct {
		Letra byte
		Total int
	}
	var probabilidades []contador

	alfabeto := "abcdefghijklmnopqrstuvxwyz"
	for _, lt := range alfabeto {
		letra := string(lt)[0]
		count := 0
		for _, palavra := range a.palavras {
			if strings.IndexByte(palavra, letra) != -1 {
				count++
			}
		}
		probabilidades = append(probabilidades, contador{Letra: letra, Total: count})
	}
	// Ordena de forma decrescente
	sort.SliceStable(probabilidades, func(i, j int) bool {
		 return probabilidades[i].Total > probabilidades[j].Total
	})

	var ret []byte
	for _, prob := range probabilidades {
		if prob.Total > 0 {
			ret = append(ret, prob.Letra)
			//fmt.Printf("%s %d\n", strings.ToUpper(string(prob.Letra)), prob.Total)
		}
	}
	return ret
}

func (a *adivinha) combinacoes(opcoes []byte) []string {
	if len(opcoes) < int(a.tamanho) {
		return []string{}
	}

	var combo []string

	MAX := uint((1 << len(opcoes))-1)

	for NUM := uint(1); MAX >= NUM; NUM++ {
		if bits.OnesCount(NUM) == int(a.tamanho) {
			game := new(adivinha)
			game.tamanho = a.tamanho
			game.palavras = make([]string, len(a.palavras))
			copy(game.palavras, a.palavras)
			game.Filtrar(a.stringCombinacaoSemRepetir(opcoes, NUM), strings.Repeat("X", int(a.tamanho)))

			combo = append(combo, game.palavras...)
		}
	}

	return combo
}

func (a *adivinha) stringCombinacaoSemRepetir(opcoes []byte, NUM uint) string {
	str := ""
	for i, letra := range opcoes {
		bit := uint(1 << i)
		if bit & NUM > 0 {
			str += string(letra)
		}
	}
	return str
}

func (a *adivinha) Filtrar(palavra, resultado string) {
	palavra = strings.ToLower(palavra)

	for pos, regra := range resultado {
		var tpFiltro TiposFiltro
		switch regra {
		case '0':
			tpFiltro = NAOCONTEM
		case '1':
			tpFiltro = NAOPERTENCE
		case '2':
			tpFiltro = PERTENCE
		case 'X': // Caso especifico para filtrar internamente chutes possiveis
			tpFiltro = CONTEM
		}
		a.filtrarInterno(tpFiltro, palavra[pos], pos+1)
	}
} 

func (a *adivinha) filtrarInterno(tp TiposFiltro, lt byte, posicao int) {
	var palavras []string

	for _, palavra := range a.palavras {
		switch tp {
		case NAOCONTEM:
			if strings.IndexByte(palavra, lt) == -1 {
				palavras = append(palavras, palavra)
			}
		case CONTEM:
			if strings.IndexByte(palavra, lt) != -1 {
				palavras = append(palavras, palavra)
			}
		case NAOPERTENCE:
			if strings.IndexByte(palavra, lt) != -1 && palavra[posicao-1] != lt {
				palavras = append(palavras, palavra)
			}
		case PERTENCE:
			if palavra[posicao-1] == lt {
				palavras = append(palavras, palavra)
			}
		}
	}

	a.palavras = palavras
}
