// Requirements:
// 1) Read a series of JSON objects from stdin

// Example input format:
// BEGIN
// {"id":1,"market":1,"price":1.8126489025824468,"volume":2529.667281360351,"is_buy":true}
// {"id":2,"market":2,"price":2.1663930707558356,"volume":3370.4751940246724,"is_buy":false}
// {"id":3,"market":11,"price":11.812182638400632,"volume":1644.641438186002,"is_buy":true}
// END
// Trade Count:  10

// 2) Send to stdout the following metrics for each market, also in JSON.
// One resulting object per market.

// Example output:
// {
//    "market":5775,
//	  "total_volume":1234567.89,
//	  "mean_price":23.33,
//	  "mean_volume":6144.299,
//	  "volume_weighted_average_price":5234.2,
//	  "percentage_buy":0.50
// }

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
	SumSpent     float32
	SumOfPrices  float32
	SumOfVolumes float32
	NumTrades    int
	NumBuys      int
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
	//log.Printf("start time: %s\n\n", startTime)

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
	mt := map[int]MarketTotals{}

	var in InputData
	var id int
	var data MarketData
	var totals MarketTotals

	// continuously compute and keep track of totals as new trades come in
	for scanner.Scan() {
		inputStr = scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		if inputStr == "END" {
			break
		}

		json.Unmarshal([]byte(inputStr), &in)

		// todo - handle non-existent map key
		id = in.MarketID
		data = MarketData{
			SumOfPrices:  md[id].SumOfPrices + in.Price,
			SumSpent:     md[id].SumSpent + in.Price*in.Volume,
			SumOfVolumes: md[id].SumOfVolumes + in.Volume,
			NumTrades:    md[id].NumTrades + 1,
			NumBuys:      md[id].NumBuys + boolToInt(in.IsBuy),
		}

		totals = MarketTotals{
			TotalVolume:            mt[id].TotalVolume + in.Volume,
			MeanPrice:              data.SumOfPrices / float32(data.NumTrades),
			MeanVolume:             data.SumOfVolumes / float32(data.NumTrades),
			VolumeWeightedAvgPrice: calcVWAP(data.SumSpent, mt[id].TotalVolume+in.Volume),
			PercentageBuy:          calcPercentage(data.NumBuys, data.NumTrades),
		}

		md[id] = data
		mt[id] = totals
	}

	totalTrades := in.ID // in.ID always increments by one for each trade

	// print all market totals as json
	for _, item := range mt {
		jsonMT, _ := json.Marshal(item) // todo - handle error
		fmt.Println(string(jsonMT))
	}

	// for debug, print first market, including source data
	resMT, _ := json.Marshal(mt[1])
	resMD, _ := json.Marshal(md[1])
	fmt.Println(string(resMT))
	fmt.Println(string(resMD))

	log.Printf("trade count: %d", totalTrades) // todo - off by one
	log.Printf("market count: %d", len(mt))
	log.Printf("\n\nduration: %s", time.Since(startTime))
}
