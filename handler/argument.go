package handler

import (
	"github.com/EAHITechnology/inf/golang/log"
	"github.com/EAHITechnology/inf/golang/log_id"
	"github.com/EAHITechnology/inf/unique-number/logic"
	"github.com/EAHITechnology/inf/unique-number/proto/pb"
	"golang.org/x/net/context"
)

func parseGetUniqueNumberArg(ctx context.Context, r *pb.UnRequest) (context.Context, logic.GetUniqueNumberRequest, error) {
	newCtx := log_id.GenCtx(ctx)
	re := logic.GetUniqueNumberRequest{
		AppName: r.GetAppName(),
	}
	log.InfofCtx(newCtx, "request GetUniqueNumber body:%v", re)
	return newCtx, re, nil
}
