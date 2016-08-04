package timeoutmap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

func BenchmarkTimeoutMapTest(b *testing.B) {
	tt := NewMap(1)
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
	tt := NewMap(1)
	for i := 1; i < 10000; i++ {
		tt.Insert(strconv.Itoa(i))
		//tt.Erase(strconv.Itoa(i))
		if i%2 == 0 {
			tt.Erase(strconv.Itoa(i))
		}
	}
	fmt.Println(len(tt.tmap))
}
