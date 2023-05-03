package dtm

type DtmConf struct {
	HTTPServer string `json:",optional"` // HTTPServer
	GRPCServer string `json:",optional"` // GRPCServer
}
