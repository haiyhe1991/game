package module

var (
	//ClientWarehouse Gateway Connection User warehouse
	clientWarehouse *InNetClientManage
)

//NewWarehouse Create a warehouse
func NewWarehouse(w *InNetClientManage) {
	clientWarehouse = w
}

//GetWarehouse result the warehouse
func GetWarehouse() *InNetClientManage {
	return clientWarehouse
}
