package btlib

import (
	"github.com/toorop/go-bittrex"
	"math"
	"github.com/murlokswarm/errors"
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