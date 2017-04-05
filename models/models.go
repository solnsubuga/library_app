package models

//Page ..
type Page struct {
	Books  []Book
	Filter string
	User   string
}

//ClassifySearchResponse ...
type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

//SearchResult ..
type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

//ClassifyBookResponse ...
type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

//Book is a what is stored in a library
type Book struct {
	Pk             int64  `db:"pk"`
	Title          string `db:"title"`
	Author         string `db:"author"`
	Classification string `db:"classification"`
	ID             string `db:"id"`
}

//User who uses a library
type User struct {
	ID       string `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
