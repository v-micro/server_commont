package serverorder

import (
	"context"
	"github.com/golang/glog"
	"google.golang.org/grpc"
	"server_order/server_common/comutil"
	"time"
)

type Client struct {
	Ctx      context.Context
	Conn     *grpc.ClientConn

	//应有的服务
	Ping PingClient
}

//公共方法
var grpcUrl = ""
var grpcClient Client

//初始化
func init() {
	grpcClient = Client{
		Ctx: context.Background(),
	}
}

//实例化
func GetClient() *Client {
	var err error

	//从服务发现获取链接地址
	var endpoints = []string{"192.168.59.131:2379"}
	ser := comutil.NewServiceDiscovery(endpoints)
	defer ser.Close()
	err = ser.WatchService("server_order")
	if err != nil {
		return nil
	}
	grpcUrl = ser.GetServices()[0]

	//grpcUrl 判断
	if grpcUrl == "" {
		glog.Error("grpcUrl 连接获取失败")
		return nil
	}

	//设置超时
	grpcClient.Ctx, _ = context.WithTimeout(context.Background(), 5*time.Minute)

	if grpcClient.Conn, err = grpc.Dial(grpcUrl, grpc.WithInsecure()); err != nil {
		glog.Error("grpc连接失败，", err , grpcUrl)
		return nil
	} else {
		//rpc链接服务
		grpcClient.Ping = NewPingClient(grpcClient.Conn)

		return &grpcClient
	}
}
