package gospider

import (
	"context"
)

type Downloader interface {
	Download(req *Request, ctx context.Context) (resp *Response, err error)
}
