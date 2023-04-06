package custom_v0

// Middleware 函数式的责任链模式
// 函数式的洋葱模式
type Middleware func(next HandleFunc) HandleFunc

// AOP 方案在不同的框架，不同的语言里面都有不同的叫法
// Middleware, Handler, Chain, Filter, Filter-Chain
// Interceptor, Wrapper

// type MiddlewareV1 interface {
// 	Invoke(next HandleFunc) HandleFunc
// }
//
// type Interceptor interface {
// 	Before(ctx *Context)
// 	After(ctx *Context)
// 	Surround(ctx *Context)
// }

// type Chain []HandleFuncV1
//
// type HandleFuncV1 func(ctx *Context) (next bool)
//
// type ChainV1 struct {
// 	handlers []HandleFuncV1
// }
//
// func (c ChainV1) Run(ctx *Context) {
// 	for _, h := range c.handlers {
// 		next := h(ctx)
// 		// 这种是中断执行
// 		if !next {
// 			return
// 		}
// 	}
// }

//
// type Net struct {
// 	handlers []HandleFuncV1
// }
//
// func (c Net) Run(ctx *Context) {
// 	var wg sync.WaitGroup
// 	for _, hdl := range c.handlers {
// 		h := hdl
// 		if h.concurrent {
// 			wg.Add(1)
// 			go func() {
// 				h.Run(ctx)
// 				wg.Done()
// 			}()
// 		} else {
// 			h.Run(ctx)
// 		}
// 	}
// 	wg.Wait()
// }

// type HandleFuncV1 struct {
// 	concurrent bool
// 	handlers []*HandleFuncV1
// }

// func (HandleFuncV1) Run(ctx *Context) {
// for _, hdl := range c.handlers {
// 	h := hdl
// 	if h.concurrent {
// 		wg.Add(1)
// 		go func() {
// 			h.Run(ctx)
// 			wg.Done()
// 		}()
// 	}
// }
// }
