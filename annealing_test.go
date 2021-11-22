package searchalg

import (
	"log"
	"math"
	"math/rand"
	"testing"
)

type quadratic struct {
	a, b, c float64
	x       float64
}

func quadratic_function(q *quadratic) float64 {
	return q.a*q.x*q.x + q.b*q.x + q.c
}

func (q *quadratic) compute() float64 {
	return -quadratic_function(q)
}

func (q *quadratic) reconfigure() {
	sign := 1.0

	if rand.Int()%2 != 0 {
		sign = -1.0
	}

	variation := (1.0 + sign*.1*rand.Float64())

	q.x = q.x * variation

	sign = 1.0

	if rand.Int()%2 != 0 {
		sign = -1.0
	}

	q.x = q.x + sign*rand.Float64()
}

func (q *quadratic) assign(g function) {
	f, ok := g.(*quadratic)

	if ok {
		q.a = f.a
		q.b = f.b
		q.c = f.c
		q.x = f.x
	}
}

func (q *quadratic) isValid() bool {
	return true
}

const maximumValue = 0.75

func TestFoundMaximum(t *testing.T) {
	log.Println("TestFoundMaximum...")

	best_func := &quadratic{-2, 3, 2, -1}
	last_func := &quadratic{-2, 3, 2, +1}

	searchAnnealing := annealing{
		temperaturaInicial: math.MaxFloat64,
		temperaturaFinal:   10e-11,
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

	searchAnnealing.run()

	log.Printf("Maximum value of f(x) = %.3f\n", best_func.x)

	if math.Abs(best_func.x-maximumValue) > float64(10e-5) {
		t.Fatalf("Annealing did not found the better value that supose to be 0.75")
	}

}

///////////////////////////////////////////////////////////////////////////////
func TestTimeTablingApproachMaximum(t *testing.T) {
	log.Println("TestTimeTablingApproachMaximum...")

	best_func := new_time_table(3, 5, 3, 5)
	last_func := new_time_table(3, 5, 3, 5)

	searchAnnealing := annealing{
		temperaturaInicial: math.MaxFloat64,
		temperaturaFinal:   10e-11,
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
		best:               &best_func,
		last:               &last_func,
	}

	startValue := best_func.f()

	log.Printf("Start value for time table: %.3f\n", startValue)

	searchAnnealing.run()

	log.Printf("Final value for time table: %.3f\n", best_func.f())

	if best_func.f() <= startValue {
		t.Fatalf("Annealing could not improve the time table proposed.")
	}
}
