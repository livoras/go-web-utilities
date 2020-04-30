package gw

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var JWT_SECRECT = []byte("JWT")

type Context struct {
	c *gin.Context
	updateFields map[string]bool
}

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SetJwtSecrect(sc string) {
	JWT_SECRECT = []byte(sc)
}

func (c *Context) PathInt (key string) int {
	ret, err := strconv.Atoi(c.c.Param(key))
	if err != nil {
		return 0
	}
	return ret
}

func (c *Context) PathUint (key string) uint {
	return uint(c.PathInt(key))
}

func (c *Context) OK(obj interface{}) {
	c.c.JSON(http.StatusOK, Result{
		Message: "OK",
		Data:    obj,
	})
}

func (c *Context) BadRequest(message string) {
	if message == "" { message = "客户端参数错误" }
	c.c.JSON(http.StatusBadRequest, Result{
		Message: message,
		Data:    nil,
	})
}

func (c *Context) Unauthorized(message string) {
	if message == "" { message = "未授权操作" }
	c.c.JSON(http.StatusUnauthorized, Result{
		Message: message,
		Data:    nil,
	})
}

func (c *Context) NotFound(message string) {
	if message == "" { message = "找不到该数据" }
	c.c.JSON(http.StatusNotFound, Result{
		Message: message,
		Data:    nil,
	})
}

func (c *Context) IntervalServerError(message string) {
	if message == "" { message = "服务器异常" }
	c.c.JSON(http.StatusInternalServerError, Result{
		Message: message,
		Data:    nil,
	})
}

func (c *Context) BindJSON(v interface{}) error {
	return c.c.BindJSON(v)
}

//func (c *Context) SetJWTByUser(user *models.User) error {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ "userId": user.ID })
//	tokenString, err := token.SignedString(JWT_SECRECT)
//	if err != nil { return err }
//	c.c.SetCookie("token", tokenString, 14 * 24 * 60 * 60, "/", "", false, false)
//	return nil
//}

func (c *Context) RemoveToken () {
	c.c.SetCookie("token", "", -1, "/", "", false, false)
}

//func (c *Context) GetUserByJWT(tokenString string) (*models.User, error) {
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//			return nil, errors.New("token 方法不对")
//		}
//		return JWT_SECRECT, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	userIdStr := token.Claims.(jwt.MapClaims)["userId"]
//	if userIdStr == nil {
//		return nil, errors.New("登录态无效")
//	}
//	userValue := reflect.ValueOf(userIdStr)
//	if userValue.Kind() != reflect.Float64 {
//		return nil, errors.New("登录态无效")
//	}
//	userId := int(userValue.Float())
//	user, err := models.FindUserG(userId)
//	if err != nil || user == nil {
//		return nil, errors.New("当前用户不存在")
//	}
//	return user, err
//}

func (c *Context) HasField (field string) bool {
	if _, ok := c.updateFields[field]; ok {
		return true
	}
	return false
}

func (c *Context) MergeUpdate (dest interface{}, src interface{})  {
	valOfDest := reflect.ValueOf(dest).Elem()
	valOfSrc := reflect.ValueOf(src)
	typeOfDest := valOfDest.Type()
	numOfFields := typeOfDest.NumField()
	for i := 0; i < numOfFields; i++ {
		field := typeOfDest.Field(i)
		fieldName := field.Name
		jsonName := field.Tag.Get("json")
		// 如果 ?fields=xxx,yyy,zzz 中没有对应的值，就不更新
		if !c.HasField(jsonName) { continue }
		valOfDest.FieldByName(fieldName).Set(
			valOfSrc.FieldByName(fieldName),
		)
	}
}

func NewContext(c *gin.Context) *Context {
	ctx := &Context{ c: c }
	ctx.updateFields = make(map[string]bool)
	fieldsStr := c.Query("fields")
	if fieldsStr != "" {
		fields := strings.Split(fieldsStr, ",")
		for _, field := range fields {
			ctx.updateFields[field] = true
		}
	}
	return ctx
}

func Rest(f func(c *Context)) gin.HandlerFunc {
	return func(context *gin.Context) {
		f(NewContext(context))
	}
}

//func RequireLogin(f func(c *Context, user *models.User)) gin.HandlerFunc {
//	return func(context *gin.Context) {
//		c := NewContext(context)
//		token, err := context.Cookie("token")
//		if err != nil {
//			c.Unauthorized("请先登录")
//			return
//		}
//		user, err := c.GetUserByJWT(token)
//		if err != nil {
//			c.Unauthorized(err.Error())
//			return
//		}
//		f(c, user)
//	}
//}

/*func RequireAdmin(f func(c *Context, admin *models.User)) gin.HandlerFunc {
	return func(context *gin.Context) {
		f(NewContext(context), nil)
	}
}
*/

