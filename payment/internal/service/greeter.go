package service

import (
	"context"
	v12 "github.com/go-kratos/kratos-layout/bff/api/helloworld/v1"

	"github.com/go-kratos/kratos-layout/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v12.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v12.HelloRequest) (*v12.HelloReply, error) {
	g, err := s.uc.CreateGreeter(ctx, &biz.Greeter{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v12.HelloReply{Message: "Hello " + g.Hello}, nil
}
