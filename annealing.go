package searchalg

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
)

const Kb = 8.6173432e-5

type function interface {
	compute() float64

	reconfigure()

	assign(f function)

	isValid() bool
}

type annealing struct {
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

	best function
	last function
}

/*
 Roda o algoritmo do simulated annealing, buscando a melhor configuração
 por um tempo determinado.
 Assume-se que foi passado para o objeto uma configuração inicial.
*/
func (ann annealing) run() {
	rand.Seed(time.Now().UnixNano())

	ann.energiaInicial = ann.best.compute() // Pega a energia inicial do sistema
	ann.energiaFinal = 0.0

	// Calcula o momento de termino T.
	fim := time.Now().Local().Add(time.Second * time.Duration(ann.prazo))

	fmt.Println(&ann.best)
	fmt.Println(ann.best)

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
			ann.prob = math.Exp((-1 * ann.delta) / (Kb * ann.temperaturaAtual))

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

	fmt.Println(ann.best)
}

////////////////////////////////////////////////////////////////////////////////////

type time_table struct {
	schedule_limits [4]int
	schedule        [][][][]int
}

func (tt *time_table) compute() float64 {
	var disciplina int
	var sala int
	var horario int
	var dia_semana int

	result := 0.0

	for disciplina = 0; disciplina < tt.schedule_limits[0]; disciplina += 1 {
		for sala = 0; sala < tt.schedule_limits[1]; sala += 1 {
			for horario = 0; horario < tt.schedule_limits[2]; horario += 1 {
				for dia_semana = 0; dia_semana < tt.schedule_limits[3]; dia_semana += 1 {
					result += float64(tt.schedule[disciplina][sala][horario][dia_semana])
				}
			}
		}
	}

	return result
}

func (tt *time_table) reconfigure() {
	disciplina := rand.Intn(tt.schedule_limits[0])
	sala := rand.Intn(tt.schedule_limits[1])
	horario := rand.Intn(tt.schedule_limits[2])
	dia_semana := rand.Intn(tt.schedule_limits[3])

	tt.schedule[disciplina][sala][horario][dia_semana] = 1 - tt.schedule[disciplina][sala][horario][dia_semana]
}

func (tt *time_table) assign(f function) {
	val := f.(time_table)

	var disciplina int
	var sala int
	var horario int
	var dia_semana int

	tt.schedule_limits[0] = val.schedule_limits[0]
	tt.schedule_limits[1] = val.schedule_limits[1]
	tt.schedule_limits[2] = val.schedule_limits[2]
	tt.schedule_limits[3] = val.schedule_limits[3]

	for disciplina = 0; disciplina < tt.schedule_limits[0]; disciplina += 1 {
		for sala = 0; sala < tt.schedule_limits[1]; sala += 1 {
			for horario = 0; horario < tt.schedule_limits[2]; horario += 1 {
				for dia_semana = 0; dia_semana < tt.schedule_limits[3]; dia_semana += 1 {
					tt.schedule[disciplina][sala][horario][dia_semana] = val.schedule[disciplina][sala][horario][dia_semana]
				}
			}
		}
	}
}

func (tt *time_table) isValid() bool {
	var disciplina int
	var sala int
	var horario int
	var dia_semana int

	disciplina_duplicada := 0

	for sala = 0; sala < tt.schedule_limits[1]; sala += 1 {
		for horario = 0; horario < tt.schedule_limits[2]; horario += 1 {
			for dia_semana = 0; dia_semana < tt.schedule_limits[3]; dia_semana += 1 {
				for disciplina = 0; disciplina < tt.schedule_limits[0]; disciplina += 1 {
					disciplina_duplicada += tt.schedule[disciplina][sala][horario][dia_semana]
				}

				if disciplina_duplicada > 1 {
					return false
				}
			}
		}
	}

	sala_com_duas_disciplinas := 0

	for horario = 0; horario < tt.schedule_limits[2]; horario += 1 {
		for dia_semana = 0; dia_semana < tt.schedule_limits[3]; dia_semana += 1 {
			for disciplina = 0; disciplina < tt.schedule_limits[0]; disciplina += 1 {
				for sala = 0; sala < tt.schedule_limits[1]; sala += 1 {
					sala_com_duas_disciplinas += tt.schedule[disciplina][sala][horario][dia_semana]
				}

				if sala_com_duas_disciplinas > 1 {
					return false
				}
			}
		}
	}

	return true
}

////////////////////////////////////////////////////////////////////////////////////
func  main()  {
	best_func := &time_table{
		schedule_limits [4]int {3, 5, 3, 5}
		schedule        [][][][]int {
			{}
		}
	
	}
	last_func := &quadratic{-2, 3, 2, +1}

	searchAnnealing := annealing{
		temperaturaInicial: math.MaxFloat64,
		temperaturaFinal:   10E-11,
		temperaturaAtual:   math.MaxFloat64,
		resfriamento:       (1 - .05),
		passos:             100,
		passoAtual:         100,
		energiaInicial:     0,
		energiaFinal:       0,
		delta:              0,
		sorteio:            0,
		prob:               0,
		prazo:              10,
		best:               best_func,
		last:               last_func,
	}

	fmt.Println(&searchAnnealing.best)
	fmt.Println(searchAnnealing.best)

	/*
		fim := time.Now().Local().Add(time.Second * time.Duration(1))

		fmt.Println(fim)
		for fim.After(time.Now()) {
			fmt.Println(time.Now())
		}
	*/

	searchAnnealing.run()

	fmt.Fprintf(os.Stdout, "Maximum value of f(x) = %.3f\n", best_func.x)

	if math.Abs(best_func.x-maximumValue) > float64(10e-5) {
		t.Fatalf("Annealing did not found the better value that supose to be 0.75")
	}
}