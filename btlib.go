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
			r := (b.Close*(2/(l+1))) + result[i-length]*(1-2/(l+1))
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
			r := (b*(2/10)) + result[i-9]*(1-2/10)
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
		//fmt.Println(v)
		macd = append(macd, v)
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


func (this *Btlib) BB(candles []bittrex.Candle, length int, stdDev float64) []BBPoint {
	len := len(candles)
	if len < length {
		return nil
	}
	result := []BBPoint{}
	for i := length - 1; i < len; i++ {
		point := calcBBPoint(candles[i-(length-1):i+1], length, stdDev)
		result = append(result, point)
	}
	return result
}

func calcBBPoint(candles []bittrex.Candle, length int, stdDev float64) BBPoint {
	size := float64(len(candles))
	if size != float64(length) {
		return BBPoint{}
	}
	total := 0.0
	for _, candle := range candles {
		total += candle.Close
	}
	middle := total / size
	total = 0.0
	for _, candle := range candles {
		total += math.Pow(candle.Close-middle, 2)
	}
	sd := math.Sqrt(total / size)
	upper := middle + sd*stdDev
	lower := middle - sd*stdDev
	bandwidth := upper - lower
	return BBPoint{middle, upper, lower, bandwidth}
}

func (this *Btlib) RSI(candles []bittrex.Candle, length int) []float64 {
	len := len(candles)
	if len < length {
		return nil
	}
	result := []float64{}
	sumGain := 0.0
	sumLoss := 0.0
	rs := 0.0
	for i := 1; i < len; i++ {
		preClose := candles[i-1].Close
		close := candles[i].Close
		change := close - preClose
		gain := 0.0
		loss := 0.0
		if change >= 0 {
			gain = change
		} else {
			loss = change * (-1.0)
		}

		if i < length - 1 {
			sumGain += gain
			sumLoss += loss
		} else {
			if i == length - 1 {
				sumGain = (sumGain + gain) / 14.0
				sumLoss = (sumLoss + loss) / 14.0
			} else {
				sumGain = (sumGain*13 + gain) / 14.0
				sumLoss = (sumLoss*13 + loss) / 14.0
			}

			if sumLoss == 0 {
				result = append(result, 100)
			} else {
				rs = sumGain / sumLoss
				result = append(result, 100 - (100 / (rs + 1)))
			}
		}
	}
	return result
}

func (this *Btlib) ADX(cd []bittrex.Candle, length int) []ADXPoint {
	size := len(cd)
	if size < length {
		return nil
	}
	result := []ADXPoint{}
	trSum := 0.0
	dmPlusSum := 0.0
	dmMinusSum := 0.0
	bv := []float64{}

	for i := 1; i < size; i++ {
		high := cd[i].High
		preHigh := cd[i-1].High
		low := cd[i].Low
		preLow := cd[i-1].Low
		preClose := cd[i-1].Close

		tr := getMax(high-low, math.Abs(high-preClose), math.Abs(low-preClose))
		dmPlus := 0.0
		if high-preHigh > preLow-low {
			dmPlus = math.Max(high-preHigh, 0.0)
		}
		dmMinus := 0.0
		if preLow-low > high-preHigh {
			dmMinus = math.Max(preLow-low, 0.0)
		}

		if i < length-1 {
			trSum += tr
			dmPlusSum += dmPlus
			dmMinusSum += dmMinus
		} else {
			if i == length-1 {
				trSum += tr
				dmPlusSum += dmPlus
				dmMinusSum += dmMinus
			} else {
				trSum = trSum - (trSum / float64(length)) + tr
				dmPlusSum = dmPlusSum - (dmPlusSum / float64(length)) + dmPlus
				dmMinusSum = dmMinusSum - (dmMinusSum / float64(length)) + dmMinus
			}
			diPlus := 100 * (dmPlusSum / trSum)
			diMinus := 100 * (dmMinusSum / trSum)
			diDiff := math.Abs(diPlus - diMinus)
			diSum := diPlus + diMinus
			dx := 100 * (diDiff / diSum)
			bv = append(bv, dx)
			result = append(result, ADXPoint{0, diPlus, diMinus})

			/*if i < length*2-1 {
				dxSum += dx
			} else if i == length*2-1 {
				adx = (dxSum + dx) / float64(length)
				result = append(result, ADXPoint{adx, diPlus, diMinus})
			} else {
				adx = (adx*float64(length-1) + dx) / float64(length)
				result = append(result, ADXPoint{adx, diPlus, diMinus})
			}*/
		}
	}
	adx := getSMAOfRSI(length, bv)
	for i := 0; i < len(adx); i++ {
		result[len(result)-1-i].ADX = adx[len(adx)-1-i]
	}
	return result
}

