/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-24 14:39:55
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-24 14:39:55
 */

package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"payment/pkg/common/enums"
	"payment/pkg/common/exception"
	"payment/pkg/common/util"
	"payment/pkg/facade"
	"payment/pkg/facade/dto"
	"payment/pkg/facade/impl"
	pb "payment/pkg/router/grpc/pb"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)


const (
	secretKey = "AEJP7gY2YO2aRTLrPW9kWtsBgDXUL3yS"
)

var (
	conMap sync.Map
) 

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		tradeDetailFacade: impl.NewTradeDetailFacadeImpl(),
	}
}

type GrpcServer struct {
	Port uint32 `yaml:"port"`
	PemPath string `yaml:"pem_path"`
	KeyPath string `yaml:"key_path"`

	serv *grpc.Server

	tradeDetailFacade facade.TradeDetailFacade
}

func (g *GrpcServer) CreateServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Port))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("failed to listen")
		return err
	}

	// TLS认证
	var opts []grpc.ServerOption

	creds, err := credentials.NewServerTLSFromFile(g.PemPath, g.KeyPath)
    if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("Failed to generate credentials")
		return err
    }

    opts = append(opts, grpc.Creds(creds))

    // 注册interceptor
    opts = append(opts, grpc.UnaryInterceptor(interceptor))

	g.serv = grpc.NewServer()
	pb.RegisterGreeterServer(g.serv, &GrpcServer{})

	// Register reflection service on gRPC server.
	reflection.Register(g.serv)
	log.Infof("Create grpc server success!port is %d", g.Port)
	if err := g.serv.Serve(lis); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatalf("failed to serve")
		return err
	}
	return nil
}

func (g *GrpcServer) StopServer() {
	if g.serv != nil {
		log.Errorln("rpc server stop")
		g.serv.Stop()
	}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    err := auth(ctx)
    if err != nil {
        return nil, err
    }
    // 继续处理请求
    return handler(ctx, req)
}


func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return grpc.Errorf(codes.Unauthenticated, "No token authentication information")
    }

	var (
		appid string
		authCode string
		timestamp string
	)
	if val, ok := md["appid"]; ok {
        appid = val[0]
    }
	if val, ok := md["authCode"]; ok {
        authCode = val[0]
    }
	if val, ok := md["timestamp"]; ok {
        timestamp = val[0]
    }

	if ok := checkAuthCode(appid, authCode, timestamp); !ok {
		return grpc.Errorf(codes.Unauthenticated, "Authentication failed")
	}
	return nil
}


/**
 * @Description: 校验授权码
 * @Author: Allen
 * @param {*} appid
 * @param {*} authCode
 * @param {string} timestamp
 * @return {*}
 * @error: 
 */
func checkAuthCode(appid, authCode, timestamp string) bool {
	if appid == "" || authCode == "" || timestamp == "" {
		log.WithFields(log.Fields{
			"appid": appid,
			"authCode": authCode,
			"timestamp": timestamp,
		}).Fatalf("Parameter cannot be empty")
		return false
	}
	//先查看conMap有无重复连接,如果有则根据时间戳参数确定连接是客户端新生成的
	if m, ok := conMap.Load(appid); ok {
		tm := m.(string)
		if tm == timestamp {
			log.WithFields(log.Fields{
				"timestamp": timestamp,
			}).Fatalf("Duplicate timestamp")
			return false
		}
	}

	code := util.Md5(fmt.Sprintf("%s;%s;%s", appid, secretKey, timestamp))
	if code == authCode {
		return true
	}
	return false
}

/**
 * @Description: 发起支付
 * @Author: Allen
 * @param {context.Context} ctx
 * @param {*pb.PayRequest} req
 * @return {*}
 * @error: 
 */
func (g *GrpcServer) Pay(ctx context.Context, req *pb.PayRequest) (*pb.PayReply, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	if util.IsEmptyOrNil(req) {
		return nil, exception.ERROR_PARAMETER
	}
	reply := &pb.PayReply{}
	reqDto := &dto.PayRequestDto{}
	err := util.Assemble(req, reqDto)
	reqDto.PayChannel = enums.PayChannelEnum(req.PayChannel)
	if err != nil {
		reply.Code = exception.ErrCode(err)
		reply.CodeMsg = exception.ErrMsg(err)
		return reply, err
	}
	result, err := g.tradeDetailFacade.Pay(reqDto)
	if err != nil {
		return reply, err
	}
	reply.Body = result
	return reply, nil
}

/**
 * @Description: 回调Ack
 * @Author: Allen
 * @param {context.Context} ctx
 * @param {*pb.AckRequest} req
 * @return {*}
 * @error: 
 */
func (g *GrpcServer) Ack(ctx context.Context, req *pb.AckRequest) (*pb.AckReply, error) {
	defer func() {
		if err := recover(); err != nil {
			log.WithFields(log.Fields{
				"err": err,
				"stack": util.CallStack(),
			}).Panicf("Method execution panic")
		}
	}()
	if util.IsEmptyOrNil(req) {
		return nil, exception.ERROR_PARAMETER
	}
	reply := &pb.AckReply{}
	reqDto := &dto.AckRequestDto{}
	err := util.Assemble(req, reqDto)
	if err != nil {
		reply.Code = exception.ErrCode(err)
		reply.CodeMsg = exception.ErrMsg(err)
		return reply, err
	}
	err = g.tradeDetailFacade.Ack(reqDto)
	if err != nil {
		reply.Code = exception.ErrCode(err)
		reply.CodeMsg = exception.ErrMsg(err)
		return reply, err
	}
	return reply, nil
}