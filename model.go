package main

type textData struct {
	Description string `json:"description"`
	Content     string `json:"content"`
}

type textID struct {
	ID string `json:"id"`
}

type textEntry struct {
	textID
	textData
}
