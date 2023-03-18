package main

import (
	"container/list"
	"fmt"
	"sync"
)

type NPC struct {
	name       [64]byte
	nameLen    uint16
	blood      uint16
	properties uint32
	x, y       float64
}

func SpawnNPC(name string, x, y float64) *NPC {
	var npc = newNPC()
	npc.nameLen = uint16(copy(npc.name[:], name))
	npc.x = x
	npc.y = y
	return npc
}

var npcPool = struct {
	sync.Mutex
	*list.List
}{
	List: list.New(),
}

func newNPC() *NPC {
	npcPool.Lock()
	defer npcPool.Unlock()
	if npcPool.Len() == 0 {
		return &NPC{}
	}
	// fmt.Println("len of pool > 0")
	return npcPool.Remove(npcPool.Front()).(*NPC)
}

func releaseNPC(npc *NPC) {
	npcPool.Lock()
	defer npcPool.Unlock()
	*npc = NPC{} // zero the released NPC
	npcPool.PushBack(npc)
}

func main() {
	var npc *NPC
	var e *list.Element
	const format = "test_%d"
	for i := 0; i < 10; i++ {
		npc = SpawnNPC(fmt.Sprintf(format, i), float64(i)+1, float64(i)+2)
		// fmt.Println(npc)
		e = npcPool.PushFront(npc)
		npcPool.InsertAfter(i, e)
		// fmt.Printf("%v\n", d.Value)

	}
	value := SpawnNPC("test_100", 11, 12)
	fmt.Println(npcPool.Len())
	// npc = SpawnNPC("test_100", 11, 12)
	// e = npcPool.PushFront(npc)
	// npcPool.InsertAfter(1, e)
	// npc = SpawnNPC("test_100", 11, 12)
	// e = npcPool.PushFront(npc)
	// npcPool.InsertAfter(2, e)

	// releaseNPC(npc)
	// npc = SpawnNPC("test_101", 14, 13)
	// SpawnNPC("test_100", 11, 12)
	// fmt.Println(npcPool.Len())
	// fmt.Println(*value)

	releaseNPC(value)
	// fmt.Println(npcPool.Len())
	// for i := 0; i < 2; i++ {
	// 	SpawnNPC("test_100", 11, 12)
	// }
	fmt.Println(npcPool.Len())
	// for i := 0; i < 9; i++ {
	// 	e = npcPool.PushBack(3)
	// 	// fmt.Printf("%v\n", d.Value)

	// }
	// fmt.Println(npcPool.Len())
	// releaseNPC(npc)

}
