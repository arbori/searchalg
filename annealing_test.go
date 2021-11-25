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

func (q *quadratic) Compute() float64 {
	return -quadratic_function(q)
}

func (q *quadratic) Reconfigure() {
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

func (q *quadratic) Assign(g Function) {
	f, ok := g.(*quadratic)

	if ok {
		q.a = f.a
		q.b = f.b
		q.c = f.c
		q.x = f.x
	}
}

func (q *quadratic) IsValid() bool {
	return true
}

func (q *quadratic) Clone() Function {
	result := quadratic{}

	result.a = q.a
	result.b = q.b
	result.c = q.c
	result.x = q.x

	return &result
}

func TestFoundMaximum(t *testing.T) {
	log.Println("TestFoundMaximum...")

	best_func := &quadratic{-2, 3, 2, -1}

	ctx := AnnealingContext{
		InitialTemperature: math.MaxFloat64,
		FinalTemperature:   10e-11,
		Cooling:            .001,
		Steps:              150,
		Deadline:           15,
	}

	SimulatedAnnealing(ctx, best_func)

	var maximumValue = 0.75

	log.Printf("Value expected to maximize f(x) is x = %.2f\n", maximumValue)
	log.Printf("Value founded for x  = %.9f\n", best_func.x)
	log.Printf("Error expected for |x - %.2f| = %.9f\n", maximumValue, float64(1e-5))
	log.Printf("Error founded  for |x - %.2f| = %.9f\n", maximumValue, math.Abs(best_func.x-maximumValue))

	if math.Abs(best_func.x-maximumValue) > float64(1e-5) {
		t.Fatalf("Annealing did not found the better value that supose to be 0.75")
	}
}
