package model

type Message struct {
	ID       int
	UserID   int
	ObjectID int
	IsGroup  int
	Type     int
	Content  string
	Image    string
	URL      string
	Number   int
}
