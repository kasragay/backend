package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"os"

	"errors"

	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/ports"
	"github.com/kasragay/backend/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

const relationalCaller = packageCaller + ".Relational"

var user_back_maximum_reference int64

func init() {
	user_back_maximum_reference_ := os.Getenv("USER_BACK_MAXIMUM_REFERENCE")
	if user_back_maximum_reference_ == "" {
		user_back_maximum_reference = 10
	} else {
		var err error
		user_back_maximum_reference, err = strconv.ParseInt(user_back_maximum_reference_, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("failed to parse USER_BACK_MAXIMUM_REFERENCE: %s", err))
		}
	}
}

type Relational struct {
	logger *utils.Logger
	client *gorm.DB
}

func NewRelationalRepo(logger *utils.Logger) ports.RelationalRepo {
	database := os.Getenv("POSTGRES_DB_DATABASE")
	if database == "" {
		logger.Fatal(context.Background(), "POSTGRES_DB_DATABASE is not set")
	}
	password := os.Getenv("POSTGRES_DB_PASSWORD")
	if password == "" {
		logger.Fatal(context.Background(), "POSTGRES_DB_PASSWORD is not set")
	}
	username := os.Getenv("POSTGRES_DB_USERNAME")
	if username == "" {
		logger.Fatal(context.Background(), "POSTGRES_DB_USERNAME is not set")
	}
	port := os.Getenv("POSTGRES_DB_PORT")
	if port == "" {
		logger.Fatal(context.Background(), "POSTGRES_DB_PORT is not set")
	}
	host := os.Getenv("POSTGRES_DB_HOST")
	if host == "" {
		logger.Fatal(context.Background(), "POSTGRES_DB_HOST is not set")
	}
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC", host, port, username, password, database)), &gorm.Config{
		Logger: gLogger.Default.LogMode(gLogger.Silent),
	})
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to connect to database: %v", err)
	}
	instance, err := db.DB()
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to get database instance: %v", err)
	}
	if err := instance.Ping(); err != nil {
		logger.Fatalf(context.Background(), "Failed to ping database: %v", err)
	}
	return &Relational{
		logger: logger,
		client: db,
	}
}

func (s *Relational) TODO_DELETE_DropTables() error {
	return s.client.Migrator().DropTable(
		&ports.AdminUserModel{},
		&ports.ClientUserModel{},
	)
}

func (s *Relational) TODO_DELETE_AddFirstUsers() error {
	_, err := s.CreateUser(
		context.Background(),
		&ports.AuthSignupPostRequest{
			Name:        "Amir Hossein Pakdaman",
			PhoneNumber: "09202400120",
			UserType:    ports.AdminUserType,
			Password:    "P@ssw0rdUnhackable",
		}, uuid.MustParse("89950e97-6b0f-4ef4-977a-8378b14bc4a7"),
	)
	if err != nil {
		return err
	}
	_, err = s.CreateUser(
		context.Background(),
		&ports.AuthSignupPostRequest{
			Name:        "Amir Hossein Pakdaman",
			PhoneNumber: "09202400120",
			UserType:    ports.ClientUserType,
			Password:    "P@ssw0rdUnhackable",
		}, uuid.MustParse("779033a2-4eaa-4817-aa81-24e24bd419f5"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Relational) AutoMigrate() error {
	return s.client.AutoMigrate(
		&ports.AdminUserModel{},
		&ports.ClientUserModel{},
	)
}

func (s *Relational) UserExists(ctx context.Context, req *ports.AuthCheckPostRequest) (resp *ports.AuthCheckPostResponse, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UserExists", err) }()
	user := ports.UserModelFromUserType(req.UserType)
	if err := s.client.WithContext(ctx).Where("username = ?", req.Username).First(user).Error; err == nil {
		if user.GetIsDeleted() {
			return &ports.AuthCheckPostResponse{
				Exists:         true,
				Deleted:        true,
				HasEmail:       false,
				HasPhoneNumber: false,
				HasPassword:    false,
			}, true, nil
		}
		return &ports.AuthCheckPostResponse{
			Exists:         true,
			Deleted:        false,
			HasEmail:       user.GetEmail() != "",
			HasPhoneNumber: user.GetPhoneNumber() != "",
			HasPassword:    user.GetPassword() != "",
		}, false, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, false, err
	}
	return &ports.AuthCheckPostResponse{
		Exists:         false,
		Deleted:        false,
		HasEmail:       false,
		HasPhoneNumber: false,
		HasPassword:    false,
	}, false, nil
}

func (s *Relational) UserExistsById(ctx context.Context, id uuid.UUID, userType ports.UserType) (exists bool, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UserExistsById", err) }()
	user := ports.UserModelFromUserType(userType)
	if err := s.client.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, false, nil
		}
		return false, false, err
	}
	if user.GetIsDeleted() {
		return true, true, nil
	}
	return true, false, nil
}

