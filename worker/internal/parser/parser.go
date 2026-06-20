package parser

import (
	"context"
)

type ParsedItem struct {
	Title string
	Price float64
	ImageURL string
	ProductURL string
}

type ItemParser interface {
	Parse(ctx context.Context, url string) ([]ParsedItem, error)
}
