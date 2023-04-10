package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/sealur/user_service/internal/models"
	"github.com/Alexander272/sealur/user_service/internal/repo"
	"github.com/Alexander272/sealur/user_service/pkg/hasher"
	"github.com/Alexander272/sealur/user_service/pkg/logger"
	"github.com/Alexander272/sealur_proto/api/user/models/user_model"
	"github.com/Alexander272/sealur_proto/api/user/user_api"
	"github.com/google/uuid"
)

type UserService struct {
	repo   repo.Users
	hasher hasher.PasswordHasher
	role   Role
}

func NewUserService(repo repo.Users, hasher hasher.PasswordHasher, role Role) *UserService {
	return &UserService{
		repo:   repo,
		hasher: hasher,
		role:   role,
	}
}

func (s *UserService) Get(ctx context.Context, req *user_api.GetUser) (*user_model.User, error) {
	user, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}

	return user, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *user_api.GetUserByEmail) (*user_model.User, error) {
	user, password, err := s.repo.GetByEmail(ctx, req)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}

	if req.Password == "" {
		return user, nil
	}

	salt := strings.Split(password, ".")[1]
	pass, err := s.hasher.Hash(req.Password, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password. error: %w", err)
	}

	if fmt.Sprintf("%s.%s", pass, salt) != password {
		return nil, models.ErrPassword
	}

	return user, nil
}

func (s *UserService) GetManager(ctx context.Context, req *user_api.GetUser) (manager *user_api.Manager, err error) {
	user, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user. error: %w", err)
	}

	m := &user_model.User{}
	if user.ManagerId != uuid.Nil.String() {
		m, err = s.repo.Get(ctx, &user_api.GetUser{Id: user.ManagerId})
		if err != nil {
			return nil, fmt.Errorf("failed to get manager. error: %w", err)
		}
	} else {
		logger.Debug("not implemented")
		// TODO получить id менеджера по региону пользователя
	}

	manager = &user_api.Manager{
		Email: m.Email,
		User:  user,
	}

	return manager, nil
}

func (s *UserService) GetManagers(ctx context.Context, req *user_api.GetNewUser) ([]*user_model.User, error) {
	users, err := s.repo.GetManagers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get managers. error: %w", err)
	}
	return users, nil
}

func (s *UserService) Create(ctx context.Context, user *user_api.CreateUser) (string, error) {
	candidate, _, err := s.repo.GetByEmail(ctx, &user_api.GetUserByEmail{Email: user.Email})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("failed to get user. error: %w", err)
	}

	if candidate != nil {
		return "", models.ErrUserExist
	}

	role, err := s.role.GetDefault(ctx)
	if err != nil {
		return "", err
	}

	salt, err := s.hasher.GenerateSalt()
	if err != nil {
		return "", fmt.Errorf("failed to create salt. error: %w", err)
	}
	pass, err := s.hasher.Hash(user.Password, salt)
	if err != nil {
		return "", fmt.Errorf("failed to hash password. error: %w", err)
	}
	user.Password = fmt.Sprintf("%s.%s", pass, salt)

	id, err := s.repo.Create(ctx, user, role.Id)
	if err != nil {
		return "", fmt.Errorf("failed to create user. error: %w", err)
	}
	return id, nil
}

func (s *UserService) Confirm(ctx context.Context, user *user_api.ConfirmUser) (*user_model.User, error) {
	if err := s.repo.Confirm(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to confirm user. error: %w", err)
	}

	u, err := s.Get(ctx, &user_api.GetUser{Id: user.Id})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) SetManager(ctx context.Context, manager *user_api.UserManager) error {
	if err := s.repo.SetManager(ctx, manager); err != nil {
		return fmt.Errorf("failed to set manager. error: %w", err)
	}
	return nil
}

func (s *UserService) Update(ctx context.Context, user *user_api.UpdateUser) error {
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

	if err := s.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user. error: %w", err)
	}

	return nil
}

// type UserService struct {
// 	hasher   hasher.PasswordHasher
// 	userRepo repo.Users
// 	roleRepo repo.Role
// 	ipRepo   repo.IP
// 	email    email_api.EmailServiceClient
// }

