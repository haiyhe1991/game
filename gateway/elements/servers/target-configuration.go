package servers

import (
	"github.com/yamakiller/magicNet/st/table"
)

const (
	constTargetMax = 2048
)

//NewTargetSet : Connection group information
func NewTargetSet() *TargetSet {
	r := &TargetSet{HashTable: table.HashTable{Mask: 0xFFFFFFFF, Max: constTargetMax, Comp: targetConnComparator}}
	r.Init()
	return r
}

func targetConnComparator(a, b interface{}) int {
	c := a.(*TargetConnection)
	if c.virtualID == b.(uint32) {
		return 0
	}
	return 1
}

//TargetSet connection set
type TargetSet struct {
	table.HashTable
}

//Push Increase connection target
func (tset *TargetSet) Push(t *TargetConnection) error {
	key, err := tset.HashTable.Push(t)
	if err != nil {
		return err
	}

	t.virtualID = key
	return nil
}

//Get Returns Target Connection Object in the Set
func (tset *TargetSet) Get(virtaulID uint32) *TargetConnection {
	v := tset.HashTable.Get(virtaulID)
	if v == nil {
		return nil
	}

	return v.(*TargetConnection)
}
