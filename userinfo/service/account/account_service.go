package account

import (
	"context"
	"database/sql"
	"errors"
	errs "errs"
	"github.com/golang-jwt/jwt/v5"
	"loggers"
	"time"
	"user-server/dao"
	"user-server/model"
)

const (
	LOGIN_EXPIRE_TIME_HOUR = 24 * time.Hour
	JWT_SIGNING_KEY        = "ofB1pXMJFKs9N11yXomfPr1Vq0h5GE80"
)

type UserClaim struct {
	jwt.RegisteredClaims
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
}

type AccountService struct {
	userDao *dao.UserDao
	logger  *logger.Logger
}

func NewAccountService(userDao *dao.UserDao, logger *logger.Logger) *AccountService {
	return &AccountService{
		userDao: userDao,
		logger:  logger,
	}
}

func (s *AccountService) Register(ctx context.Context, email string, password string) error {
	s.logger.Info(ctx, "Call AccountService.Register, email: ", email)
	// 1. check whether email has been registered.
	user, err := s.userDao.GetUserByEmail(ctx, email)
	//   1.1 if user exists, return error.
	if user != nil && user.Id != 0 {
		s.logger.Error(ctx, "This email has already been registered.")
		return errs.New(errs.ERR_EMAIL_IS_REGISTERED)
	}
	//   1.2 if err is sql DB internal error, return error.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		s.logger.Error(ctx, "sql DB internal error, err: ", err.Error())
		return errs.New(errs.ERR_REGISTER_INTERNAL)
	}

	// 2. save email and password.
	user = &model.User{
		Password: password,
		Email:    email,
	}
	err = s.userDao.Insert(ctx, user)
	if err != nil {
		s.logger.Error(ctx, "Insert user failed, err: ", err.Error())
		return errs.New(errs.ERR_REGISTER_INTERNAL)
	}
	return nil
}

func (s *AccountService) Login(ctx context.Context, email string, password string) (*string, error) {
	s.logger.Info(ctx, "Call AccountService.Login, email: ", email)
	// 1. verify email and password.
	user, err := s.userDao.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(ctx, "Get user failed, err: ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.New(errs.ERR_LOGIN_NO_USER)
		} else {
			return nil, errs.New(errs.ERR_LOGIN_INTERNAL)
		}
	}
	if user.Password != password {
		s.logger.Error(ctx, "Password mismatch.")
		return nil, errs.New(errs.ERR_PASSWORD_MISMATCH)
	}
	// 2. login succeed, return token.
	claim := &UserClaim{}
	claim.UserId = user.Id
	claim.Email = email
	claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(LOGIN_EXPIRE_TIME_HOUR))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(JWT_SIGNING_KEY))
	if err != nil {
		s.logger.Error(ctx, "Generate token failed, err: ", err.Error())
		return nil, errs.New(errs.ERR_LOGIN_INTERNAL)
	}
	s.logger.Info(ctx, "Call AccountService.Login succeed.")
	return &tokenString, nil
}

func (s *AccountService) Authenticate(ctx context.Context, token string) (uint64, string, error) {
	s.logger.Info(ctx, "Call AccountService.Authenticate, token: ", token)
	claim := &UserClaim{}
	tk, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SIGNING_KEY), nil
	})
	if err != nil {
		s.logger.Error(ctx, "Parse token failed, err: ", err.Error())
		return 0, "", errs.New(errs.ERR_AUTH_FAILED)
	}
	if tk == nil || !tk.Valid || claim.UserId == 0 {
		s.logger.Error(ctx, "Token invalid.")
		return 0, "", errs.New(errs.ERR_AUTH_FAILED)
	}

	// check whether expired.
	if time.Now().After(claim.ExpiresAt.Time) {
		s.logger.Error(ctx, "Token expired.")
		return 0, "", errs.New(errs.ERR_TOKEN_EXPIRED)
	}
	s.logger.Info(ctx, "Call AccountService.Authenticate succeed.")
	return claim.UserId, claim.Email, nil
}
