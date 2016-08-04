package timeoutmap

import (
	//"fmt"
	"sync"
	"time"
)

type DataValue struct {
	Timestamp int64
	Ch        chan bool
}
type TTMap struct {
	tmap    map[string]DataValue
	lock    sync.Mutex
	timeout int
}

func NewMap(t int) *TTMap {
	mm := new(TTMap)
	mm.tmap = make(map[string]DataValue)
	mm.timeout = t
	return mm
}
func (tt *TTMap) Insert(k string) bool {
	var datavalue DataValue
	tt.lock.Lock()
	defer tt.lock.Unlock()
	_, ok := tt.tmap[k]
	if ok {
		return false
	} else {
		timenow := time.Now().Unix()
		datavalue.Timestamp = timenow
		datavalue.Ch = make(chan bool, 1)
		tt.tmap[k] = datavalue
		//fmt.Println("insert")
		go tt.timeoutscan(k, datavalue.Ch)
		//tmap[k] = timenow
		return true
	}
}
func (tt *TTMap) Erase(k string) bool {
	tt.lock.Lock()
	defer tt.lock.Unlock()
	data, ok := tt.tmap[k]
	if ok {
		data.Ch <- true
		delete(tt.tmap, k)
		//fmt.Println("erase")
		return true
	} else {
		return false
	}
}

func (tt *TTMap) timeoutscan(k string, ch chan bool) {
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(tt.timeout) * time.Second) // 等待1秒钟
		timeout <- true
		//fmt.Println("timeout")
	}()
	select {
	case <-ch:
		//fmt.Println("erase ch")
		return
	case <-timeout:
		// 一直没有从ch中读取到数据，但从timeout中读取到了数据
		//fmt.Println("start to timeout erase")
		tt.Erase(k)
		//fmt.Println("timeout:", k)
	}
}
