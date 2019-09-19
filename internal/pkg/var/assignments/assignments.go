package assigns

var (
	//                 invkID
	assignments = make(map[string]assign)
)

type assign struct {
	id       string           // ID of the assignment
	assignee string           // username of the assignee
	address  string           // IP of the assignee
	holder   chan interface{} // channel to pass result
}
