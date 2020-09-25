### go 原生支持rpc 调用

1、net/rpc/server.go 

   注册函数: Register
   
   HTTP处理器: HandleHTTP
 
2、net/rpc/client.go

   连接函数: DialHTTP
   
   同步呼叫函数: Call
   
   异步呼叫函数: Go
 
3、服务端代码如下

具体实现业务逻辑代码如下:

```
// 定义rpc的请求结构体
type StringRequest struct {
	SA string
	SB string
}
// 定义rpc接口
type Service interface {
	// Concat sa and sb
	Concat(req StringRequest, ret *string) error
}
type StringService struct {
}
// 实现Service接口
func (stringService *StringService) Concat(req StringRequest, ret *string) error {
	// len(sa + sb) > StrMaxSize to ErrMaxSize
	if len(req.SA)+len(req.SB) > StrMaxSize {
		*ret = ""
		return ErrMaxSize
	}
	time.Sleep(time.Second * 2)
	*ret = req.SA + req.SB
	return nil
}
``` 

服务端启动代码如下:

```
func main() {
	// 构建实现结构体
	stringService := new(service.StringService)
	err := rpc.Register(stringService)
	if err != nil {
		log.Printf("register rpc err : %s\n", err)
	}
	rpc.HandleHTTP()

	listen, err := net.Listen("tcp", "localhost:9527")
	if err != nil {
		log.Printf("listen address err : %s\n", err)
	}
	http.Serve(listen, nil)
}
```

4、客户端代码如下

客户端调用服务端代码如下:
```
func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:9527")
	if err != nil {
		log.Printf("client listen tcp err : %s\n", err)
	}

	sr := &service.StringRequest{"hello ", "go rpc"}

    // 同步调用
	var reply string
	err = client.Call("StringService.Concat", sr, &reply)
	if err != nil {
		log.Printf("client call remote method error : %s\n", err)
	}
	fmt.Println(" -------------------- ", reply)

	// 异步调用
	var asynReply string
	future := make(chan *rpc.Call, 1)
	completeableFuture := client.Go("StringService.Concat", sr, &asynReply, future)
	fmt.Println(" ------------------- ", asynReply)
	_ = <-completeableFuture.Done
	fmt.Println(" ------------------- ", asynReply)
}
```

### Go RPC 服务端原理

服务端的 RPC 代码主要分为两个部分：

    1、服务方法注册，包括调用注册接口，通过反射处理将方法取出，并存到 map 中；
    2、处理网络调用，主要是监听端口、读取数据包、解码请求和调用反射处理后的方法，将返回值编码，返回给客户端。
      处理流程如下: 
        service.StringService (创建一个StringService结构体)    
                 |
          注册结构体实现的接口方法     
                 |
              反射处理
                 |
              Stub保存    
                    
register函数解析:

```
func (server *Server) register(rcvr interface{}, name string, useName bool) error {
  // 如果服务为空，默认注册一个 
  if server.serviceMap == nil { 
    server.serviceMap = make(map[string]*service) 
  }
  // 获取注册服务的反射信息 
  s := new(service) 
  s.typ = reflect.TypeOf(rcvr) 
  s.rcvr = reflect.ValueOf(rcvr) 
  // 可以使用自定义名称 
  sname := reflect.Indirect(s.rcvr).Type().Name() 
  if useName { 
    sname = name 
  } 
  // 方法必须是暴露的，既服务名首字符大写；不允许重复注册。代码有省略
  if !isExported(sname) && !useName { 
  }   
  if _, present := server.serviceMap[sname]; present { 
  } 
  
  s.name = sname 
  // 开始注册 rpc struct 内部的方法存根 
  s.method = suitableMethods(s.typ, true) 
  if len(s.method) == 0 { 
     // 如果struct内部一个方法也没，那么直接报错，打印详细的错误信息
  }
  // 保存在server的serviceMap中 
  server.serviceMap[s.name] = s 
  return nil 
} 
```

Server 处理请求流程:

