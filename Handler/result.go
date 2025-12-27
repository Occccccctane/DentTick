package Handler

type Result struct {
	// 2 成功 3 重定向 4 客户端错误 5 服务端错误
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"date"`
}
