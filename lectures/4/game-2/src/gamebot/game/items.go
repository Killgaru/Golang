// items
package game

// import "fmt"

var (
	MapOfItems map[string]*Item // MOI
)

type Item struct {
	Name        string
	iType       string
	description map[*Room]string
}

func createItem(n, it string, des map[*Room]string) {
	item := Item{
		Name:        n,
		iType:       it,
		description: des,
	}
	item.addToMOI()
	item.SpawnItem()
	return
}

func (it *Item) SpawnItem() {
	for r := range it.description {
		r.Items = append(r.Items, it)
		r.refreshRoomDes()
	}
	return
}

func (it *Item) addToMOI() {
	if _, ok := MapOfItems[it.Name]; !ok {
		MapOfItems[it.Name] = it
	}
	return
}

func InitItems() {
	MapOfItems = map[string]*Item{}
	createItem("ключи", "key", map[*Room]string{
		World["домой"]["домой"]["комната"]: "на столе: ",
	},
	)
	createItem("конспекты", "object", map[*Room]string{
		World["домой"]["домой"]["комната"]: "на столе: ",
	},
	)
	createItem("рюкзак", "clothes", map[*Room]string{
		World["домой"]["домой"]["комната"]: "на стуле - ",
	},
	)
	return
}
