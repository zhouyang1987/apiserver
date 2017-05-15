package componentconfig

type ApiserverConfig struct {
	HttpAddr string
	HttpPort int
	RpcAddr  string
	RpcPort  int
	Version  string
}

type BuildConfig struct {
	HttpAddr string
	HttpPort int
	RpcAddr  string
	RpcPort  int
	Endpoint string
	Version  string
}

type RegistryConfig struct {
	HttpAddr string
	HttpPort int
	RpcAddr  string
	RpcPort  int
	Endpoint string
	Version  string
}
