package ghcamodule

import (
	//"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var tmap map[string]*TT
var timespan time.Duration

type TT struct {
	t    *time.Timer
	key  string
	flag bool
}

func (tt *TT) f() {
	mutex.Lock()
	if tt.flag {
		delete(tmap, tt.key)
	}
	mutex.Unlock()
}

func NewTT(s string) {
	tt := new(TT)
	tt.key = s
	tt.flag = true
	mutex.Lock()
	tmap[s] = tt
	mutex.Unlock()
	tt.t = time.AfterFunc(timespan*time.Second, tt.f)
}
func Erase(key string) {
	mutex.Lock()
	value, ok := tmap[key]
	if ok {
		value.t.Stop()
		value.flag = false
		//fmt.Println("erase", value.key)
		delete(tmap, key)
	}
	mutex.Unlock()
}
func Find(key string) (*TT, bool) {
	value, ok := tmap[key]
	return value, ok
}
func Init(t int64) {
	tmap = make(map[string]*TT)
	timespan = time.Duration(t)
}

type DataValue struct {
	Timestamp int64
	Ch        chan bool
}
type TTMap struct {
	tmap     map[string]DataValue
	lock     sync.Mutex
	timeout  int
	Datachan chan string
}

func NewMap(t int, cache int) *TTMap {
	mm := new(TTMap)
	mm.tmap = make(map[string]DataValue)
	mm.timeout = t
	mm.Datachan = make(chan string, cache)
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
		tt.Datachan <- k
		//fmt.Println("timeout:", k)
	}
}
