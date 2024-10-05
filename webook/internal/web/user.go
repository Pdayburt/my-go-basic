package web

import (
	"errors"
	"example.com/mod/webook/internal/domain"
	"example.com/mod/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const (
	biz = "login"
)

type UserHandler struct {
	svc           service.UserService
	codeSvc       service.CodeService
	emailRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailRegex := regexp.MustCompile(emailRegexPattern, regexp.Debug)
	passwordRegex := regexp.MustCompile(passwordRegexPattern, regexp.Debug)
	return &UserHandler{
		svc:           svc,
		codeSvc:       codeSvc,
		emailRegex:    emailRegex,
		passwordRegex: passwordRegex,
	}
}

func (u *UserHandler) RegisterRoutes(ginServer *gin.Engine) {

	group := ginServer.Group("/users")

	group.POST("/signup", u.SignUp)

	//group.POST("/login", u.Login)
	group.POST("/login", u.LoginJWT)

	group.POST("/edit", u.Edit)

	group.GET("/profile", u.Profile)
	group.GET("/profileJWT", u.ProfileJWT)
	group.POST("/login_sms/code/send", u.SendLoginSMSCode)
	group.POST("/login_sms", u.LoginSMS)

}
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 2,
		Msg:  "登陆成功",
	})
	return

}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	err := u.codeSvc.Send(ctx, biz, req.Phone)

	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 0,
			Msg:  "发送成功",
		})
		return
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email" `
		Password        string `json:"password" `
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	matched, err := u.emailRegex.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统内部错误")
		return
	}
	if !matched {
		ctx.String(http.StatusOK, "邮件格式错误")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}

	match, err := u.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统内部错误")
		return
	}
	if !match {
		ctx.String(http.StatusOK, "密码必须大于8,并包含特殊字符")
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱已被注册")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常！！")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%#v\n", req)

}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	//{"email":"11","password":"11"}
	type LoginParam struct {
		Email    string `json:"email" `
		Password string `json:"password" `
	}
	var param LoginParam
	if err := ctx.Bind(&param); err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	user, err := u.svc.Login(ctx, domain.User{
		Email:    param.Email,
		Password: param.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统内部错！")
		return
	}
	/*session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Save()*/
	err = u.setJWTToken(ctx, user.Id)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
	}
	//fmt.Printf("tokensignedString %#v\n", signedString)
	fmt.Printf("user :%#v\n", user)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, userId int64) error {
	claims := UserClaims{
		Uid:       userId,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.
		SignedString([]byte("30pzPuWsJCXJi5eryywAYltH5AS4GcOAA7aBwkDGpu0vGSqnVxjFLOmlLLNaWNsF"))
	ctx.Header("x-jwt-token", signedString)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return err
	}
	return err
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func (u *UserHandler) Login(ctx *gin.Context) {
	//{"email":"11","password":"11"}
	type LoginParam struct {
		Email    string `json:"email" `
		Password string `json:"password" `
	}
	var param LoginParam
	if err := ctx.Bind(&param); err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    param.Email,
		Password: param.Password,
	})

	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或者密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统内部错！")
		return
	}

	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	value, ok := ctx.Get("userId")
	if !ok {
		//监控这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	val, ok := value.(int64)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Printf("val %#v\n", val)
	ctx.String(http.StatusOK, "这是你的profile")
}
func (u *UserHandler) Profile(ctx *gin.Context) {

	userId, ok := ctx.Get("userId")
	if !ok {
		//监控这里
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	id := userId.(int64)
	profile, err := u.svc.Profile(ctx, id)
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
