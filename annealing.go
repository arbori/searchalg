package searchalg

import (
	"math"
	"math/rand"
	"time"
)

const BOLTZMAN_CONSTANT = 8.6173432e-5

type Function interface {
	compute() float64

	reconfigure()

	assign(f Function)

	isValid() bool
}

type Annealing struct {
	// Temperatura inicial do processo.
	TemperaturaInicial float64
	// Temperatura final do processo.
	TemperaturaFinal float64
	// Definição da temperatura máxima do processo de annealing.
	TemperaturaAtual float64
	// Taxa de Resfriamento do processo.
	Resfriamento float64

	// Número de tentativas de baixar a temperatura.
	Passos int
	// Indica qual o passo atual do processo.
	PassoAtual int

	// Variáveis de monitoração do processo
	EnergiaInicial float64
	EnergiaFinal   float64
	Delta          float64
	Sorteio        float64
	Prob           float64

	Prazo int64

	Best Function
	Last Function
}

/*
 Roda o algoritmo do simulated annealing, buscando a melhor configuração
 por um tempo determinado.
 Assume-se que foi passado para o objeto uma configuração inicial.
*/
func (ann Annealing) run() {
	rand.Seed(time.Now().UnixNano())

	ann.EnergiaInicial = ann.Best.compute() // Pega a energia inicial do sistema
	ann.EnergiaFinal = 0.0

	// Calcula o momento de termino T.
	fim := time.Now().Local().Add(time.Second * time.Duration(ann.Prazo))

	// Processa o resfriamento do sistema.
	for fim.After(time.Now()) {
		// Busca uma configuração para a dada temperatura
		//um certo número de vezes.
		for ann.PassoAtual = ann.Passos; ann.PassoAtual >= 0; ann.PassoAtual -= 1 {
			ann.Best.reconfigure()

			// Calcula a energia atual do sistema.
			ann.EnergiaFinal = ann.Best.compute()
			// Calcula a variação de energia.
			ann.Delta = (ann.EnergiaFinal - ann.EnergiaInicial)

			// Compute some probability
			ann.Sorteio = rand.Float64()
			// Compute Boltzman probability
			ann.Prob = math.Exp((-1 * ann.Delta) / (BOLTZMAN_CONSTANT * ann.TemperaturaAtual))

			// Verifica se aceita a nova configuração.
			// Para uma nova configuração ser aceita além da variação de energia e da probabilidade
			// deve atender as restrições do problema.
			if (ann.Delta <= 0 || ann.Sorteio < ann.Prob) && ann.Best.isValid() {
				if ann.Delta != 0 {
					ann.EnergiaInicial = ann.EnergiaFinal
					ann.Last.assign(ann.Best)
				}
			} else {
				ann.EnergiaFinal = ann.EnergiaInicial
			}
		}

		ann.Best.assign(ann.Last)

		ann.PassoAtual = ann.Passos
		ann.TemperaturaAtual = ann.Resfriamento * ann.TemperaturaAtual
	}
}
