package ghcamodule

import (
	"fmt"
	"strconv"
	//"fmt"
	"testing"
	"time"
)

/*func TestInsert(t *testing.T) {
	var exit chan bool
	num := 1000000
	NewSessionM(num, 20)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println(GetMapSize())
		}
	}()
	var se Session
	se.Id = "1234"
	se.Uname = "ghca"
	se.Baseip = "1.1.1.1"
	se.Timet = time.Now().Unix()
	fmt.Println("start")
	for i := 0; i < num; i++ {
		se.Id = strconv.Itoa(i)
		InsertSession(se)
		if i > num/2 {
			se.Uname = "ghca" + se.Id
			//fmt.Println(se.Uname)
		}
		if i == num/2 {
			//	DeleteByUserName("ghca")
		}
	}
	fmt.Println(GetMapSize())
	<-exit
}
*/
func BenchmarkSessionM(b *testing.B) {
	num := 1000000
	NewSessionM(num, 20)
	var se Session
	se.Id = "1234"
	se.Uname = "ghca"
	se.Baseip = "1.1.1.1"
	se.Timet = time.Now().Unix()
	for i := 0; i < b.N; i++ {
		se.Id = strconv.Itoa(i)
		InsertSession(se)
	}
	fmt.Println(GetMapSize())
}
