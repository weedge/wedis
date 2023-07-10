package config

type RespCmdServiceOptins struct {
	Addr                  string `mapstructure:"addr"`
	AuthPassword          string `mapstructure:"authPassword"`
	ConnKeepaliveInterval int    `mapstructure:"connKeepaliveInterval"`
}

func DefaultRespCmdServiceOptins() *RespCmdServiceOptins {
	return &RespCmdServiceOptins{
		Addr: "127.0.0.1:6666",
		//defualt 0 disable and not check
		ConnKeepaliveInterval: 0,
	}
}
