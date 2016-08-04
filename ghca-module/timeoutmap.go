package timeoutmap

import (
	"fmt"
	"time"
)

type DataValue struct {
	Timestamp int64
	Ch        chan bool
}

var tmap map[string]DataValue

func New() {
	tmap = make(map[string]DataValue)
}
func Insert(k string) bool {
	var datavalue DataValue
	timenow := time.Now().Unix()
	_, ok := tmap[k]
	if ok {
		return false
	} else {
		datavalue.Timestamp = timenow
		datavalue.Ch = make(chan bool)
		tmap[k] = datavalue
		go timeoutscan(k, datavalue.Ch)
		//tmap[k] = timenow
		return true
	}
}
func Erase(k string) bool {
	data, ok := tmap[k]
	if ok {
		data.Ch <- true
		delete(tmap, k)
		return true
	} else {
		return false
	}
}

func timeoutscan(k string, ch chan bool) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1e9) // 等待1秒钟
		timeout <- true
	}()
	select {
	case <-ch:
		return
	case <-timeout:
		// 一直没有从ch中读取到数据，但从timeout中读取到了数据
		Erase(k)
		fmt.Println("timeout:", k)
	}
}
