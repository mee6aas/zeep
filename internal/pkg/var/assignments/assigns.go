package assigns

var (
	//                 invkID
	assigns = make(map[string]assign)
)

type assign struct {
	id       string
	assignee string
	address  string
	holder   chan interface{}
}
