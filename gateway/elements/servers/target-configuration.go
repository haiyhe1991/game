package servers

import (
	"github.com/yamakiller/magicNet/st/table"
)

const (
	constTargetMax = 2048
)

//NewTargetGroup : Connection group information
func NewTargetGroup() *TargetGroup {
	return &TargetGroup{HashTable: table.HashTable{Mask: 0xFFFFFFFF, Max: constTargetMax, Comp: targetConnComparator}}
}

func targetConnComparator(a, b interface{}) int {
	c := a.(*TargetConnection)
	if c.VirtaulID == b.(uint32) {
		return 0
	}
	return 1
}

//TargetGroup connection group
type TargetGroup struct {
	table.HashTable
}

//Push Increase connection target
func (group *TargetGroup) Push(t *TargetConnection) error {
	key, err := group.HashTable.Push(t)
	if err != nil {
		return err
	}

	t.VirtaulID = key
	return nil
}

//Get Returns Target Connection Object in the group
func (group *TargetGroup) Get(virtaulID uint32) *TargetConnection {
	v := group.HashTable.Get(virtaulID)
	if v == nil {
		return nil
	}

	return v.(*TargetConnection)
}
