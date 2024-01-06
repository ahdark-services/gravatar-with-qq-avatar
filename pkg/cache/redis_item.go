package cache

import (
	"encoding/gob"
	"github.com/AH-dark/gravatar-with-qq-support/pkg/utils"
)

type redisItem struct {
	Value interface{}
}

func init() {
	gob.Register(&redisItem{})
}

func encodeRedisItem(item redisItem) ([]byte, error) {
	return utils.GobEncode(item)
}

func decodeRedisItem(b []byte) (item redisItem, err error) {
	err = utils.GobDecode(b, &item)
	return
}
