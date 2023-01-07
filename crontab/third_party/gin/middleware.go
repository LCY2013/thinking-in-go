package _gin

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"reflect"
)

var (
	ctxType = reflect.TypeOf(new(context.Context)).Elem()
	errType = reflect.TypeOf(new(error)).Elem()
)

type HttpResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Wrapper 封装具体的函数实现
func Wrapper(function any) gin.HandlerFunc {
	ftype := reflect.TypeOf(function)
	fv := reflect.ValueOf(function)

	return func(ctx *gin.Context) {
		if ftype == nil {
			log.WithFields(log.Fields{
				"gin-Wrapper": "err",
			}).Error(errors.New("can't invoke an untyped nil"))
			return
		}
		if ftype.Kind() != reflect.Func {
			log.WithFields(log.Fields{
				"gin-Wrapper": "err",
			}).Errorf("can't invoke non-function %v (type %v)", function, ftype)
			return
		}

		fnInNum, fnOutNum := ftype.NumIn(), ftype.NumOut()
		argv := parseArgs(ctx, ftype, fnInNum)

		fnRet := fv.Call(argv)

		for i := 0; i < fnOutNum; i++ {
			idxRet := fnRet[i].Interface()
			if fnRet[i].Type().Implements(errType) && idxRet != nil {
				err := idxRet.(error)
				ctx.JSON(http.StatusOK, &HttpResponse{
					Code: -1,
					Msg:  err.Error(),
				})
				goto OUT
			}
			if idxRet == nil {
				ctx.JSON(http.StatusOK, &HttpResponse{
					Code: 0,
					Data: struct{}{},
				})
				goto OUT
			}
			ctx.JSON(http.StatusOK, &HttpResponse{
				Code: 0,
				Data: idxRet,
			})
		}
	OUT:
		if fnOutNum == 0 {
			ctx.JSON(http.StatusOK, &HttpResponse{
				Code: 0,
				Data: struct{}{},
			})
		}
	}
}

var (
	typeParses = map[reflect.Kind]func(ctx *gin.Context, rtype reflect.Type) reflect.Value{
		reflect.Struct: typeParsesStruct,
		reflect.Map:    typeParsesMap,
		reflect.String: typeParsesString,
	}
)

func typeParsesString(c *gin.Context, rtype reflect.Type) reflect.Value {
	input := reflect.New(rtype).Interface()
	_ = c.ShouldBindQuery(input)
	return reflect.ValueOf(input)
}

func typeParsesMap(c *gin.Context, rtype reflect.Type) reflect.Value {
	input := reflect.New(rtype).Interface()
	if c.ContentType() == gin.MIMEJSON {
		_ = c.ShouldBindJSON(input)
	} else {
		_ = c.ShouldBind(input)
	}
	return reflect.ValueOf(input)
}

func typeParsesStruct(c *gin.Context, rtype reflect.Type) reflect.Value {
	input := reflect.New(rtype).Interface()
	if c.ContentType() == gin.MIMEJSON {
		_ = c.ShouldBindJSON(input)
	} else {
		_ = c.ShouldBind(input)
	}
	return reflect.ValueOf(input)
}

// parseArgs 解析请求参数信息
func parseArgs(ctx *gin.Context, ftype reflect.Type, fnInNum int) []reflect.Value {
	argv := make([]reflect.Value, fnInNum)
	for i := 0; i < fnInNum; i++ {
		inType := ftype.In(i).Elem()

		// 参数是ctx时
		if inType.Implements(ctxType) {
			argv[i] = reflect.ValueOf(ctx)
			continue
		}

		// 其他参数类型
		if typeParse, ok := typeParses[inType.Kind()]; ok {
			argv[i] = typeParse(ctx, inType)
			continue
		}

		input := reflect.New(inType).Interface()
		_ = ctx.ShouldBindQuery(input)

		argv[i] = reflect.ValueOf(input)
	}
	return argv
}
