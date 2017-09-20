package btlib

type MACD struct {
	Macd, Signal, Histogram float64
}

type BBPoint struct {
	Middle    float64
	Upper     float64
	Lower     float64
	BandWidth float64
}