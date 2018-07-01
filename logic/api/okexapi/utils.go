package okexapi

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func (this *OkexApi) getKey(contract string, table string, ttype string) string {
	return fmt.Sprintf("%s|%s|%s", contract, table, ttype)
}

func (this *OkexApi) checkExt(key string) (ok bool) {
	_, ok = this.upstreamConns.Load(key)
	return
}

func (this *OkexApi) _getLoginTable(channel string) string {
	fields := strings.Split(channel, "_")
	return fields[len(fields)-1]
}

func (this *OkexApi) Sign(params map[string]string, api_secret string) string {
	var (
		keyLst     = make([]string, 0)
		sortParams string
	)
	for key, _ := range params {
		keyLst = append(keyLst, key)
	}
	sort.Strings(keyLst)
	for _, key := range keyLst {
		sortParams += key + "=" + params[key] + "&"
	}
	sortParams += "secret_key=" + api_secret
	h := md5.New()
	io.WriteString(h, sortParams)
	sign := strings.ToUpper(fmt.Sprintf("%x", h.Sum(nil)))
	return sign
}

func (this *OkexApi) generateCommonQuoteData(msg string, dataRes []DataResponse, depth int) ([]byte, error) {
	var (
		dataPush DataPush
	)
	dataPush.Msg = msg
	dataPush.Market = this.Market
	dataPush.Table = this.Table
	dataPush.Contract = this.Symbol
	dataPush.Type = this.Type
	dataPush.Optional = Option{
		Period: "",
		Depth:  depth,
	}
	dataPush.Data = make([]QuotePush, 0)
	for _, data := range dataRes {
		var tmpQuotePush QuotePush
		tmpQuotePush.Ts = data.Data.Ts
		tmpQuotePush.Asks = data.Data.Asks
		tmpQuotePush.Bids = data.Data.Bids
		dataPush.Data = append(dataPush.Data, tmpQuotePush)
	}
	buf, err := json.Marshal(dataPush)
	return buf, err
}

func (this *OkexApi) generateCommonData(msg string, dataCommonRes []DataCommonRes, pushTable string) ([]byte, error) {
	var (
		dataCommonPush DataCommonPush
	)
	dataCommonPush.Msg = msg
	dataCommonPush.Market = this.Market
	dataCommonPush.Table = pushTable
	dataCommonPush.Contract = this.Symbol
	dataCommonPush.Type = this.Type
	dataCommonPush.Data = make([]interface{}, 0)
	//structure抽象
	for _, data := range dataCommonRes {
		dataCommonPush.Data = append(dataCommonPush.Data, data.Data)
	}
	buf, err := json.Marshal(dataCommonPush)
	return buf, err

}
