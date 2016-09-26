package ghcamodule

import (
	"sync"
	//"time"
)

type Session struct {
	Id     string
	Uname  string
	Baseip string
	Timet  int64
}

var (
	rwmutex     sync.RWMutex
	sessionM    map[string]Session
	delipchan   chan string
	delnamechan chan string
)

func NewSessionM(catch int, proc int) {
	sessionM = make(map[string]Session)
	delipchan = make(chan string, catch)
	delnamechan = make(chan string, catch)
	for i := 0; i < proc; i++ {
		go delprocess()
	}

}

func InsertSession(se Session) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	sessionM[se.Id] = se
}

func DeleteSession(key string) {
	rwmutex.Lock()
	defer rwmutex.Unlock()
	delete(sessionM, key)
	//fmt.Println(len(sessionM), key)
}

func FindSession(key string) (Session, bool) {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	v, ok := sessionM[key]
	return v, ok
}

func GetMapSize() int {
	return len(sessionM)
}
func DeleteByBaseIP(ip string) {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	for k, v := range sessionM {
		if v.Baseip == ip {
			delipchan <- k
		}
	}
}

func DeleteByUserName(name string) {
	rwmutex.RLock()
	defer rwmutex.RUnlock()
	for k, v := range sessionM {
		if v.Uname == name {
			//fmt.Println("find ", v)
			delnamechan <- k
		}
	}
}

func delprocess() {
	for {
		var key string
		select {
		case key = <-delnamechan:
		case key = <-delipchan:
		}
		//fmt.Println("start to delete ", key)
		DeleteSession(key)
	}
}
