package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hololee2cn/wxpub/v1/src/captcha/internal/consts"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/hololee2cn/wxpub/v1/src/captcha/pkg/grpcIFace"

	"google.golang.org/grpc"
)

var (
	client pb.CaptchaServiceClient
)

// 先启动rpc服务, 再启动这个http服务, 然后访问 http://localhost:3333/ 测试
func main() {

	http.Handle("/", http.FileServer(http.Dir("./static")))

	// api for create captcha
	// 创建图像验证码api
	http.HandleFunc("/get", get)

	// api for verify captcha
	http.HandleFunc("/verify", verify)

	fmt.Println("Server is at localhost:3333")
	if err := http.ListenAndServe("localhost:3333", nil); err != nil {
		log.Fatal(err)
	}
}

func init() {
	conn, err := grpc.Dial(consts.RPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	// defer conn.Close()
	client = pb.NewCaptchaServiceClient(conn)
}

func get(w http.ResponseWriter, r *http.Request) {
	httpResp := HttpResponse{}
	line := r.URL.Query().Get("line")

	lineOptions, err := strconv.Atoi(line)
	if err != nil {
		httpResp.Code = 1
		httpResp.Msg = "wrong line options"
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(httpResp)
		return
	}

	rpcResp, err := client.Get(context.Background(), &pb.GetCaptchaRequest{
		MaxAge:          300,
		ShowLineOptions: int32(lineOptions),
	})

	if err != nil {
		httpResp.Code = 1
		httpResp.Msg = err.Error()
	} else {
		httpResp.Data = map[string]string{
			"id":    rpcResp.GetID(),
			"value": rpcResp.GetBase64Value(),
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(httpResp)
}

func verify(w http.ResponseWriter, r *http.Request) {
	var httpResp HttpResponse

	id := r.URL.Query().Get("id")
	answer := r.URL.Query().Get("answer")

	rpcResp, err := client.Verify(context.Background(), &pb.VerifyCaptchaRequest{
		ID:     id,
		Answer: answer,
	})
	if err != nil {
		httpResp.Code = 1
		httpResp.Msg = err.Error()
	} else {
		httpResp.Data = map[string]bool{
			"result": rpcResp.GetData(),
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(httpResp)
}

type RequestBody struct {
}

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
