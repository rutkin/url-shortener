package handlers

import (
	context "context"

	"github.com/rutkin/url-shortener/internal/app/logger"
	"github.com/rutkin/url-shortener/internal/app/service"
	"go.uber.org/zap"
)

type GRPCHanlder struct {
	UnimplementedGRPCHandlerServer
	service service.Service
}

func NewGRPCHandler() (*GRPCHanlder, error) {
	s, err := service.NewURLService()
	if err != nil {
		logger.Log.Error("failed to create url service", zap.String("error", err.Error()))
		return nil, err
	}

	return &GRPCHanlder{service: s}, nil
}

func (grpc *GRPCHanlder) CreateURL(ctx context.Context, in *CreateURLRequest) (*CreateURLResponse, error) {
	var result CreateURLResponse
	shortURL, err := grpc.service.CreateURL([]byte(in.LongUrl), in.UserId)
	if err != nil {
		result.Error = err.Error()
	} else {
		result.ShortUrl = shortURL
	}
	return &result, nil
}

func (grpc *GRPCHanlder) CreateURLS(ctx context.Context, in *CreateURLSRequest) (*CreateURLSResponse, error) {
	var result CreateURLSResponse
	resp, err := grpc.service.CreateURLS(in.LongUrl, in.UserId)
	if err != nil {
		result.Error = err.Error()
	} else {
		result.ShortUrl = resp
	}
	return &result, nil
}

func (grpc *GRPCHanlder) GetURL(ctx context.Context, in *GetURLRequest) (*GetURLResponse, error) {
	var result GetURLResponse
	resp, err := grpc.service.GetURL(in.ShortUrl)
	if err != nil {
		result.Error = err.Error()
	} else {
		result.LongUrl = resp
	}
	return &result, nil
}

func (grpc *GRPCHanlder) DeleteURLS(ctx context.Context, in *DeleteURLSRequest) (*DeleteURLSResponse, error) {
	var result DeleteURLSResponse
	err := grpc.service.DeleteURLS(in.ShortUrl, in.UserId)
	if err != nil {
		result.Error = err.Error()
	}
	return &result, nil
}

func (grpc *GRPCHanlder) GetStats(context.Context, *Empty) (*GetStatsResponse, error) {
	var result GetStatsResponse
	resp, err := grpc.service.GetStats()
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Urls = int64(resp.URLS)
		result.Users = int64(resp.Users)
	}
	return &result, nil
}
