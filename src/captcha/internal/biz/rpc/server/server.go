package server

import (
	"net"

	pb "github.com/hololee2cn/captcha/pkg/grpcIFace"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type rpcServer struct {
	address          string
	captchaSvcServer *captchaSvcServer
	s                *grpc.Server
	opts             []grpc.ServerOption
}

func NewRpcServer(address string, captchaSvcServer *captchaSvcServer, opts ...grpc.ServerOption) *rpcServer {
	return &rpcServer{address: address, captchaSvcServer: captchaSvcServer, opts: opts}
}

func (rs *rpcServer) Start() {
	listener, err := net.Listen("tcp", rs.address)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(rs.opts...)
	pb.RegisterCaptchaServiceServer(s, rs.captchaSvcServer)
	err = s.Serve(listener)
	if err != nil {
		panic(err)
	}
	log.Info("rpc server started.")

	rs.s = s
}

func (rs *rpcServer) Stop() {
	if rs.s != nil {
		rs.s.GracefulStop()
	}
}
