package srv

func Ping(h *ConnClient, params [][]byte) (res interface{}, err error) {
	res = "PONG"
	return
}
