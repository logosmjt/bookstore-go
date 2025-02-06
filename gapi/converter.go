package gapi

import (
	db "github.com/logosmjt/bookstore-go/db/sqlc"
	"github.com/logosmjt/bookstore-go/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}
