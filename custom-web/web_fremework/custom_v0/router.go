package custom_v0

import (
	"fmt"
	"strings"
)

// 用来支持对路由树的操作
// 代表路由树（森林）
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// GET 有一棵树，POST 也有一棵树

	// http method => 路由树根节点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{},
	}
}

// 加一些限制：
// path 必须以 / 开头，不能以 / 结尾，中间也不能有连续的 //
func (r *router) addRoute(method string, path string, handleFunc HandleFunc) {
	if path == "" {
		panic("path cannot be empty")
	}

	// find method tree
	methodRoot, ok := r.trees[method]

	if !ok {
		// not found method methodRoot, create
		methodRoot = &node{
			path: "/",
		}

		r.trees[method] = methodRoot
	}

	// The beginning cannot be without /
	if path[0] != '/' {
		panic("the beginning cannot be without /")
	}

	// cannot end with /
	if path != "/" && path[len(path)-1] == '/' {
		panic("cannot end with /")
	}

	// Multiple consecutive /s cannot exist in the path
	if strings.Contains(path, "//") {
		panic("multiple consecutive /s cannot exist in the path")
	}

	// methodRoot path
	// special treatment for the root node
	if path == "/" {
		// root node duplicate registration
		if methodRoot.handler != nil {
			panic("web: routing conflict, duplicate registration[/]")
		}
		methodRoot.handler = handleFunc
		methodRoot.route = "/"
		return
	}

	// /user/home 被切割成三段
	// 切割这个 path
	segs := strings.Split(path[1:], "/")
	for _, seg := range segs {
		if seg == "" {
			panic("web: 不能有连续的 /")
		}
		// 递归下去，找准位置
		// 如果中途有节点不存在，你就要创建出来
		child := methodRoot.childOrCreate(seg)
		methodRoot = child
	}
	if methodRoot.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突，重复注册[%s]", path))
	}
	methodRoot.handler = handleFunc
	methodRoot.route = path
}

func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	// 基本上是不是也是沿着树深度查找下去？
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}

	// 这里把前置和后置的 / 都去掉
	path = strings.Trim(path, "/")

	// 按照斜杠切割
	segs := strings.Split(path, "/")
	var pathParams map[string]string
	for _, seg := range segs {
		child, paramChild, found := root.childOf(seg)
		if !found {
			return nil, false
		}
		// 命中了路径参数
		if paramChild {
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			// path 是 :id 这种形式
			pathParams[child.path[1:]] = seg
		}
		root = child
	}
	// 代表我确实有这个节点
	// 但是节点是不是用户注册的有 handler 的，就不一定了
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true

	// return root, root.handler != nil
}

func (n *node) childOrCreate(seg string) *node {

	if seg[0] == ':' {
		if n.starChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有通配符匹配")
		}
		n.paramChild = &node{
			path: seg,
		}
		return n.paramChild
	}

	if seg == "*" {
		if n.paramChild != nil {
			panic("web: 不允许同时注册路径参数和通配符匹配，已有路径参数")
		}
		n.starChild = &node{
			path: seg,
		}
		return n.starChild
	}

	if n.children == nil {
		n.children = map[string]*node{}
	}
	res, ok := n.children[seg]
	if !ok {
		// 要新建一个
		res = &node{
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

// childOf 优先考虑静态匹配，匹配不上，再考虑通配符匹配
// 第一个返回值是子节点
// 第二个是标记是否是路径参数
// 第三个标记命中了没有
func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	child, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.starChild, false, n.starChild != nil
	}
	return child, false, ok
}

type tree struct {
	root *node
}

type node struct {
	route string

	path string

	// 静态匹配的节点
	// 子 path 到子节点的映射
	children map[string]*node

	// 通配符匹配
	starChild *node

	// 加一个路径参数
	paramChild *node

	// 用户注册的处理逻辑
	handler HandleFunc
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}
