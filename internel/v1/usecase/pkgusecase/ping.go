package pkgusecase

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/Inari3301/ippinger/internel/v1/model"
	"github.com/Inari3301/ippinger/pkg/ping"
	"strconv"
	"time"
)

var (
	image = []string{"ip", "timeout"}
)

type P struct{}

func (p *P) Ping(ip string, timeout time.Duration) (model.PingResult, error) {
	d, err := ping.Ping(ip, timeout)
	if err != nil {
		return model.PingResult{}, err
	}

	return model.PingResult{
		IP:       ip,
		Duration: d,
	}, nil
}

func (p *P) PingByCsv(csvDoc []byte) (model.PingByCsvResult, error) {
	buf := bytes.NewBuffer(csvDoc)
	csvReader := csv.NewReader(buf)
	records, err := csvReader.ReadAll()
	if err != nil {
		return model.PingByCsvResult{}, err
	}
	if len(records) == 0 {
		return model.PingByCsvResult{}, fmt.Errorf("emty csv file")
	}

	if len(records[0]) < len(image) {
		return model.PingByCsvResult{}, fmt.Errorf("bad scv file, need columns %s|%s", image[0], image[1])
	}

	if records[0][0] != image[0] || records[0][1] != image[1] {
		return model.PingByCsvResult{}, fmt.Errorf("bad scv file, need columns %s|%s", image[0], image[1])
	}

	records = records[1:]
	recordsNum := len(records[0])
	c := make(chan model.PingResult, recordsNum)
	for i := 0; i < recordsNum; i++ {
		go func(index int) {
			ip := records[0][index]
			timeout, err := strconv.ParseUint(records[1][index], 10, 64)
			if err != nil {
				c <- model.PingResult{
					IP:       ip,
					Duration: 0,
					Error:    err.Error(),
				}
				return
			}

			lag, err := ping.Ping(ip, time.Duration(timeout)*time.Second)
			if err != nil {
				c <- model.PingResult{
					IP:       ip,
					Duration: 0,
					Error:    err.Error(),
				}
				return
			}

			c <- model.PingResult{
				IP:       ip,
				Duration: lag,
				Error:    "",
			}
		}(i)
	}
	res := model.PingByCsvResult{}
	res.PingResults = make([]model.PingResult, recordsNum)
	i := 0
	for pgRes := range c {
		res.PingResults[i] = pgRes
		i++
	}

	return res, nil
}
