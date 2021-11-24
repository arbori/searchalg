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
	temperaturaInicial float64
	// Temperatura final do processo.
	temperaturaFinal float64
	// Definição da temperatura máxima do processo de annealing.
	temperaturaAtual float64
	// Taxa de resfriamento do processo.
	resfriamento float64

	// Número de tentativas de baixar a temperatura.
	passos int
	// Indica qual o passo atual do processo.
	passoAtual int

	// Variáveis de monitoração do processo
	energiaInicial float64
	energiaFinal   float64
	delta          float64
	sorteio        float64
	prob           float64

	prazo int64

	best Function
	last Function
}

/*
 Roda o algoritmo do simulated annealing, buscando a melhor configuração
 por um tempo determinado.
 Assume-se que foi passado para o objeto uma configuração inicial.
*/
func (ann Annealing) run() {
	rand.Seed(time.Now().UnixNano())

	ann.energiaInicial = ann.best.compute() // Pega a energia inicial do sistema
	ann.energiaFinal = 0.0

	// Calcula o momento de termino T.
	fim := time.Now().Local().Add(time.Second * time.Duration(ann.prazo))

	// Processa o resfriamento do sistema.
	for fim.After(time.Now()) {
		// Busca uma configuração para a dada temperatura
		//um certo número de vezes.
		for ann.passoAtual = ann.passos; ann.passoAtual >= 0; ann.passoAtual -= 1 {
			ann.best.reconfigure()

			// Calcula a energia atual do sistema.
			ann.energiaFinal = ann.best.compute()
			// Calcula a variação de energia.
			ann.delta = (ann.energiaFinal - ann.energiaInicial)

			// Compute some probability
			ann.sorteio = rand.Float64()
			// Compute Boltzman probability
			ann.prob = math.Exp((-1 * ann.delta) / (BOLTZMAN_CONSTANT * ann.temperaturaAtual))

			// Verifica se aceita a nova configuração.
			// Para uma nova configuração ser aceita além da variação de energia e da probabilidade
			// deve atender as restrições do problema.
			if (ann.delta <= 0 || ann.sorteio < ann.prob) && ann.best.isValid() {
				if ann.delta != 0 {
					ann.energiaInicial = ann.energiaFinal
					ann.last.assign(ann.best)
				}
			} else {
				ann.energiaFinal = ann.energiaInicial
			}
		}

		ann.best.assign(ann.last)

		ann.passoAtual = ann.passos
		ann.temperaturaAtual = ann.resfriamento * ann.temperaturaAtual
	}
}