// func NewUserService(user repo.Users, role repo.Role, ip repo.IP, hasher hasher.PasswordHasher, email email_api.EmailServiceClient) *UserService {
// 	return &UserService{
// 		userRepo: user,
// 		roleRepo: role,
// 		ipRepo:   ip,
// 		hasher:   hasher,
// 		email:    email,
// 	}
// }

// func (s *UserService) Get(ctx context.Context, req *user_api.GetUserRequest) (u *user_api.User, err error) {
// 	user, err := s.userRepo.Get(ctx, req)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, err
// 		}
// 		return nil, fmt.Errorf("failed to get user. error: %w", err)
// 	}

// 	if req.Login != "" {
// 		salt := strings.Split(user.Password, ".")[1]
// 		pass, err := s.hasher.Hash(req.Password, salt)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to hash password. error: %w", err)
// 		}

// 		if fmt.Sprintf("%s.%s", pass, salt) != user.Password {
// 			return nil, models.ErrPassword
// 		}
// 	}

// 	roles, err := s.roleRepo.Get(ctx, &user_api.GetRolesRequest{UserId: user.Id})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get roles. error: %w", err)
// 	}

// 	var userRoles []*user_api.Role
// 	for _, r := range roles {
// 		ur := user_api.Role{
// 			Id:      r.Id,
// 			Service: r.Service,
// 			Role:    r.Role,
// 		}
// 		userRoles = append(userRoles, &ur)
// 	}

// 	u = &user_api.User{
// 		Id:           user.Id,
// 		Organization: user.Organization,
// 		Name:         user.Name,
// 		Email:        user.Email,
// 		City:         user.City,
// 		Position:     user.Position,
// 		Phone:        user.Phone,
// 		Roles:        userRoles,
// 	}

// 	return u, nil
// }

// func (s *UserService) GetAll(ctx context.Context, req *user_api.GetAllUserRequest) ([]*user_api.User, int, error) {
// 	users, err := s.userRepo.GetAll(ctx, req)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get all users. error: %w", err)
// 	}

// 	if len(users) == 0 {
// 		return nil, 0, models.ErrUsersEmpty
// 	}

// 	roles, err := s.roleRepo.GetAll(ctx, &user_api.GetAllRolesRequest{})
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get all roles. error: %w", err)
// 	}

// 	ips, err := s.ipRepo.GetAll(ctx, &user_api.GetAllIpRequest{})
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to get all ip. error: %w", err)
// 	}

// 	var u []*user_api.User
// 	for _, item := range users {
// 		var userRoles []*user_api.Role
// 		for j := 0; j < len(roles); j++ {
// 			if roles[j].UserId == item.Id {
// 				ur := user_api.Role{
// 					Id:      roles[j].Id,
// 					Service: roles[j].Service,
// 					Role:    roles[j].Role,
// 				}
// 				userRoles = append(userRoles, &ur)
// 			}
// 		}

// 		var userIp []*user_api.Ip
// 		for j := 0; j < len(ips); j++ {
// 			if ips[j].UserId == item.Id {
// 				ip := user_api.Ip{
// 					Ip:   ips[j].Ip,
// 					Date: ips[j].Date,
// 				}
// 				userIp = append(userIp, &ip)
// 			}
// 		}

// 		user := user_api.User{
// 			Id:           item.Id,
// 			Organization: item.Organization,
// 			Name:         item.Name,
// 			Email:        item.Email,
// 			City:         item.City,
// 			Position:     item.Position,
// 			Phone:        item.Phone,
// 			Login:        item.Login,
// 			Roles:        userRoles,
// 			Ip:           userIp,
// 		}
// 		u = append(u, &user)
// 	}

// 	return u, users[0].Count, nil
// }

// func (s *UserService) GetNew(ctx context.Context, req *user_api.GetNewUserRequest) ([]*user_api.User, error) {
// 	users, err := s.userRepo.GetNew(ctx, req)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get new users. error: %w", err)
// 	}

// 	if len(users) == 0 {
// 		return []*user_api.User{}, nil
// 	}

// 	var u []*user_api.User
// 	for _, item := range users {
// 		user := user_api.User{
// 			Id:           item.Id,
// 			Organization: item.Organization,
// 			Name:         item.Name,
// 			Email:        item.Email,
// 			City:         item.City,
// 			Position:     item.Position,
// 			Phone:        item.Phone,
// 		}
// 		u = append(u, &user)
// 	}

