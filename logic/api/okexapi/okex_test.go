package okexapi

import (
	"testing"
	"xframe/log"
)

func Test_Okex(t *testing.T) {
	log.InitLogger("", "", "", 0, "DEBUG", "stdout")
	okexApi := NewOkexApi()
	err := okexApi.Start("btc", "orderbook")
	if err != nil {
		t.Error(err)
	}
	t.Log("complete")
}
