// the requests built from blocks, which should be parseable directly.
package request

// type Event struct {
// 	Listen string `json:"listen"`
// 	Script Script `json:"script"`
// }

// type Script struct {
// 	Exec []string `json:"exec"`
// 	Type string   `json:"type"`
// }

// type Basic struct {
// 	Username     string `json:"username"`
// 	Password     string `json:"password"`
// 	ShowPassword bool   `json:"showPassword"`
// }

// type Bearer struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// 	Type  string `json:"type"`
// }

// type Auth struct {
// 	Type   string   `json:"type"`
// 	Basic  []Basic  `json:"basic"` // should be implementation specific
// 	Bearer []Bearer `json:"bearer"`
// }

// type Header struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// type Urlencoded struct {
// 	Key      string `json:"key"`
// 	Value    string `json:"value"`
// 	Type     string `json:"type"` // text, number, bool, null
// 	Disabled bool   `json:"disabled,omitempty"`
// }

// type Body struct {
// 	Mode       string       `json:"mode"`
// 	Urlencoded []Urlencoded `json:"urlencoded,omitempty"`
// 	Raw        string       `json:"raw,omitempty"`
// }

// type Options struct {
// 	Raw Raw `json:"raw"`
// }

// type Query struct {
// 	Key   string `json:"key"`
// 	Value string `json:"value"`
// }

// type URL struct {
// 	Raw      string   `json:"raw"`
// 	Protocol string   `json:"protocol"`
// 	Host     []string `json:"host"`
// 	Path     []string `json:"path"`
// 	Query    []Query  `json:"query"`
// }

// type Raw struct {
// 	Language string `json:"language"`
// }

// "apikey",
// "awsv4",
// "basic",
// "bearer",
// "digest",
// "edgegrid",
// "hawk",
// "noauth",
// "oauth1",
// "oauth2",
// "ntlm"

// ref: https://www.postman.com/collections/15f01ec6b24aee48cc19
