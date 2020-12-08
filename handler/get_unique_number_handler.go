package handler

import (
	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/unique-number/logic"
	"github.com/EAHITechnology/inf/unique-number/proto/pb"
	"golang.org/x/net/context"
)

type GetUniqueNumberServer struct{}

func (this *GetUniqueNumberServer) GetUniqueNumber(ctx context.Context, request *pb.UnRequest) (*pb.UnResponse, error) {
	resp := &pb.UnResponse{}
	newCtx, arg, err := parseGetUniqueNumberArg(ctx, request)
	if err != nil {
		log.ErrorfCtx(newCtx, "GetUniqueNumber parseGetUniqueNumberArg err:%s", err.Error())
		resp.ErrorCode = pb.ErrorCode_CLIENT_ERROR
		resp.ErrorMsg = "参数错误:" + err.Error()
		return resp, nil
	}
	if err := logic.GetUniqueNumberLogic(newCtx, arg); err != nil {
		log.ErrorfCtx(ctx, "GetUniqueNumber GetUniqueNumberLogic err:%s", err.Error())
		resp.ErrorCode = pb.ErrorCode_SERVER_ERROR
		resp.ErrorMsg = "服务器错误:" + err.Error()
		return resp, nil
	}
	resp.ErrorCode = pb.ErrorCode_SUCCESS
	resp.ErrorMsg = "请求成功"
	return resp, nil
}
