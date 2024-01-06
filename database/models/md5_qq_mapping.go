package models

import (
	"encoding/gob"
	"github.com/bytedance/sonic"
	"github.com/scylladb/gocqlx/v2/table"
	"reflect"
)

type MD5QQMapping struct {
	EmailMD5 string `db:"email_md5"`
	QQId     int64  `db:"qq_id"`
}

var MD5QQMappingTable = table.New(table.Metadata{
	Name: "md5_qq_mapping",
	Columns: []string{
		"email_md5",
		"qq_id",
	},
	PartKey: []string{
		"email_md5",
	},
	SortKey: []string{
		"qq_id",
	},
})

func init() {
	gob.Register(MD5QQMapping{})
	if err := sonic.Pretouch(reflect.TypeOf(MD5QQMapping{})); err != nil {
		panic(err)
	}
}
