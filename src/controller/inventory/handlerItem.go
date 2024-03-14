package inventory

import "time"

type handlerItem struct {
	list      *[]InventorySell
	index     int
	lastIndex int
}

func NewHandlerItem(list *[]InventorySell) *handlerItem {
	return &handlerItem{list: list, index: 0, lastIndex: len(*list) - 1}
}
func (hi *handlerItem) IsEmpty() bool          { return hi.lastIndex < 0 }
func (hi *handlerItem) IsExist() bool          { return hi.index <= hi.lastIndex }
func (hi *handlerItem) Next()                  { hi.index++ }
func (hi *handlerItem) GetItem() InventorySell { return (*hi.list)[hi.index] }
func (hi *handlerItem) GetDate() time.Time {
	date, _ := time.Parse("02/01/2006", (*hi.list)[hi.index].Date)
	return date
}
