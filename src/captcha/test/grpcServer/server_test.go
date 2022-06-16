package grpcServer

import (
	"context"
	"testing"
	"time"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/rpc/server"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/service"
	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/biz/store"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hololee2cn/captcha/pkg/grpcIFace"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRPCServer(t *testing.T) {
	const Address = ":30002"
	captchaSvc := service.NewDefaultCaptchaSvc(store.NewMemoryStore(3))

	go func() {
		rs := server.NewRpcServer(Address, server.NewCaptchaSvcServer(captchaSvc))
		defer rs.Stop()
		rs.Start()

		t.Log("server end")
	}()

	time.Sleep(time.Second * 2) // wait some while
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err)
	defer conn.Close()

	client := grpcIFace.NewCaptchaServiceClient(conn)
	resp, err := client.Get(context.Background(), &grpcIFace.GetCaptchaRequest{
		Debug:  true,
		Length: 5,
	})
	assert.Nil(t, err)
	t.Logf("answer: %s", resp.GetAnswer())

	resp2, err := client.Verify(context.Background(), &grpcIFace.VerifyCaptchaRequest{
		ID:     resp.GetID(),
		Answer: resp.GetAnswer(),
	})
	assert.Nil(t, err)
	assert.True(t, resp2.Data)

	// 再测试一次难错误的
	resp, err = client.Get(context.Background(), &grpcIFace.GetCaptchaRequest{
		Debug:  true,
		Length: 5,
	})
	assert.Nil(t, err)
	t.Logf("answer: %s", resp.GetAnswer())

	resp2, err = client.Verify(context.Background(), &grpcIFace.VerifyCaptchaRequest{
		ID:     resp.GetID(),
		Answer: resp.GetAnswer() + "a", // 故意猜错
	})
	assert.Nil(t, err)
	assert.False(t, resp2.Data)
}
