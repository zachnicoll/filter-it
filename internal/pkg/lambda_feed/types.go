package lambda_feed

type FeedRequestBody struct {
	Filters []string `json:"filters"`
}

const (
	GREYSCALE int = 0
	SEPIA         = 1
	INVERT        = 2
)

const (
	READY      int = 0
	PROCESSING     = 1
	DONE           = 2
	FAILED         = 3
)

type ImageDocument struct {
	Id          string `json:"id"`
	DateCreated int    `json:"date_created"`
	Filters     []int  `json:"filters"`
	Progress    int    `json:"progress"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Image       string `json:"image"`
}