func (s *Relational) CreateUser(ctx context.Context, req *ports.AuthSignupPostRequest, forceId ...uuid.UUID) (resp *ports.User, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".CreateUser", err) }()
	user := ports.UserModelFromUserType(req.UserType)
	err = s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			existsUser := ports.UserModelFromUserType(req.UserType)
			if err := tx.WithContext(ctx).Where("username = ?", req.Username).First(existsUser).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}
			} else {
				return utils.UsernameAlreadyExistsResponse.Clone().
					WithReason("username", req.Username)
			}
			hashedPass, err := utils.HashPassword(req.Password)
			if err != nil {
				return err
			}
			userId := uuid.New()
			if len(forceId) > 0 {
				userId = forceId[0]
			}
			user = user.New(userId, req.Username, req.Name, req.PhoneNumber, req.Email, hashedPass, req.Avatar != "")
			if err := tx.WithContext(ctx).Create(user).Error; err != nil {
				return err
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	return user.ToUser(), nil
}

func (s *Relational) AuthUserPassword(ctx context.Context, req *ports.AuthSigninPasswordPostRequest) (resp *ports.User, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".AuthUserPassword", err) }()
	user := ports.UserModelFromUserType(req.UserType)
	if err := s.client.WithContext(ctx).Where("username = ?", req.Username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.UsernameNotFoundResponse.Clone().
				WithReason("username", req.Username)
		}
		return nil, err
	}
	if !utils.VerifyPassword(user.GetPassword(), req.Password) {
		return nil, utils.PasswordIncorrectResponse.Clone()
	}
	return user.ToUser(), nil
}

func (s *Relational) GetUserById(ctx context.Context, id uuid.UUID, userType ports.UserType) (user ports.UserModel, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".GetUserById", err) }()
	user = ports.UserModelFromUserType(userType)
	if err := s.client.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if user.GetIsDeleted() {
		return nil, true, nil
	}
	return user, false, nil
}

func (s *Relational) GetUserByUsername(ctx context.Context, username string, userType ports.UserType) (user ports.UserModel, isDeleted bool, err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".GetUserByUsername", err) }()
	user = ports.UserModelFromUserType(userType)
	if err := s.client.WithContext(ctx).Where("username = ?", username).First(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if user.GetIsDeleted() {
		return nil, true, nil
	}
	return user, false, nil
}

