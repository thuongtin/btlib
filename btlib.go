package btlib

import (
	"github.com/toorop/go-bittrex"
	"math"
)

type Btlib struct {
	client *bittrex.Bittrex
}

func (this Btlib) NewClient() *Btlib {
	bt := bittrex.New("", "")
	return &Btlib{bt}
}

func (this Btlib) HeikinAshi(candles []bittrex.Candle) {
	fistCandle := candles[0]
	fistCandle.Close = (candles[0].Low + candles[0].High + candles[0].Open + candles[0].Close) / 4
	fistCandle.Open = (candles[0].Open + candles[0].Close) / 2
	result := []bittrex.Candle{fistCandle}
	for i:=1; i<len(candles)-1; i++ {
		haCandle := bittrex.Candle{}
		//<b>HA-Close = (Open(0) + High(0) + Low(0) + Close(0)) / 4</b>
		haCandle.Close = (candles[i].Open + candles[i].High + candles[i].Low + candles[i].Close) / 4
		//<b>HA-Open = (HA-Open(-1) + HA-Close(-1)) / 2</b>
		haCandle.Open = (result[i-1].Open + result[i-1].Close) / 2
		//<b>HA-High = Maximum of the High(0), HA-Open(0) or HA-Close(0) </b>
		haCandle.High = math.Max(math.Max(haCandle.Close, haCandle.Open), candles[i].High)
		//<b>HA-Low = Minimum of the Low(0), HA-Open(0) or HA-Close(0) </b>
		haCandle.Low = math.Min(math.Min(haCandle.Close, haCandle.Open), candles[i].Low)
	}
}