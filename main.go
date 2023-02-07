package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

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

type MarketTotals struct {
	TotalVolume            float32 `json:"total_volume"`
	MeanPrice              float32 `json:"mean_price"`
	MeanVolume             float32 `json:"mean_volume"`
	VolumeWeightedAvgPrice float32 `json:"volume_weighted_average_price"`
	PercentageBuy          float32 `json:"percentage_buy"`
}

type InputData struct {
	ID       int     `json:"id"`
	MarketID int     `json:"market"`
	Price    float32 `json:"price"`
	Volume   float32 `json:"volume"`
	IsBuy    bool    `json:"is_buy"`
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

func main() {
	startTime := time.Now()
	scanner := bufio.NewScanner(os.Stdin)
	var inputStr string

	// look for the start
	for scanner.Scan() {
		inputStr = scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
		if inputStr == "BEGIN" {
			break
		}
	}

	md := map[int]MarketData{}
	var in InputData
	var data MarketData

	// continuously compute and keep track of totals as new trades come in
	// assumptions:
	// valid json input, no zero values

	if err := scanner.Err(); err != nil { // todo - move?
		log.Fatal(err)
	}

	for scanner.Scan() {
		inputStr = scanner.Text()
		if inputStr == "END" {
			break
		}

		json.Unmarshal([]byte(inputStr), &in)

		data = md[in.MarketID] // todo - handle non-existent map key
		data.update(in)
		md[in.MarketID] = data
	}

	totalTrades := in.ID // in.ID always increments by one for each trade

	// print all market totals as json
	for _, item := range md {
		jsonMT, _ := json.Marshal(item.getMarketTotals()) // todo - handle error
		fmt.Println(string(jsonMT))
	}

	// for debug, print first market, including source data
	resMD, _ := json.Marshal(md[1])
	fmt.Println(string(resMD))

	log.Printf("trade count: %d", totalTrades) // todo - off by one
	log.Printf("market count: %d", len(md))
	log.Printf("\n\nduration: %s", time.Since(startTime))
}
