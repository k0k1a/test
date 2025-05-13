package src

import (
	"fmt"
	"testing"
)

func TestCacheManager(t *testing.T) {

	go NewTLSServer()
	//NewTLSServer()

	ll, err := BcjClient([]string{"a", "b", "c"})
	if err != nil {
		fmt.Println(err)
	}
	for i, item := range ll {
		fmt.Println(ll[i], item)
	}
}
