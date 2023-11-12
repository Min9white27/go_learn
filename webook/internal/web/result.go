package web

// Result 501001 => 代表验证码
type Result struct {
	// 业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
