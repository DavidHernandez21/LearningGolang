package rp

import (
	"fmt"
	"src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/gotr/resource_pool/myOwnPool/model"
)

const CR_POOL_SIZE = 300

var pool chan *model.ClientReq
var alloced, reused uint

func init() {
	pool = make(chan *model.ClientReq, CR_POOL_SIZE)
}

func Alloc() *model.ClientReq {
	select {
	case cr := <-pool:
		reused++
		return cr
	default:
		alloced++
		cr := &model.ClientReq{}
		return cr
	}
}

func Release(cr *model.ClientReq) {
	select {
	case pool <- cr:
	default:
	}
}

func Stats() {
	fmt.Printf("Total: %v, Allocated: %v, Reused: %v\n", alloced+reused, alloced, reused)
}
