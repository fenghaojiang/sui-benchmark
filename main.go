package main

import (
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type EventEnvelope struct {
	Timestamp     uint64 `gorm:"column:timestamp"`
	SeqNum        uint64 `gorm:"column:seq_num"`
	EventNum      uint64 `gorm:"column:event_num"`
	TxDigest      string `gorm:"column:tx_digest"`
	EventType     uint64 `gorm:"column:event_type"`
	PackageID     string `gorm:"column:package_id"`
	ModuleName    string `gorm:"column:module_name"`
	Function      string `gorm:"column:function"`
	ObjectType    string `gorm:"column:object_type"`
	ObjectID      string `gorm:"column:object_id"`
	Fields        string `gorm:"column:fields"`
	MoveEventName string `gorm:"column:move_event_name"`
	Contents      []byte `gorm:"column:contents"`
	Sender        []byte `gorm:"column:sender"`
	Recipient     string `gorm:"column:recipient"`
}

func main() {
	dataDir := "/datadir/events.db"
	database, err := gorm.Open(sqlite.Open(dataDir), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	go func() {
		for {
			select {
			case <-time.After(1 * time.Second):
				//new object events
				var events []EventEnvelope
				err := database.Table("events").Select("*").Where("event_type = ?", 8).Order("seq desc ").Offset(0).Limit(10).Scan(&events).Error
				if err != nil {
					fmt.Println(err.Error())
				} else {
					for i := range events {
						fmt.Println(fmt.Sprintf("new object events %+v", events[i]))
					}
				}
			}
		}
	}()
	http.ListenAndServe(":9999", nil)
}
