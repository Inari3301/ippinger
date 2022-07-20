package pkgusecase

import (
	"fmt"
	"testing"
)

func TestPing_PingByCsv(t *testing.T) {
	csvFile :=
		"ip, timeout\n" +
			"87.240.190.67, 10"
	p := Ping{}
	r, err := p.PingByCsv([]byte(csvFile))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
}
