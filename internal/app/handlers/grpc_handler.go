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
