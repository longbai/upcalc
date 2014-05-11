package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"time"
)

type Log struct {
	Time     int64 `json:"stamp"`
	Duration int64 `json:"dur"`
	Length   int64 `json:"len"`
}

type TransferRec struct {
	Time     int64
	Transfer int64
}

var zone *time.Location = time.FixedZone("CST", 8*3600)
var start time.Time = time.Date(2014, 4, 1, 0, 0, 0, 0, zone)
var startUnix int64 = start.Unix()
var end time.Time = time.Date(2014, 5, 1, 0, 0, 0, 1, zone)
var endUnix int64 = end.Unix()

const interval = 5 * 60

const pointNumber = 288 * 30

func (t *TransferRec) timeString() string {
	_t := time.Unix(t.Time, 0)
	return _t.In(zone).String()
}

func (t *TransferRec) Bandwidth() int64 {
	return t.Transfer * 8 / interval
}

func (t *TransferRec) print() {
	fmt.Println(t.Time, t.timeString(), t.Bandwidth())
}

type TransferArray []TransferRec

func (p TransferArray) Len() int           { return len(p) }
func (p TransferArray) Less(i, j int) bool { return p[i].Transfer < p[j].Transfer }
func (p TransferArray) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p TransferArray) Sort()              { sort.Sort(p) }

func valid(utc int64) bool {
	return utc >= startUnix && utc < endUnix
}

func initTransferPoint() []TransferRec {
	recs := make([]TransferRec, pointNumber)
	fromUtc := startUnix + interval
	for i := 0; i < pointNumber; i++ {
		recs[i].Time = fromUtc + int64(i*interval)
	}
	return recs
}

func main() {
	input := flag.String("i", "", "input Json")
	flag.Parse()
	recs := initTransferPoint()
	_data, err := ioutil.ReadFile(*input)
	if err != nil {
		flag.PrintDefaults()
		fmt.Println(err)
		return
	}
	var logs []Log
	err = json.Unmarshal(_data, &logs)
	if err != nil {
		fmt.Println(err)
		return
	}
	var total int64
	for _, v := range logs {
		t := v.Time / 1e7
		if !valid(t) {
			fmt.Println("invalid", t)
			continue
		}
		index := (t - startUnix) / interval
		recs[index].Transfer += v.Length
		total += v.Length
	}
	TransferArray(recs).Sort()
	pos := int(float32(pointNumber) * 0.95)
	fmt.Print(pos, " ")
	recs[pos].print()

	fmt.Print(pos+1, " ")
	recs[pos+1].print()
	fmt.Print(pointNumber, " ")
	recs[pointNumber-1].print()
	fmt.Println("total upload bytes", total)
}
