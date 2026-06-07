package parser

type ParsedItem struct {
	Title string
	Price float64
	ImageUrl string
}

type ItemParser interface {
	Parse(url string) (ParsedItem, error)
}