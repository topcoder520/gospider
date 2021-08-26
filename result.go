package gospider

type Result struct {
	TargetRequests []Request
	TargetItems    map[string]interface{}
}

func (hdl *Result) AddItem(key string, val interface{}) {
	hdl.TargetItems[key] = val
}

func (hdl *Result) AddTargetUrl(target string) {
	req := NewRequest()
	req.Url = target
	req.Skip = false
	hdl.TargetRequests = append(hdl.TargetRequests, req)
}

func (hdl *Result) AddTargetRequest(target Request) {
	hdl.TargetRequests = append(hdl.TargetRequests, target)
}
