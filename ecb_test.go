package ecb_test

import (
	"log"
	"testing"

	"github.com/xpt-nl/ecb"
)

func TestExchangeRate(t *testing.T) {
	rate, err := ecb.EUR(ecb.USD)
	if err != nil {
		t.Fail()
	}
	log.Println(rate)
}
