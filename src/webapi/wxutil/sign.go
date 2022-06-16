package wxutil

import (
	"sort"
	"strings"

	"github.com/hololee2cn/wxpub/v1/src/utils"
)

func CalcSign(params ...string) (sign string) {
	sort.Strings(params)
	var b strings.Builder
	for _, v := range params {
		b.WriteString(v)
	}
	return utils.Sha1(b.String())
}

func CheckSign(sign1, sign2 string) bool {
	return len(sign1) > 0 && sign1 == sign2
}