```
1、接受请求 - net/rpc/server.go

func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Print("rpc.Serve: accept:", err.Error())
			return
		}
        // 没接受一个客户端的连接就启动一个gorutine去处理
		go server.ServeConn(conn)
	}
}

2、读取并解析请求 - net/rpc/server.go

func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	buf := bufio.NewWriter(conn)
	srv := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewEncoder(buf),
		encBuf: buf,
	}
    // 按不同的协议编解码
	server.ServeCodec(srv)
}

func (server *Server) ServeCodec(codec ServerCodec) { 
  sending := new(sync.Mutex) 
  for { 
    // 解析请求 
    service, mtype, req, argv, replyv, keepReading, err := server.readRequest (codec)
    if err != nil { 
      if debugLog && err != io.EOF { 
        log.Println("rpc:", err) 
      } 
      if !keepReading { 
        break 
      } 
      // send a response if we actually managed to read a header. 
      // 如果当前请求错误了，我们应该返回信息，然后继续处理 
      if req != nil { 
        server.sendResponse(sending, req, invalidRequest, codec, err.Error())
        server.freeRequest(req) 
      } 
      continue 
    } 
    // 因为需要继续处理后续请求，所以开一个gorutine处理rpc方法 
    go service.call(server, sending, mtype, req, argv, replyv, codec) 
  }
  // 如果连接关闭了需要释放资源 
  codec.Close() 
} 

3、执行远程方法并返回响应 - net/rpc/server.go

func (s *service) call(server *Server, sending *sync.Mutex, mtype *methodType, req *Request, argv, replyv reflect.Value, codec ServerCodec) {
  function := mtype.method.Func 
  // 这里是真正调用rpc方法的地方 
  returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv}) 
  errInter := returnValues[0].Interface() 
  errmsg := "" 
  // 处理返回请求了 
  server.sendResponse(sending, req, replyv.Interface(), codec, errmsg) 
  server.freeRequest(req) 
} 

```

### Go RPC 客户端原理

客户端发送 RPC 请求原理

```
1、同步调用 - net/rpc/client.go

func (client *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
 // 同步 直接调用了 Go 方法
 call := <-client.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
 return call.Error
}

2、异步调用 - net/rpc/client.go

先创建并初始化了 Call 对象，记录下此次调用的方法、参数和返回值，并生成 DoneChannel；
然后调用 Client 的 send 方法进行真正的请求发送处理  

// 异步调用实现
func (client *Client) Go(serviceMethod string, args interface{}, reply interface{}, done chan *Call) *Call {
 // 初始化 Call 
 call := new(Call)
 call.ServiceMethod = serviceMethod
 call.Args = args
 call.Reply = reply
 if done == nil {
   done = make(chan *Call, 10) // buffered.
 } else {
   if cap(done) == 0 {
     log.Panic("rpc: done channel is unbuffered")
   }
 }
 call.Done = done
  // 调用 Client 的 send 方法
 client.send(call)
 return call
}
type Call struct {
  ServiceMethod string   // 服务名及方法名 格式:服务.方法
  Args     interface{} // 函数的请求参数 (*struct).
  Reply     interface{} // 函数的响应参数 (*struct).
  Error     error    // 方法完成后 error的状态.
  Done     chan *Call // 
}

3、请求参数编码 - net/rpc/client.go
Client 的 send 函数首先会判断客户端实例的状态，如果处于关闭状态，则直接返回结果;
否则会生成唯一的 seq 值，将 Call 保存到客户端的哈希表 pending 中，然后调用客户端编码器的WriteRequest 来编码请求并发送;

func (client *Client) send(call *Call) {
  // 省略锁等操作 ...
  //生成seq,每次调用均生成唯一的seq,在服务端返回结果后会通过该值进行匹配
  seq := client.seq
  client.seq++
  client.pending[seq] = call
 
  // 请求并发送请求
  client.request.Seq = seq
  client.request.ServiceMethod = call.ServiceMethod
  err := client.codec.WriteRequest(&client.request, call.Args)
  if err != nil {
    //发送请求错误时,将map中call对象删除.
    client.mutex.Lock()
    call = client.pending[seq]
    delete(client.pending, seq)
    client.mutex.Unlock()
  }
}

4、接收返回值 - net/rpc/client.go

客户端的 input 函数接收服务端返回的响应值，它进行无限 for 循环，不断调用 codec 也就是 gobClientCodecd 的 ReadResponseHeader 函数，
然后根据其返回数据中的 seq 来判断是否是本客户端发出请求的响应值。
如果是，则获取对应的 Call 对象，并将其从 pending 哈希表中删除，继续调用 codec 的 ReadReponseBody 方法获取返回值 Reply 对象，并调用 Call 对象的 done 方法

func (client *Client) input() {
  var err error
  var response Response
  for err == nil {
    response = Response{}
    //通过response中的 Seq获取call对象
    seq := response.Seq
    client.mutex.Lock()
    call := client.pending[seq]
    delete(client.pending, seq)
    client.mutex.Unlock()
 
    switch {
    case call == nil:
    case response.Error != "":
      //上述两个case，一个处理call为nil，另外处理服务端返回的错误,直接将错误返回
    default:
      //通过编码器,将Resonse的body部分解码成reply对象.
      err = client.codec.ReadResponseBody(call.Reply)
      if err != nil {
        call.Error = errors.New("reading body " + err.Error())
      }
      call.done()
    }
  }
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
		if debugLog {
			log.Println("rpc: discarding Call reply due to insufficient Done chan capacity")
		}
	}
}

```




