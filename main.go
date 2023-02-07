package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type InputData struct {
	ID       int     `json:"id"`
	MarketID int     `json:"market"`
	Price    float32 `json:"price"`
	Volume   float32 `json:"volume"`
	IsBuy    bool    `json:"is_buy"`
}

type MarketData struct {
	SumSpent               float32
	SumOfPrices            float32
	SumOfVolumes           float32
	NumTrades              int
	NumBuys                int
	TotalVolume            float32
	MeanPrice              float32
	MeanVolume             float32
	VolumeWeightedAvgPrice float32
	PercentageBuy          float32
}

type MarketTotals struct {
	TotalVolume            float32 `json:"total_volume"`
	MeanPrice              float32 `json:"mean_price"`
	MeanVolume             float32 `json:"mean_volume"`
	VolumeWeightedAvgPrice float32 `json:"volume_weighted_average_price"`
	PercentageBuy          float32 `json:"percentage_buy"`
}

func (md *MarketData) update(in InputData) {
	md.SumOfPrices += in.Price
	md.SumSpent += in.Price * in.Volume
	md.SumOfVolumes += in.Volume
	md.NumTrades++
	md.NumBuys += boolToInt(in.IsBuy)
	md.TotalVolume += in.Volume

	md.MeanPrice = md.SumOfPrices / float32(md.NumTrades) // NumTrades will always be >= 1
	md.MeanVolume = md.SumOfVolumes / float32(md.NumTrades)
	md.VolumeWeightedAvgPrice = calcVWAP(md.SumSpent, md.TotalVolume+in.Volume)
	md.PercentageBuy = calcPercentage(md.NumBuys, md.NumTrades)
}

func (md *MarketData) getMarketTotals() MarketTotals {
	return MarketTotals{
		TotalVolume:            md.TotalVolume,
		MeanPrice:              md.MeanPrice,
		MeanVolume:             md.MeanVolume,
		VolumeWeightedAvgPrice: md.VolumeWeightedAvgPrice,
		PercentageBuy:          md.PercentageBuy,
	}
}

// calculates percentage num/den
func calcPercentage(num int, den int) float32 {
	if den == 0 {
		return 0
	}
	return float32(num) / float32(den)
}

func boolToInt(in bool) int {
	if in {
		return 1
	}
	return 0
}

func calcVWAP(sumSpent float32, cummulativeVolume float32) float32 {
	if cummulativeVolume == 0 || sumSpent == 0 {
		return 0
	}
	vwap := sumSpent / cummulativeVolume
	return vwap
}

func processTrade(inputStr string, md map[int]MarketData) {
	in := InputData{}
	_ = json.Unmarshal([]byte(inputStr), &in) // todo - check error

	data := md[in.MarketID]
	data.update(in)
	md[in.MarketID] = data
}

func main() {
	startTime := time.Now()
	scanner := bufio.NewScanner(os.Stdin)
	inputStr := ""

	for scanner.Scan() {
		inputStr = scanner.Text()
		if inputStr == "BEGIN" {
			break
		}
	}

	if err := scanner.Err(); err != nil { // todo - move?
		log.Fatal(err)
	}

	md := map[int]MarketData{}
	totalTrades := 0

	for scanner.Scan() {
		inputStr = scanner.Text()
		if inputStr == "END" {
			break
		}
		processTrade(inputStr, md)
		totalTrades++
	}

	for _, item := range md {
		m, _ := json.Marshal(item.getMarketTotals()) // todo - handle error
		fmt.Println(string(m))
	}

	log.Printf("\n\ntrade count: %d", totalTrades) // todo - off by one
	log.Printf("market count: %d", len(md))
	log.Printf("total duration: %s", time.Since(startTime))
}
