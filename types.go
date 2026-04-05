package artwork

type LensResponse struct {
	AboutThisImage AboutThisImage `json:"about_this_image"`
}

type AboutThisImage struct {
	Sections []Section `json:"sections"`
}

type Section struct {
	Title       string       `json:"title"`
	PageResults []PageResult `json:"page_results"`
}

type PageResult struct {
	Title   string `json:"title"`
	Source  string `json:"source"`
	Snippet string `json:"snippet"`
}

type ArtworkResponse struct {
	Artist             string `json:"artist"`
	ArtworkTitle       string `json:"artworkTitle"`
	DateOfCreation     string `json:"dateOfCreation"`
	ArtworkDescription string `json:"artworkDescription"`
	WikipediaLink      string `json:"wikipediaLink"`
}