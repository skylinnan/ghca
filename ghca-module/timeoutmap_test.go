package ghcamodule

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

var exit chan bool

func BenchmarkTimeoutMapTest(b *testing.B) {
	tt := NewMap(1, 100000)
	for i := 1; i < b.N; i++ {
		v := strconv.Itoa(i)
		tt.Insert(v)
		if i%2 == 0 {
			tt.Erase(v)
		}
	}
	//fmt.Println(len(tt.tmap))
}
func BenchmarkTimeoutGoTest(b *testing.B) {
	tt := make(map[string]string)
	var lock sync.Mutex
	for i := 1; i < b.N; i++ {
		v := strconv.Itoa(i)
		lock.Lock()
		_, ok := tt[v]
		if !ok {
			tt[v] = v
		}
		lock.Unlock()
		lock.Lock()
		_, ok = tt[v]
		if ok {
			delete(tt, v)
		}
		lock.Unlock()
	}
}

func TestMap(t *testing.T) {
	/*tt := NewMap(1, 1000)
	for i := 1; i < 1000; i++ {
		tt.Insert(strconv.Itoa(i))
		//tt.Erase(strconv.Itoa(i))
		if i%2 == 0 {
			tt.Erase(strconv.Itoa(i))
		}
	}*/

	Init(5)
	NewTT("a")
	time.Sleep(2 * time.Second)
	//Erase("a")
	_, ok := Find("a")
	if ok {
		t.Log("pass")
	} else {
		t.Error("test failed")
	}
	time.Sleep(4 * time.Second)
	_, ok = Find("a")
	if ok {
		t.Error("test timeout failed")
	} else {
		t.Log("test timeout success.")
	}
	/*for range tt.Datachan {
		data := <-tt.Datachan
		fmt.Println(data)
	}*/
}