func getMax(elems ...float64) float64 {
	const MinFloat = float64(math.MinInt64)
	max := MinFloat
	for _, e := range elems {
		if max < e {
			max = e
		}
	}
	return max
}



func getSMAOfRSI(length int, rsiArray []float64) []float64 {
	size := len(rsiArray)
	if size < length {
		return nil
	}
	result := []float64{}
	for i := length - 1; i < size; i++ {
		sum := 0.0
		for j := i - (length - 1); j <= i; j++ {
			sum += rsiArray[j]
		}
		result = append(result, sum/float64(length))
	}
	return result
}

func getBBAroundSMAArray(rsiArray []float64, bandLength int, stdDev float64) []TDIPoint {
	size := len(rsiArray)
	if size < bandLength {
		return nil
	}
	result := []TDIPoint{}
	for i := bandLength - 1; i < size; i++ {
		point := getBBAroundSMAPoint(rsiArray[i-(bandLength-1):i+1], bandLength, stdDev)
		result = append(result, point)
	}
	return result
}

func getBBAroundSMAPoint(rsiArray []float64, bandLength int, stdDev float64) TDIPoint {
	size := float64(len(rsiArray))
	if size != float64(bandLength) {
		return TDIPoint{}
	}
	total := 0.0
	for _, rsiPoint := range rsiArray {
		total += rsiPoint
	}
	average := total / size
	total = 0.0
	for _, rsiPoint := range rsiArray {
		total += math.Pow(rsiPoint-average, 2)
	}
	sd := math.Sqrt(total / size)
	upper := average + sd*stdDev
	lower := average - sd*stdDev
	middle:= (upper+lower)/2
	return TDIPoint{middle, upper, lower, 0,0}
}

func (this *Btlib) TDI(candles []bittrex.Candle, rsiPeriod, bandLength, fast, slow int) []TDIPoint {
	rsi := this.RSI(candles, rsiPeriod)
	result := getBBAroundSMAArray(rsi, bandLength, 1.6185)
	fastArray := getSMAOfRSI(fast, rsi)
	slowArray := getSMAOfRSI(slow, rsi)

	for i := 0; i < len(result); i++{
		result[len(result) - 1 - i].FastMA = fastArray[len(fastArray) - 1 - i]
		result[len(result) - 1 - i].SlowMA = slowArray[len(slowArray) - 1 - i]
	}
	return result
}
func (this *Btlib) Ichimoku(candles []bittrex.Candle) []IchimokuCloud {
	result := []IchimokuCloud{}
	lowestLow := candles[0].Low
	highestHigh := candles[0].High

	for i:=0; i<len(candles); i++ {
		tenkan := 0.0
		kijun := 0.0
		chikou := 0.0
		senkouA := 0.0
		if i < len(candles) - 26 {
			chikou = candles[i+26].Close
		}
		if i >= 8 {
			highestHigh = findMax(candles[i-8:i])
			lowestLow = findMin(candles[i-8:i])
			tenkan = (highestHigh + lowestLow) / 2
			if i >= 25 {
				highestHigh = findMax(candles[i-25:i])
				lowestLow = findMin(candles[i-25:i])
				kijun = (highestHigh + lowestLow) / 2
				if i >= 77 {
					senkouA = (result[i-51].Tenkan + result[i-51].Kijun)/2
				}
			}
		}
		result = append(result, IchimokuCloud{Tenkan: tenkan, Kijun:kijun, Chikou:chikou, SenkouA:senkouA})
	}
	return result
}

func findMax(items []bittrex.Candle) float64 {
	max := items[0].High
	for _, item := range items {
		max = math.Max(max, item.High)
	}
	return max
}

func findMin(items []bittrex.Candle) float64 {
	min := items[0].Low
	for _, item := range items {
		min = math.Min(min, item.Low)
	}
	return min
}