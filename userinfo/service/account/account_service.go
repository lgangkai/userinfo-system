package account

import (
	"database/sql"
	"errors"
	errs "errs"
	"github.com/asim/go-micro/v3/logger"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"user-server/dao"
	"user-server/model"
)

const (
	LOGIN_EXPIRE_TIME_HOUR = 24 * time.Hour
	JWT_CLAIM_KEY_ID       = "user_id"
	JWT_SIGNING_KEY        = "ofB1pXMJFKs9N11yXomfPr1Vq0h5GE80"
)

type UserClaim struct {
	jwt.RegisteredClaims
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
}

type AccountService struct {
	userDao *dao.UserDao
}

func NewAccountService(userDao *dao.UserDao) *AccountService {
	return &AccountService{
		userDao: userDao,
	}
}

func (s *AccountService) Register(email string, password string) error {
	logger.Info("Call AccountService.Register, email: ", email)
	// 1. check whether email has been registered.
	user, err := s.userDao.GetUserByEmail(email)
	//   1.1 if user exists, return error.
	if user != nil && user.Id != 0 {
		logger.Error("This email has already been registered.")
		return errs.New(errs.ERR_EMAIL_IS_REGISTERED)
	}
	//   1.2 if err is sql DB internal error, return error.
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("sql DB internal error, err: ", err.Error())
		return errs.New(errs.ERR_REGISTER_GET_USER_FAILED)
	}

	// 2. save email and password.
	user = &model.User{
		Password: password,
		Email:    email,
	}
	err = s.userDao.Insert(user)
	if err != nil {
		logger.Error("Insert user failed, err: ", err.Error())
		return errs.New(errs.ERR_REGISTER_INTERNAL)
	}
	return nil
}

func (s *AccountService) Login(email string, password string) (*string, error) {
	logger.Info("Call AccountService.Login, email: ", email)
	// 1. verify email and password.
	user, err := s.userDao.GetUserByEmail(email)
	if err != nil {
		logger.Error("Get user failed, err: ", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.New(errs.ERR_LOGIN_NO_USER)
		} else {
			return nil, errs.New(errs.ERR_LOGIN_INTERNAL)
		}
	}
	if user.Password != password {
		logger.Error("Password mismatch.")
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
		logger.Error("Generate token failed, err: ", err.Error())
		return nil, errs.New(errs.ERR_LOGIN_INTERNAL)
	}
	logger.Info("Call AccountService.Login succeed.")
	return &tokenString, nil
}

func (s *AccountService) Authenticate(token string) (uint64, string, error) {
	logger.Info("Call AccountService.Authenticate, token: ", token)
	claim := &UserClaim{}
	tk, err := jwt.ParseWithClaims(token, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWT_SIGNING_KEY), nil
	})
	if err != nil {
		logger.Error("Parse token failed, err: ", err.Error())
		return 0, "", errs.New(errs.ERR_AUTH_FAILED)
	}
	if tk == nil || !tk.Valid || claim.UserId == 0 {
		logger.Error("Token invalid.")
		return 0, "", errs.New(errs.ERR_AUTH_FAILED)
	}

	// check whether expired.
	if time.Now().After(claim.ExpiresAt.Time) {
		logger.Error("Token expired.")
		return 0, "", errs.New(errs.ERR_TOKEN_EXPIRED)
	}
	logger.Info("Call AccountService.Authenticate succeed.")
	return claim.UserId, claim.Email, nil
}
