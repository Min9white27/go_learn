package web

import (
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/domain"
	"gitee.com/geektime-geekbang_admin/geektime-basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
)

//var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

// UserHandler 准备在这上面定义跟用户有关的路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	//用分组路由来简化注册，比较便利，不容易写错
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
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

	isEmail, err := h.emailExp.MatchString(req.Email)
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

	isPassword, err := h.passwordExp.MatchString(req.Password)
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
	err = h.svc.SignUp(ctx, domain.User{
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

//func (h *UserHandler) LoginJWT(ctx *gin.Context) {
//	type Req struct {
//		Email    string `json:"email"`
//		Password string `json:"password"`
//	}
//	var req Req
//	if err := ctx.Bind(&req); err != nil {
//		return
//	}
//	u, err := h.svc.Login(ctx, req.Email, req.Password)
//	switch err {
//	case nil:
//		uc := UserClaims{
//			Uid:       u.Id,
//			UserAgent: ctx.GetHeader("User-Agent"),
//			RegisteredClaims: jwt.RegisteredClaims{
//				// 1 分钟过期
//				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
//			},
//		}
//		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
//		tokenStr, err := token.SignedString(JWTKey)
//		if err != nil {
//			ctx.String(http.StatusOK, "系统错误")
//		}
//		ctx.Header("x-jwt-token", tokenStr)
//		ctx.String(http.StatusOK, "登录成功")
//	case service.ErrInvalidUserOrPassword:
//		ctx.String(http.StatusOK, "用户名或者密码不对")
//	default:
//		ctx.String(http.StatusOK, "系统错误")
//	}
//}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := h.svc.Login(ctx, req.Email, req.Password)
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
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	//type EditReq struct {
	//	Nickname        string `json:"nickname"`
	//	Birthday        string `json:"birthday"`
	//	PersonalProfile string `json:"personalProfile"`
	//}
	//
	//var req EditReq
	//if err := ctx.Bind(&req); err != nil {
	//	return
	//}
	//if req.Nickname == "" {
	//	ctx.String(http.StatusOK, "昵称不能为空")
	//	return
	//}
	//if len(req.PersonalProfile) > 1024 {
	//	ctx.String(http.StatusOK, "个人简介不能过长")
	//	return
	//}
	//
	//uc, ok := ctx.MustGet("user").(UserClaims)
	//if !ok {
	//	ctx.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//// DateOnly 可以将生日的格式转化成 “2006-01-02 ”，并返回
	//birthday, err := time.Parse(time.DateOnly, req.Birthday)
	//if err != nil {
	//	// 这里其实没有直接校验生日的具体格式，而是检查生日能够转化过来，就说明没有问题
	//	ctx.String(http.StatusOK, "生日格式不对")
	//	return
	//}
	//err = h.svc.UpdateNonSensitiveInfo(ctx, domain.User{
	//	Id:              uc.Uid,
	//	Nickname:        req.Nickname,
	//	Birthday:        birthday,
	//	PersonalProfile: req.PersonalProfile,
	//})
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统异常")
	//	return
	//}
	//ctx.String(http.StatusOK, "更新成功")

}

func (h *UserHandler) Profile(ctx *gin.Context) {

	//ctx.String(http.StatusOK, "这是你的 profile")

	//uc, ok := ctx.MustGet("user").(UserClaims)
	//if !ok {
	//	ctx.AbortWithStatus(http.StatusUnauthorized)
	//	return
	//}
	//u, err := h.svc.FindById(ctx, uc.Uid)
	//if err != nil {
	//	ctx.String(http.StatusOK, "系统异常")
	//	return
	//}
	//type ProfileReq struct {
	//	Email           string `json:"email"`
	//	Nickname        string `json:"nickname"`
	//	Birthday        string `json:"birthday"`
	//	PersonalProfile string `json:"personalProfile"`
	//}
	//ctx.JSON(http.StatusOK, ProfileReq{
	//	Email:           u.Email,
	//	Nickname:        u.Nickname,
	//	Birthday:        u.Birthday.Format(time.DateOnly),
	//	PersonalProfile: u.PersonalProfile,
	//})

}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
