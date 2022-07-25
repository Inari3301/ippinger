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
	pingList, err := getPingListFromCsv(csvDoc)
	if err != nil {
		return model.PingByCsvResult{}, err
	}

	if pingList == nil {
		return model.PingByCsvResult{}, nil
	}

	c := make(chan model.PingResult, len(pingList))
	defer close(c)
	for _, pingEntity := range pingList {
		go func(ip string, timeout int64) {
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
		}(pingEntity.ip, pingEntity.timeout)
	}

	res := model.PingByCsvResult{}
	res.PingResults = []model.PingResult{}
	for i := 0; i < len(pingList); i++ {
		pgRes := <-c
		res.PingResults = append(res.PingResults, pgRes)
	}

	return res, nil
}

func getPingListFromCsv(csvDoc []byte) ([]struct {
	ip      string
	timeout int64
}, error) {
	buf := bytes.NewBuffer(csvDoc)
	csvReader := csv.NewReader(buf)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	if records == nil {
		return nil, nil
	}

	if len(records) > 2 {
		return nil, fmt.Errorf("invalid csv")
	}

	if records[0][0] == image[0] && records[0][1] == image[1] {
		if len(records[0]) < 2 {
			return nil, fmt.Errorf("invalid csv")
		}
		records = records[1:]
	}

	out := make([]struct {
		ip      string
		timeout int64
	}, len(records))

	for i := 0; i < len(out); i++ {
		out[i].ip = records[i][0]
		out[i].timeout, err = strconv.ParseInt(records[i][1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
