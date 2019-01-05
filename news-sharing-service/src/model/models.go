package model

//NewItemSubmission ...
type NewItemSubmission struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	//SubmittedBy string `json:"submittedBy"`
}

//NewsItem ...
type NewsItem struct {
	NewsID      string `json:"newsID"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	SubmittedBy string `json:"submittedBy"`
	Upvotes     string `json:"numUpvotes"`
	Comments    string `json:"numComments"`
}

//NewsItemComments ...
type NewsItemComments struct {
	NewsID   string   `json:"newsID"`
	Comments []string `json:"comments"`
}
