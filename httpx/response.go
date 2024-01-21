package httpx

func Success(data any) Result {
	// 数据脱敏
	if v, ok := data.(Desensitization); ok {
		v.Desensitization()
	}

	return Result{
		Data: data,
		Code: 10000,
		Msg:  "OK",
	}
}

func Failed(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
	}
}
