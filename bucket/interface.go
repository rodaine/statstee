package bucket

type Interface interface {
	Add(v float64)

	Sum() float64
	Freq() float64
	Unique() float64

	Last() float64

	Mean() float64
	Min() float64
	Max() float64

	Median() float64
	P75() float64
	P95() float64
	P99() float64

	Reset()
}
