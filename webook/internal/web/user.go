package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web/jwt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const biz = "login"

// 确保 UserHandler 上实现了 handler 接口
var _ handler = &UserHandler{}

// 这个更优雅
var _ handler = (*UserHandler)(nil)

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

// UserHandler 准备在这上面定义跟用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
	ijwt.Handler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService, jwtHdl ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		phoneRegexPattern    = `^1([38][0-9]|4[^23]|[59][^4]|6[2567]|7[0-8])\d{8}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		phoneExp:    phoneExp,
		Handler:     jwtHdl,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	//用分组路由来简化注册，比较便利，不容易写错
	//ug.GET("/profile", u.Profile)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/signup", u.SignUp)
	//ug.POST("/login", u.Login)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.LogoutJWT)
	ug.POST("/edit", u.Edit)
	// PUT "/login/sms/code" 发送验证码
	// POST “/login/sms/code” 校验验证码
	// /sms/login/code /code/sms
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}

// RefreshToken 可以同时刷新长短 token，用 redis 来记录是否有效，即 refresh_token 是一次性的
// 参考登录校验部分，比较 User-Agent 来增强安全性
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	// 只有这个接口，拿出来的才是 refresh_token,其他地方都是 access_token
	refreshToken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		// 要么 redis 有问题，要么已经退出登录
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	// 搞个新的 access_token
	err = u.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		// 这种含糊的日志，尽量少打，信息量不足
		zap.L().Error("系统异常", zap.Error(err))
		// 正常来说，msg 的部分就应该包含足够的定位信息
		zap.L().Error("ijoihpidf 设置 JWT token 出现异常", zap.Error(err),
			zap.String("method", "UserHandler:RefreshToken"))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "刷新成功",
	})
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	// Bind 方法会根据 content-Type 来解析你的数据到 req 里面
	//解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入得密码不一致")
		return
	}

	isPassword, err := u.passwordExp.MatchString(req.Password)
	if err != nil {
		//记录日志
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码长度必须大于8位，包含数字、特殊字符")
		return
	}
	//调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
	//这边属于数据库操作
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不正确")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	// 在这里用 JWT 设置登录态，即生成一个 JWT token

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "登录成功")
	return
}
func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登录失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登录",
	})
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码不正确")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	// 在这里登录成功了
	// 设置 session
	sess := sessions.Default(ctx)
	// 可以随便设置值，就是要放在 session 里面的值
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		Secure:   true,
		HttpOnly: true,
		MaxAge:   60,
	})
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Options(sessions.Options{
		//Secure: true,
		//HttpOnly: true,
		MaxAge: -1,
	})
	sess.Save()
	ctx.String(http.StatusOK, "退出登录成功")
}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 用正则表达式校验手机号码输入正确与否
	isPhone, err := u.phoneExp.MatchString(req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !isPhone {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号码输入错误",
		})
		return
	}
	err = u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		zap.L().Warn("短信发送太频繁", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		zap.L().Error("短信发送失败", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
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
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("校验验证码出错", zap.Error(err),
			// 不能这样打，因为手机号码是敏感数据，你不能达到日志里面
			// 打印加密后的串
			// 脱敏，152****1234
			zap.String("手机号码", req.Phone))
		// 最多最多就这样
		zap.L().Debug("", zap.String("手机号码", req.Phone))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}

	// FindOrCreate 要想到手机号码并不一定能找得到，因为这个手机号码可能还未进行注册，所以取名 FindOrCreate
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	// 设置 JWTToken 进行登陆
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		// 记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		Nickname        string `json:"nickname"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"personalProfile"`
	}

	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	sess := sessions.Default(ctx)
	uid, ok := sess.Get("userId").(int64)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if uid == 0 {
		//	没有登录
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	if req.Nickname == "" {
		ctx.String(http.StatusOK, "昵称不能为空")
		return
	}
	if len(req.PersonalProfile) > 1024 {
		ctx.String(http.StatusOK, "个人简介不能过长")
		return
	}
	//
	//uc, ok := ctx.MustGet("user").(UserClaims)
	//if !ok {
	//	ctx.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	// DateOnly 可以将生日的格式转化成 “2006-01-02 ”，并返回
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		// 这里其实没有直接校验生日的具体格式，而是检查生日能够转化过来，就说明没有问题
		ctx.String(http.StatusOK, "生日格式不对")
		return
	}
	err = u.svc.UpdateNonSensitiveInfo(ctx, domain.User{
		Id:              uid,
		Nickname:        req.Nickname,
		Birthday:        birthday,
		PersonalProfile: req.PersonalProfile,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "更新成功")

}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	//ctx.String(http.StatusOK, "这是你的 profile")
	c, _ := ctx.Get("claims")
	// 可以断定，必然有 claims
	//if !ok {
	//	// 在这里设置个监控，用来判定有没有拿到 claims
	//	ctx.String(http.StatusOK, "系统错误")
	//	return
	//}
	// 通过类型断言，用 ok 来判断是不是
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 在这里设置个监控，用来判定有没有拿到 claims
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	println(claims.Uid)
	ctx.String(http.StatusOK, "你的 profile")
}

func (u *UserHandler) Profile(ctx *gin.Context) {

	//ctx.String(http.StatusOK, "这是你的 profile")

	type ProfileReq struct {
		Nickname        string `json:"nickname"`
		Birthday        string `json:"birthday"`
		PersonalProfile string `json:"personalProfile"`
	}

	var req ProfileReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	sess := sessions.Default(ctx)
	uid := sess.Get("userId")

	if uid == nil {
		//	没有登录
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	//} else {
	//	sess.Set("userId", uid)
	//}

	uc, err := u.svc.FindById(ctx, uid)

	//uc, ok := ctx.MustGet("user").(UserClaims)
	//if !ok {
	//	ctx.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//u, err := u.svc.FindById(ctx, uc.Uid)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	//type ProfileReq struct {
	//	Email           string `json:"email"`
	//	Nickname        string `json:"nickname"`
	//	Birthday        string `json:"birthday"`
	//	PersonalProfile string `json:"personalProfile"`
	//}
	ctx.JSON(http.StatusOK, ProfileReq{
		Nickname:        uc.Nickname,
		Birthday:        uc.Birthday.Format(time.DateOnly),
		PersonalProfile: uc.PersonalProfile,
	})

}

//func (u *UserHandler) Do(fn func(ctx *gin.Context) (any, error)) {
//	data, err := fn(ctx)
//	if err != nil {
// 在这里打日志
//	}
//}
