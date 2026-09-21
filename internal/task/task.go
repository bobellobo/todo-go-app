package task

// todo item
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Group string `json:"group,omitempty"` // default is empty ""
}