func (s *Relational) UpdateUserPasswordById(ctx context.Context, id uuid.UUID, userType ports.UserType, password string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdateUserPasswordById", err) }()
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := ports.UserModelFromUserType(userType)
			if err := tx.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.UserNotFoundResponse.Clone().
						WithReason("id", id.String())
				}
				return err
			}
			if err := tx.WithContext(ctx).Model(user).Updates(
				map[string]any{
					"password":   hashedPass,
					"updated_at": time.Now().UTC(),
				},
			).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Relational) UpdateUserPasswordByUsername(ctx context.Context, username string, userType ports.UserType, password string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdateUserPasswordByUsername", err) }()
	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return err
	}
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := ports.UserModelFromUserType(userType)
			if err := tx.WithContext(ctx).Where("username = ?", username).First(user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.UsernameNotFoundResponse.Clone().
						WithReason("username", username)
				}
				return err
			}
			if err := tx.WithContext(ctx).Model(user).Updates(
				map[string]any{
					"password":   hashedPass,
					"updated_at": time.Now().UTC(),
				},
			).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Relational) UpdateUserProfileById(ctx context.Context, id uuid.UUID, username, name, avatar string, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdateUserProfileById", err) }()
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := ports.UserModelFromUserType(userType)
			if err := tx.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.UserNotFoundResponse.Clone().
						WithReason("id", id.String())
				}
				return err
			}
			if err := tx.WithContext(ctx).Model(user).Updates(
				map[string]any{
					"username":   username,
					"name":       name,
					"has_avatar": avatar != "",
					"updated_at": time.Now().UTC(),
				},
			).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Relational) DeleteUserById(ctx context.Context, id uuid.UUID, userType ports.UserType) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".DeleteUserById", err) }()
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			return ports.UserModelFromUserType(userType).Delete(ctx, tx, id)
		},
	)
}

func (s *Relational) UpdateUserPhoneById(ctx context.Context, id uuid.UUID, userType ports.UserType, phoneNumber string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdateUserPhoneById", err) }()
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := ports.UserModelFromUserType(userType)
			if err := tx.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.UserNotFoundResponse.Clone().
						WithReason("id", id.String())
				}
				return err
			}
			var count int64
			if err := tx.WithContext(ctx).Where("phone_number = ?", phoneNumber).Count(&count).Error; err != nil {
				return err
			} else if count >= user_back_maximum_reference {
				return utils.PhoneNumberIsTakenByMultipleResponse.Clone().
					WithReason("phone_number", phoneNumber)
			}
			if err := tx.WithContext(ctx).Model(user).Updates(
				map[string]any{
					"phone_number": phoneNumber,
					"updated_at":   time.Now().UTC(),
				},
			).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Relational) UpdateUserEmailById(ctx context.Context, id uuid.UUID, userType ports.UserType, email string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".UpdateUserEmailById", err) }()
	return s.client.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := ports.UserModelFromUserType(userType)
			if err := tx.WithContext(ctx).Where("id = ?", id).First(user).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utils.UserNotFoundResponse.Clone().
						WithReason("id", id.String())
				}
				return err
			}
			var count int64
			if err := tx.WithContext(ctx).Where("email = ?", email).Count(&count).Error; err != nil {
				return err
			} else if count >= user_back_maximum_reference {
				return utils.EmailIsTakenByMultipleResponse.Clone().
					WithReason("email", email)
			}
			if err := tx.WithContext(ctx).Model(user).Updates(
				map[string]any{
					"email":      email,
					"updated_at": time.Now().UTC(),
				},
			).Error; err != nil {
				return err
			}
			return nil
		},
	)
}

func (s *Relational) CheckUserEmailLimit(ctx context.Context, email string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".CheckUserEmailLimit", err) }()
	var count int64
	if err := s.client.WithContext(ctx).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	} else if count >= user_back_maximum_reference {
		return utils.EmailIsTakenByMultipleResponse.Clone().
			WithReason("email", email)
	}
	return nil
}
func (s *Relational) CheckUserPhoneLimit(ctx context.Context, phoneNumber string) (err error) {
	defer func() { err = utils.FuncPipe(relationalCaller+".CheckUserPhoneLimit", err) }()
	var count int64
	if err := s.client.WithContext(ctx).Where("phone_number = ?", phoneNumber).Count(&count).Error; err != nil {
		return err
	} else if count >= user_back_maximum_reference {
		return utils.PhoneNumberIsTakenByMultipleResponse.Clone().
			WithReason("phone_number", phoneNumber)
	}
	return nil
}

func (s *Relational) Close() error {
	s.logger.Info(context.Background(), "Closing database connection")
	instance, err := s.client.DB()
	if err != nil {
		return err
	}
	return instance.Close()
}
