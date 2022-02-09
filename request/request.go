package request

type Request struct {
	Auth   Auth     `json:"auth"`
	Method string   `json:"method"`
	Header []Header `json:"header"`
	Body   Body     `json:"body"`
	URL    URL      `json:"url"`
}
