package btlib

import (
	"github.com/toorop/go-bittrex"
	"math"
	"github.com/murlokswarm/errors"
	"fmt"
)

type Btlib struct {
	client *bittrex.Bittrex
}

func (this Btlib) NewClient() *Btlib {
	bt := bittrex.New("", "")
	return &Btlib{bt}
}

func (this Btlib) HeikinAshi(candles []bittrex.Candle) ([]bittrex.Candle, error) {
	if len(candles) == 0 {
		return nil, nil
	}
	fistCandle := candles[0]
	fistCandle.Close = (candles[0].Low + candles[0].High + candles[0].Open + candles[0].Close) / 4
	fistCandle.Open = (candles[0].Open + candles[0].Close) / 2
	result := []bittrex.Candle{fistCandle}
	for i:=1; i<len(candles); i++ {
		haCandle := candles[i]
		//<b>HA-Close = (Open(0) + High(0) + Low(0) + Close(0)) / 4</b>
		haCandle.Close = (candles[i].Open + candles[i].High + candles[i].Low + candles[i].Close) / 4
		//<b>HA-Open = (HA-Open(-1) + HA-Close(-1)) / 2</b>
		haCandle.Open = (result[i-1].Open + result[i-1].Close) / 2
		//<b>HA-High = Maximum of the High(0), HA-Open(0) or HA-Close(0) </b>
		haCandle.High = math.Max(math.Max(haCandle.Close, haCandle.Open), candles[i].High)
		//<b>HA-Low = Minimum of the Low(0), HA-Open(0) or HA-Close(0) </b>
		haCandle.Low = math.Min(math.Min(haCandle.Close, haCandle.Open), candles[i].Low)
		result = append(result, haCandle)
	}
	return result, nil
}

func (this *Btlib) GetCandles(pair, interVal string) ([]bittrex.Candle, error) {
	candles, err := this.client.GetTicks(pair, interVal)
	if err != nil {
		return nil, nil
	}
	return candles, nil
}

func (this *Btlib) GetLimitCandles(pair, interVal string, limit int) ([]bittrex.Candle, error) {
	candles, err := this.client.GetTicks(pair, interVal)
	if err != nil {
		return nil, nil
	}
	cLen := len(candles)
	if cLen > limit {
		candles = candles[cLen-limit:]
	}
	return candles, nil
}

func (this *Btlib) MA(candles []bittrex.Candle, length int) ([]float64, error) {
	//Giá đóng cửa trung bình của length nến
	result := []float64{}
	l := float64(length)
	cLen := len(candles)
	if cLen <= length {
		return nil, errors.New("Out of range")
	}
	a := candles[0].Close
	b := 0.0
	for i:=0; i<length; i++ {
		b += candles[i].Close
	}

	for i := length-1; i<cLen-1; i++ {
		result = append(result, b/l)
		b = b - a + candles[i+1].Close
		a = candles[i-length+2].Close
	}
	return result, nil
}

func (this *Btlib) maOfMacd(macd []float64) []float64 {
	result := []float64{}

	a := macd[0]
	b := 0.0
	for i:=0; i<9; i++ {
		b += macd[i]
	}

	for i := 8; i<len(macd)-1; i++ {
		result = append(result, b/9)
		b = b - a + macd[i+1]
		a = macd[i-7]
	}
	return result
}

func (this *Btlib) Vol(candles []bittrex.Candle, maLength int) ([]float64, error) {
	result := []float64{}
	l := float64(maLength)
	cLen := len(candles)
	if cLen <= maLength {
		return nil, errors.New("Out of range")
	}
	a := candles[0].BaseVolume
	b := 0.0
	for i:=0; i<maLength; i++ {
		b += candles[i].BaseVolume
	}
	for i := maLength-1; i<cLen-1; i++ {
		result = append(result, b/l)
		b = b - a + candles[i+1].BaseVolume
		a = candles[i-maLength+2].BaseVolume
	}
	return result, nil
}

//func (this *Btlib) Ema(candles []bittrex.Candle, length int) ([]float64, error) {
//	var smoothing float64
//	result := []float64{}
//	l := float64(length)
//	smoothing = 2/(l+1)
//	cLen := len(candles)
//	if cLen <= length {
//		return nil, errors.New("Out of range")
//	}
//
//	sma,_ := this.MA(candles, length)
//	result = append(result, sma[0])
//	for i:=1; i<=len(sma); i++ {
//		result = append(result, smoothing * (candles[i+length-2].Close - result[i-1]) + result[i-1])
//	}
//	return result, nil
//}

func (this *Btlib) Ema(candles []bittrex.Candle, length int) ([]float64, error) {
	result := []float64{}
	l := float64(length)
	cLen := len(candles)
	if cLen <= length {
		return nil, errors.New("Out of range")
	}
	a := 0.0
	for i, b := range candles {
		if i < length - 1 {
			a += b.Close
		} else if i == length-1 {
			a += b.Close
			result = append(result, a/l)
		} else {
			r := (b.Close*(2/(l+1))) + result[i-length]*(1-(2/(l+1)))
			result = append(result, r)
		}
	}
	return result, nil
}


func (this *Btlib) emaOfMacd(macd []float64) []float64 {
	result := []float64{}
	a := 0.0
	for i, b := range macd {
		if i < 8 {
			a += b
		} else if i == 8 {
			a += b
			result = append(result, a/9)
		} else {
			r := (b*(2/10)) + result[i-9]*(1-(2/10))
			result = append(result, r)
		}
	}
	return result
}

func (this *Btlib) Macd(candles []bittrex.Candle) ([]MACD, error) {
	result := []MACD{}
	cLen := len(candles)
	if cLen <= 26 {
		return nil, errors.New("Out of range")
	}
	ema12, _ := this.Ema(candles, 12)
	ema26, _ := this.Ema(candles, 26)

	ema12 = ema12[14:]
	macd := []float64{}

	for i, _ := range ema26 {
		v := ema12[i] - ema26[i]
		macd = append(macd, v)
		fmt.Println(v)
	}

	signals := this.emaOfMacd(macd)
	macd = macd[8:]

	for i, _ := range macd {
		result = append(result, MACD{
			macd[i],
			signals[i],
			macd[i] - signals[i]})
	}

	return result, nil
}