// 	return u, nil
// }

// func (s *UserService) Create(ctx context.Context, user *user_api.CreateUserRequest) (*user_api.SuccessResponse, error) {
// 	if err := s.userRepo.Create(ctx, user); err != nil {
// 		return nil, fmt.Errorf("failed to create user. error: %w", err)
// 	}

// 	_, err := s.email.SendConfirm(ctx, &email_api.ConfirmUserRequest{
// 		Organization: user.Organization,
// 		Name:         user.Name,
// 		Position:     user.Position,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send email. error: %w", err)
// 	}

// 	return &user_api.SuccessResponse{Success: true}, nil
// }

// func (s *UserService) Confirm(ctx context.Context, user *user_api.ConfirmUserRequest) (*user_api.SuccessResponse, error) {
// 	candidate, err := s.userRepo.Get(ctx, &user_api.GetUserRequest{Login: user.Login})
// 	if err != nil && !errors.Is(err, sql.ErrNoRows) {
// 		return nil, fmt.Errorf("failed to get user. error: %w", err)
// 	}
// 	if (candidate != models.User{}) {
// 		return nil, models.ErrUserExist
// 	}

// 	origPas := user.Password

// 	salt, err := s.hasher.GenerateSalt()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create salt. error: %w", err)
// 	}
// 	pass, err := s.hasher.Hash(user.Password, salt)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to hash password. error: %w", err)
// 	}
// 	user.Password = fmt.Sprintf("%s.%s", pass, salt)

// 	confirmUser, err := s.userRepo.Confirm(ctx, user)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to verify user. error: %w", err)
// 	}

// 	services := make([]string, 0, len(user.Roles))
// 	var roles []*user_api.CreateRoleRequest

// 	for _, item := range user.Roles {
// 		role := user_api.CreateRoleRequest{
// 			UserId:  user.Id,
// 			Service: item.Service,
// 			Role:    item.Role,
// 		}
// 		roles = append(roles, &role)

// 		switch item.Service {
// 		case "pro":
// 			services = append(services, "Sealur Pro")
// 		case "moment":
// 			services = append(services, "Расчет момента затяжки")
// 		}
// 	}

// 	if err := s.roleRepo.Create(ctx, roles); err != nil {
// 		return nil, fmt.Errorf("failed to create roles. error: %w", err)
// 	}

// 	_, err = s.email.SendJoin(ctx, &email_api.JoinUserRequest{
// 		Name:     confirmUser.Name,
// 		Login:    user.Login,
// 		Password: origPas,
// 		Email:    confirmUser.Email,
// 		Services: services,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to send email. error: %w", err)
// 	}

// 	return &user_api.SuccessResponse{Success: true}, nil
// }

// func (s *UserService) Update(ctx context.Context, user *user_api.UpdateUserRequest) error {
// 	if user.Password != "" {
// 		salt, err := s.hasher.GenerateSalt()
// 		if err != nil {
// 			return fmt.Errorf("failed to create salt. error: %w", err)
// 		}
// 		pass, err := s.hasher.Hash(user.Password, salt)
// 		if err != nil {
// 			return fmt.Errorf("failed to hash password. error: %w", err)
// 		}
// 		user.Password = fmt.Sprintf("%s.%s", pass, salt)
// 	}

// 	if err := s.userRepo.Update(ctx, user); err != nil {
// 		return fmt.Errorf("failed to update user. error: %w", err)
// 	}

// 	return nil
// }

// func (s *UserService) Delete(ctx context.Context, user *user_api.DeleteUserRequest) error {
// 	if _, err := s.userRepo.Delete(ctx, user); err != nil {
// 		return fmt.Errorf("failed to delete user. error: %w", err)
// 	}

// 	return nil
// }

// func (s *UserService) Reject(ctx context.Context, user *user_api.DeleteUserRequest) error {
// 	deleteUser, err := s.userRepo.Delete(ctx, user)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete user. error: %w", err)
// 	}

// 	_, err = s.email.SendReject(ctx, &email_api.RejectUserRequest{
// 		Name:  deleteUser.Name,
// 		Email: deleteUser.Email,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to send email. error: %w", err)
// 	}

// 	return nil
// }
