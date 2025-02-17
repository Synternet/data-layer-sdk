package main

import (
	"context"

	v1 "github.com/synternet/data-layer-sdk/examples/rpc/types/example/v1"
	"github.com/synternet/data-layer-sdk/pkg/service"
	"github.com/synternet/data-layer-sdk/x/synternet/rpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.UserServiceServer = (*UserService)(nil)

type UserService struct {
	v1.UnimplementedUserServiceServer
	pub   *service.Service
	users []*v1.User
}

// Add implements v1.UserServiceServer.
func (u *UserService) Add(_ context.Context, r *v1.AddRequest) (*v1.AddResponse, error) {
	user := &v1.User{
		Id:      uint64(len(u.users)),
		Name:    r.Name,
		Surname: r.Surname,
	}
	u.users = append(u.users, user)
	u.pub.Publish(&v1.RegistrationsResponse{User: user}, "users.registrations")
	u.pub.Logger.Info("User Add", "user", user)
	return &v1.AddResponse{
		Result: &v1.AddResponse_User{User: user},
	}, nil
}

// Get implements v1.UserServiceServer.
func (u *UserService) Get(_ context.Context, r *v1.GetRequest) (*v1.GetResponse, error) {
	u.pub.Logger.Info("User Get", "r", r)
	if r.Id >= uint64(len(u.users)) {
		return &v1.GetResponse{Result: &v1.GetResponse_Error{Error: &rpc.Error{Error: "invalid ID"}}}, nil
	}
	return &v1.GetResponse{Result: &v1.GetResponse_User{User: u.users[r.Id]}}, nil
}

// Registrations implements v1.UserServiceServer.
func (u *UserService) Registrations(context.Context, *emptypb.Empty) (*v1.RegistrationsResponse, error) {
	return nil, nil
}
