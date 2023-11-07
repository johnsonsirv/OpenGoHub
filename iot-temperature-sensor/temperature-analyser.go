package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	"github.com/wcharczuk/go-chart"
)

type IOTData struct {
	Name string `json:"name"`
	MinTemp float32 `json:"tempMin"`
	MaxTemp float32 `json:"tempMax"`
	Interval int `json:"interval"`
	Values []Value `json:"values"`
}

type Value struct {
	Message int `json:"messageId"`
	Temperature float32 `json:"temperature"`
	EnqueuedTime string `json:"enqueuedTime"`
}

type TemperatureReading struct {
	Hour int
	Normal float32
	OutOfRange float32
}

func main()  {
	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer jsonFile.Close()

	byteVal, _ := io.ReadAll(jsonFile)
	var data IOTData

	json.Unmarshal(byteVal, &data)

	tempMap := make(map[int][]float32)
	fmt.Println("Total sensor readings", len(data.Values))

	for _, e := range data.Values {
		t, err := time.Parse("2006-01-02 15:04:05", e.EnqueuedTime)

		if err != nil {
			log.Fatal(err.Error())
		}

		h := t.Hour()
		tempMap[h] = append(tempMap[h], e.Temperature)
	}

	var readings []TemperatureReading
	var normalReading, outOfRange float32

	for h, v := range tempMap {
		normalReading, outOfRange = 0.0, 0.0
		for _, r := range v {
			if r >= data.MinTemp  && r <= data.MaxTemp{
				normalReading++
			} else {
				outOfRange++
			}
		}

		reading := TemperatureReading{h, normalReading, outOfRange}
		readings = append(readings, reading)
	}

	sort.Slice(readings, func(i, j int) bool {
		return readings[i].Hour < readings[j].Hour
	})

	printTableOnConsole(readings)
	printChart(readings)
}

func printTableOnConsole(r []TemperatureReading)  {
	fmt.Printf("Hour\tTotal\tNormal\tOut of Range\tPercent\n")

	for _, val := range r {
		total := val.Normal + val.OutOfRange
		percent := val.OutOfRange / total * 100
		fmt.Printf("%v\t%v\t%v\t%5v\t\t%5.1f\n", val.Hour, total, val.Normal, val.OutOfRange, percent )
	}
}

func printChart(r []TemperatureReading)  {
	var bars []chart.StackedBar
	
	for _, v := range r {
		label := fmt.Sprintf("Hour %d", v.Hour)
		bar := chart.StackedBar{
			Name: label,
			Values: []chart.Value{
				{ Label: "Green", Value: float64(v.Normal)},
				{ Label: "Red", Value: float64(v.OutOfRange)},
			},
		}

		bars = append(bars, bar)
	}

	sbc := chart.StackedBarChart{
		Title: "IOT Sensor Bar Chart",
		Bars: bars,
		Height: 512,
		Width: 2000,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
	}

	f, _ := os.Create("iot-temperature-sensor-bar-chart.png")
	defer f.Close()

	sbc.Render(chart.PNG, f)
}