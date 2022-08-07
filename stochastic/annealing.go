package stochastic

import (
	"math"
	"math/rand"
	"time"
)

const BOLTZMAN_CONSTANT = 8.6173432e-5

// Interface that has behavior for a model used to search the optimum solution
type Function interface {
	// Return the value of objective function for the problem.
	Compute() float64
	// Reconfigure the value(s) of the model, representing other point in the solution domine.
	Reconfigure()
	// Copy the value and states from another histance, other point in the solution domine.
	Assign(f Function)
	// Check if configuration is a valid point in the solution domine.
	IsValid() bool
	Clone() Function
}

type AnnealingContext struct {
	// The temperature when the process start.
	InitialTemperature float64
	// The temperature when the process finish.
	FinalTemperature float64
	// Cooling percentage after each iteration.
	Cooling float64
	// Number of attempts to lower the temperature.
	Steps int
	// Time to finish, even final temperature could not be achived.
	Deadline int64
}

/*
Given an annealing context and two copies of the model that we want to
optimize, SimulatedAnnealing function run the simulation connecting the
annealing, a physical process, with the behavior of the model implemented
the Function interface.
*/
func SimulatedAnnealing(ctx AnnealingContext, best Function) {
	rand.Seed(time.Now().UnixNano())

	var passoAtual int
	var delta float64
	var sorteio float64
	var prob float64

	last := best.Clone()
	energiaInicial := best.Compute() // Pega a energia inicial do sistema
	energiaFinal := 0.0

	// Calcula o momento de termino T.
	fim := time.Now().Local().Add(time.Second * time.Duration(ctx.Deadline))

	// Processa o resfriamento do sistema.
	for temperatura := ctx.InitialTemperature; fim.After(time.Now()) && temperatura > ctx.FinalTemperature; temperatura = (1 - ctx.Cooling) * temperatura {
		// Busca uma configuração para a dada temperatura
		//um certo número de vezes.
		for passoAtual = ctx.Steps; passoAtual >= 0; passoAtual -= 1 {
			best.Reconfigure()

			// Calcula a energia atual do sistema.
			energiaFinal = best.Compute()
			// Calcula a variação de energia.
			delta = (energiaFinal - energiaInicial)

			// Compute some probability
			sorteio = rand.Float64()
			// Compute Boltzman probability
			prob = math.Exp((-1 * delta) / (BOLTZMAN_CONSTANT * temperatura))

			// Verifica se aceita a nova configuração.
			// Para uma nova configuração ser aceita além da variação de energia e da probabilidade
			// deve atender as restrições do problema.
			if (delta <= 0 || sorteio < prob) && best.IsValid() {
				if delta != 0 {
					energiaInicial = energiaFinal
					last.Assign(best)
				}
			} else {
				energiaFinal = energiaInicial
			}
		}

		best.Assign(last)

		passoAtual = ctx.Steps
	}
}
