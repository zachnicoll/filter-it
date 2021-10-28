package util

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
	Id          string `json:"id,omitempty"`
	DateCreated int64  `json:"date_created,omitempty"`
	Filter      int    `json:"filter"`
	Progress    int    `json:"progress,omitempty"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Image       string `json:"image"`
}

type QueueResponse struct {
	DocumentID string `json:"id"`
}
