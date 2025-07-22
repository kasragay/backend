package ports

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kasragay/backend/internal/utils"
	"gorm.io/gorm"
)

var s3Prefix string

// #TODO: it's make a Error
// go run cmd/tui/main.go
// panic: DOMAIN is not set

// goroutine 1 [running]:
// github.com/kasragay/backend/internal/ports.init.0()
//         /home/arman/bemula/kasragay/internal/ports/models.user.go:19 +0xe8
// exit status 2
// make: *** [Makefile:177: tui] Error 1

// func init() {
// 	domain := os.Getenv("DOMAIN")
// 	if domain == "" {
// 		panic("DOMAIN is not set")
// 	}
// 	version := os.Getenv("VERSION")
// 	if version == "" {
// 		panic("VERSION is not set")
// 	}
// 	s3Prefix = "https://api." + domain + "/" + version + "/s3/avatars/"
// }

func GetAvatarUrl(id uuid.UUID, userType UserType) string {
	return s3Prefix + string(userType) + "/" + id.String() + ".png"
}

func UserModelFromUserType(userType UserType) UserModel {
	switch userType {
	case AdminUserType:
		return &AdminUserModel{}
	case ClientUserType:
		return &ClientUserModel{}
	default:
		panic("unknown userType")
	}
}

type UserModel interface {
	New(id uuid.UUID, username, name, phoneNumber, email, password string, hasAvatar bool) UserModel
	ToUser() *User

	GetId() uuid.UUID
	GetUsername() string
	GetName() string
	GetPhoneNumber() string
	GetEmail() string
	GetPassword() string
	GetHasAvatar() bool
	GetUpdatedAt() time.Time
	GetCreatedAt() time.Time
	GetIsDeleted() bool

	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) (err error)
}

type BaseUserModel struct {
	Id          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	Username    string    `json:"username" gorm:"unique"`
	Name        string    `json:"name" gorm:"not null"`
	HasAvatar   bool      `json:"has_avatar" gorm:"not null"`
	PhoneNumber *string   `json:"phone_number"`
	Email       *string   `json:"email"`
	Password    *string   `json:"password"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	IsDeleted   bool      `json:"is_deleted" gorm:"not null"`
}

func newBaseUserModelWithId(id uuid.UUID, username, name, phoneNumber, email, password string, hasAvatar bool) *BaseUserModel {
	now := time.Now().UTC()
	return &BaseUserModel{
		Id:          id,
		Username:    username,
		Name:        name,
		HasAvatar:   hasAvatar,
		PhoneNumber: utils.StringSetNullIfEmpty(phoneNumber),
		Email:       utils.StringSetNullIfEmpty(email),
		Password:    utils.StringSetNullIfEmpty(password),
		UpdatedAt:   now,
		CreatedAt:   now,
	}
}

func (u BaseUserModel) GetId() uuid.UUID {
	return u.Id
}

func (u BaseUserModel) GetUsername() string {
	return u.Username
}

func (u BaseUserModel) GetName() string {
	return u.Name
}

func (u BaseUserModel) GetPhoneNumber() string {
	return utils.GetStringIfNotNull(u.PhoneNumber)
}

func (u BaseUserModel) GetEmail() string {
	return utils.GetStringIfNotNull(u.Email)
}

func (u BaseUserModel) GetPassword() string {
	return utils.GetStringIfNotNull(u.Password)
}

func (u BaseUserModel) GetHasAvatar() bool {
	return u.HasAvatar
}

func (u BaseUserModel) GetCreatedAt() time.Time {
	return u.CreatedAt
}

func (u BaseUserModel) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

func (u BaseUserModel) GetIsDeleted() bool {
	return u.IsDeleted
}

type AdminUserModel struct {
	BaseUserModel `json:",inline" gorm:"embedded"`
}

func (u AdminUserModel) New(id uuid.UUID, username, name, phoneNumber, email, password string, hasAvatar bool) UserModel {
	return &AdminUserModel{
		BaseUserModel: *newBaseUserModelWithId(id, username, name, phoneNumber, email, password, hasAvatar),
	}
}

func (u AdminUserModel) ToUser() *User {
	var avatar string
	if u.HasAvatar {
		avatar = GetAvatarUrl(u.Id, AdminUserType)
	}
	return &User{
		Id:       u.Id,
		Username: u.Username,
		Name:     u.Name,
		Avatar:   avatar,
		UserType: AdminUserType,
	}
}

func (u AdminUserModel) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) (err error) {
	defer func() {
		const caller = packageCaller + "AdminUserModel.Delete"
		if err != nil {
			var uErr *utils.Error
			if errors.As(err, &uErr) {
				err = uErr.WithCaller(caller)
			} else {
				err = utils.NewInternalError(err).WithCaller(caller)
			}
		}
	}()
	if err_ := tx.WithContext(ctx).Model(u).Where("id = ?", id).Updates(
		map[string]any{
			"email":      nil,
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	).Error; err_ != nil {
		if errors.Is(err_, gorm.ErrRecordNotFound) {
			return utils.UserNotFoundResponse.Clone().
				WithReason("id", id.String())
		}
		return err_
	}
	return nil
}

type ClientUserModel struct {
	BaseUserModel `json:",inline" gorm:"embedded"`
}

func (u ClientUserModel) New(id uuid.UUID, username, name, phoneNumber, email, password string, hasAvatar bool) UserModel {
	return &ClientUserModel{
		BaseUserModel: *newBaseUserModelWithId(id, username, name, phoneNumber, email, password, hasAvatar),
	}
}

func (u ClientUserModel) ToUser() *User {
	var avatar string
	if u.HasAvatar {
		avatar = GetAvatarUrl(u.Id, ClientUserType)
	}
	return &User{
		Id:       u.Id,
		Name:     u.Name,
		Avatar:   avatar,
		UserType: ClientUserType,
	}
}

func (u ClientUserModel) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) (err error) {
	defer func() {
		const caller = packageCaller + "ClientUserModel.Delete"
		if err != nil {
			var uErr *utils.Error
			if errors.As(err, &uErr) {
				err = uErr.WithCaller(caller)
			} else {
				err = utils.NewInternalError(err).WithCaller(caller)
			}
		}
	}()
	if err_ := tx.WithContext(ctx).Model(u).Where("id = ?", id).Updates(
		map[string]any{
			"email":      nil,
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	).Error; err_ != nil {
		if errors.Is(err_, gorm.ErrRecordNotFound) {
			return utils.UserNotFoundResponse.Clone().
				WithReason("id", id.String())

		}
		return err_
	}
	return nil
}
