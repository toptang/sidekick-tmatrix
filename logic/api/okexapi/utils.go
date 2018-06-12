package okexapi

import (
	"fmt"
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
