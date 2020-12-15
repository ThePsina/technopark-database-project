package entity

//easyjson:json
type Info struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread,omitempty"`
}
