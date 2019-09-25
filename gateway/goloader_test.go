package gateway

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/actor"
)

//TestLoader Load balancer test
func TestLoader(t *testing.T) {
	targetV := make(map[uint32]int)

	tlset := servers.NewLoadSet()
	tlset.Add("sign/in", servers.NewLoader(200))
	tlset.Add("wold/logic", servers.NewLoader(20))

	for j := 1; j < 10; j++ {
		tlset.Get("sign/in").AddTarget("sign/in#"+strconv.Itoa(j),
			&servers.TargeObject{ID: uint32(j),
				Target: &actor.PID{}})
	}

	for i := 0; i < 10000; i++ {
		t := tlset.Get("sign/in").GetTarget(strconv.Itoa(rand.Intn(10000)))
		n, ok := targetV[t.ID]
		if !ok {
			targetV[t.ID] = 1
		} else {
			targetV[t.ID] = n + 1
		}
	}

	for k, v := range targetV {
		fmt.Printf("ID:%d,Num:%d\n", k, v)
	}
}
