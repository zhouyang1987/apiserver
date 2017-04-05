package componentconfig

type ApiserverConfig struct {
	HttpAddr string
	HttpPort int
	RpcAddr  string
	RpcPort  int
}

type BuildConfig struct {
	HttpAddr string
	HttpPort int
	RpcAddr  string
	RpcPort  int
	Endpoint string
	Version  string
}
