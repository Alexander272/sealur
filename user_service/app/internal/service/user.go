package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur/user_service/internal/repo"
	proto_user "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto"
	proto_email "github.com/Alexander272/sealur/user_service/internal/transport/grpc/proto/email"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
)

type UserService struct {
	hasher   hasher.PasswordHasher
	userRepo repo.Users
	roleRepo repo.Role
	ipRepo   repo.IP
	email    proto_email.EmailServiceClient
}

func NewUserService(user repo.Users, role repo.Role, ip repo.IP, hasher hasher.PasswordHasher, email proto_email.EmailServiceClient) *UserService {
	return &UserService{
		userRepo: user,
		roleRepo: role,
		ipRepo:   ip,
		hasher:   hasher,
		email:    email,
	}
}

func (s *UserService) Get(ctx context.Context, req *proto_user.GetUserRequest) (u *proto_user.User, err error) {
	user, err := s.userRepo.Get(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}

	if req.Login != "" {
		salt := strings.Split(user.Password, ".")[1]
		pass, err := s.hasher.Hash(req.Password, salt)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password. error: %w", err)
		}

		if fmt.Sprintf("%s.%s", pass, salt) != user.Password {
			return nil, models.ErrPassword
		}
	}

	roles, err := s.roleRepo.Get(ctx, &proto_user.GetRolesRequest{UserId: user.Id})
	if err != nil {
		return nil, fmt.Errorf("failed to get roles. error: %w", err)
	}

	var userRoles []*proto_user.Role
	for _, r := range roles {
		ur := proto_user.Role{
			Id:      r.Id,
			Service: r.Service,
			Role:    r.Role,
		}
		userRoles = append(userRoles, &ur)
	}

	u = &proto_user.User{
		Id:           user.Id,
		Organization: user.Organization,
		Name:         user.Name,
		Email:        user.Email,
		City:         user.City,
		Position:     user.Position,
		Phone:        user.Phone,
		Roles:        userRoles,
	}

	return u, nil
}

func (s *UserService) GetAll(ctx context.Context, req *proto_user.GetAllUserRequest) ([]*proto_user.User, error) {
	users, err := s.userRepo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users. error: %w", err)
	}

	if len(users) == 0 {
		return nil, models.ErrUsersEmpty
	}

	roles, err := s.roleRepo.GetAll(ctx, &proto_user.GetAllRolesRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles. error: %w", err)
	}

	ips, err := s.ipRepo.GetAll(ctx, &proto_user.GetAllIpRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get all ip. error: %w", err)
	}

	var u []*proto_user.User
	for i, item := range users {
		var userRoles []*proto_user.Role
		for j := i; j < len(roles); j++ {
			if roles[j].UserId == item.Id {
				ur := proto_user.Role{
					Id:      roles[j].Id,
					Service: roles[j].Service,
					Role:    roles[j].Role,
				}
				userRoles = append(userRoles, &ur)
			}
		}

		var userIp []*proto_user.Ip
		for j := 0; j < len(ips); j++ {
			if ips[j].UserId == item.Id {
				ip := proto_user.Ip{
					Ip:   ips[j].Ip,
					Date: ips[j].Date,
				}
				userIp = append(userIp, &ip)
			}
		}

		user := proto_user.User{
			Id:           item.Id,
			Organization: item.Organization,
			Name:         item.Name,
			Email:        item.Email,
			City:         item.City,
			Position:     item.Position,
			Phone:        item.Phone,
			Login:        item.Login,
			Roles:        userRoles,
			Ip:           userIp,
		}
		u = append(u, &user)
	}

	return u, nil
}

func (s *UserService) GetNew(ctx context.Context, req *proto_user.GetNewUserRequest) ([]*proto_user.User, error) {
	users, err := s.userRepo.GetNew(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get new users. error: %w", err)
	}

	if len(users) == 0 {
		return []*proto_user.User{}, nil
	}

	var u []*proto_user.User
	for _, item := range users {
		user := proto_user.User{
			Id:           item.Id,
			Organization: item.Organization,
			Name:         item.Name,
			Email:        item.Email,
			City:         item.City,
			Position:     item.Position,
			Phone:        item.Phone,
		}
		u = append(u, &user)
	}

	return u, nil
}

func (s *UserService) Create(ctx context.Context, user *proto_user.CreateUserRequest) (*proto_user.SuccessResponse, error) {
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user. error: %w", err)
	}

	_, err := s.email.SendConfirm(ctx, &proto_email.ConfirmUserRequest{
		Organization: user.Organization,
		Name:         user.Name,
		Position:     user.Position,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send email. error: %w", err)
	}

	return &proto_user.SuccessResponse{Success: true}, nil
}

func (s *UserService) Confirm(ctx context.Context, user *proto_user.ConfirmUserRequest) (*proto_user.SuccessResponse, error) {
	candidate, err := s.userRepo.Get(ctx, &proto_user.GetUserRequest{Login: user.Login})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}
	if (candidate != models.User{}) {
		return nil, models.ErrUserExist
	}

	origPas := user.Password

	salt, err := s.hasher.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("failed to create salt. error: %w", err)
	}
	pass, err := s.hasher.Hash(user.Password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password. error: %w", err)
	}
	user.Password = fmt.Sprintf("%s.%s", pass, salt)

	confirmUser, err := s.userRepo.Confirm(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to verify user. error: %w", err)
	}

	services := make([]string, 0, len(user.Roles))
	var roles []*proto_user.CreateRoleRequest

	for _, item := range user.Roles {
		role := proto_user.CreateRoleRequest{
			UserId:  user.Id,
			Service: item.Service,
			Role:    item.Role,
		}
		roles = append(roles, &role)

		switch item.Service {
		case "pro":
			services = append(services, "Sealur Pro")
		case "moment":
			services = append(services, "Расчет момента затяжки")
		}
	}

	if err := s.roleRepo.Create(ctx, roles); err != nil {
		return nil, fmt.Errorf("failed to create roles. error: %w", err)
	}

	_, err = s.email.SendJoin(ctx, &proto_email.JoinUserRequest{
		Name:     confirmUser.Name,
		Login:    user.Login,
		Password: origPas,
		Email:    confirmUser.Email,
		Services: services,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send email. error: %w", err)
	}

	return &proto_user.SuccessResponse{Success: true}, nil
}

func (s *UserService) Update(ctx context.Context, user *proto_user.UpdateUserRequest) error {
	if user.Password != "" {
		salt, err := s.hasher.GenerateSalt()
		if err != nil {
			return fmt.Errorf("failed to create salt. error: %w", err)
		}
		pass, err := s.hasher.Hash(user.Password, salt)
		if err != nil {
			return fmt.Errorf("failed to hash password. error: %w", err)
		}
		user.Password = fmt.Sprintf("%s.%s", pass, salt)
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user. error: %w", err)
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, user *proto_user.DeleteUserRequest) error {
	if _, err := s.userRepo.Delete(ctx, user); err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}

	return nil
}

func (s *UserService) Reject(ctx context.Context, user *proto_user.DeleteUserRequest) error {
	deleteUser, err := s.userRepo.Delete(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}

	_, err = s.email.SendReject(ctx, &proto_email.RejectUserRequest{
		Name:  deleteUser.Name,
		Email: deleteUser.Email,
	})
	if err != nil {
		return fmt.Errorf("failed to send email. error: %w", err)
	}

	return nil
}
