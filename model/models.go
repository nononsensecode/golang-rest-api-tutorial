package model

// Article type exposes our article object
type Article struct {
	ID      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

//Articles slice of articles
var Articles []Article
