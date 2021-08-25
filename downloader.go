package gospider

import (
	"context"
)

type Downloader interface {
	Download(url string, ctx context.Context) (html string, err error)
	SetHeader(header map[string][]string)
}
