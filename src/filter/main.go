package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Log struct {
	Time     int64 `json:"stamp"`
	Duration int64 `json:"dur"`
	Length   int64 `json:"len"`
}

var invalid error = errors.New("invalid")

type body struct {
	Length string `json:"Content-Length"`
}

func parseLog(line string) (log Log, err error) {
	rec := strings.Split(line, "\t")
	status := rec[7]
	if status != "200" {
		err = invalid
		return
	}
	log.Time, err = strconv.ParseInt(rec[2], 10, 64)
	if err != nil {
		return
	}
	var b body
	err = json.Unmarshal([]byte(rec[5]), &b)
	if err != nil {
		return
	}
	log.Length, err = strconv.ParseInt(b.Length, 10, 64)
	if err != nil {
		return
	}
	if len(rec) < 12 {
		fmt.Println(line)
		err = errors.New("short")
		return
	}
	end := rec[11]
	vars := strings.Split(end, " ")
	log.Duration, err = strconv.ParseInt(vars[0], 10, 64)
	return
}

func readFile(logs []Log, path string) (ret []Log, err error) {
	fmt.Println("parse file", path)
	file, err := os.OpenFile(path, os.O_RDONLY, 0700)
	if err != nil {
		return
	}
	ret = logs
	scan := bufio.NewScanner(file)
	count := 0
	for scan.Scan() {
		line := scan.Text()
		log, err1 := parseLog(scan.Text())
		count++
		if err1 == invalid {
			continue
		}
		if err1 != nil {
			err = err1
			fmt.Println("line", count, line, err)
			return
		}
		ret = append(ret, log)
	}
	return
}

func main() {
	out := flag.String("o", "", "output log")
	flag.Parse()

	if *out == "" {
		flag.PrintDefaults()
		fmt.Println("invalid args")
		return
	}

	var logs []Log
	var err error
	for _, v := range flag.Args() {
		logs, err = readFile(logs, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	d, _ := json.MarshalIndent(logs, "", "")
	ioutil.WriteFile(*out, d, 0666)

}
