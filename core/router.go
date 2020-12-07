package core

import (
	"regexp"
	"strings"
)

type routeNode struct {
	part            string
	originTotalPath string
	pattern         *regexp.Regexp
	handler         *MethodsHandler
	children        []*routeNode
}

type Router map[string]MethodsHandler

type RouterManager struct {
	OriginRouter  *Router
	rootRouteNode *routeNode
}

func CreateMethodsHandler(methods HttpMethods, handler NekoHandlerFunc) MethodsHandler {
	r := MethodsHandler{
		HttpMethod: methods,
		Handler:    &handler,
	}
	return r
}

/*
判断节点有没有这个子节点
*/
func (r *routeNode) CheckHaveChildren(part string) (*routeNode, bool) {
	for _, c := range r.children {
		if c != nil && c.part == part {
			return c, true
		}
	}
	return nil, false
}

/*
正则需要全匹配
*/
func pathToRegexp(part string) string {
	if !strings.HasPrefix(part, "^") {
		part = "^" + part
	}
	if !strings.HasSuffix(part, "$") {
		part = part + "$"
	}
	return part
}

/*
判断是否匹配上这个子节点

最后一个节点要匹配方法
*/
func (r *routeNode) CheckMatchChildren(
	methods string,
	patternList []string,
	pathParams *map[string]string,
) (*routeNode, bool) {
	part := patternList[0]
	for _, c := range r.children {
		if c != nil && c.pattern.MatchString(part) {
			if len(patternList) == 1 {
				// 最后一个节点要匹配方法
				if c.handler == nil || string(c.handler.HttpMethod) != methods {
					continue
				}
			}

			// 处理path传参的情况
			match := c.pattern.FindStringSubmatch(part)
			groupNames := c.pattern.SubexpNames()

			p := *pathParams
			for i, param := range groupNames {
				if i != 0 && param != "" { // 第一个分组为空（也就是整个匹配）
					p[param] = match[i]
				}
			}

			return c, true
		}
	}
	return nil, false
}

/*
合并router
*/
func (r Router) MergeRouter(router *Router) *Router {
	for pattern, handler := range *router {
		r[pattern] = handler
	}
	return &r
}

func (r RouterManager) addRouteNode(
	patternList []string,
	handler MethodsHandler,
	currencyNode *routeNode,
	originPath string,
) {
	if len(patternList) == 0 {
		currencyNode.handler = &handler
		currencyNode.originTotalPath = originPath
	} else {
		// 长度不为0代表有子节点
		c, in := currencyNode.CheckHaveChildren(patternList[0])
		if in == true {
			r.addRouteNode(patternList[1:], handler, c, originPath)
		} else {
			// 不在就创建一个子节点
			child := routeNode{
				part:    patternList[0],
				pattern: regexp.MustCompile(pathToRegexp(patternList[0])),
			}
			currencyNode.children = append(currencyNode.children, &child)
			r.addRouteNode(patternList[1:], handler, &child, originPath)
		}
	}
}

func splitPath(path string) []string {
	pathList := strings.Split(path, "/")
	if pathList[len(pathList)-1] == "" {
		pathList = pathList[1 : len(pathList)-1]
	} else {
		pathList = pathList[1:]
	}
	return pathList
}

func initRouterManager(router *Router) *RouterManager {
	// 先创建一个根节点
	root := routeNode{
		part: "",
	}
	r := RouterManager{
		OriginRouter: router,
		rootRouteNode: &root,
	}
	for pattern, handler := range *router {
		patternList := splitPath(pattern)
		r.addRouteNode(patternList, handler, &root, pattern)
	}
	return &r
}

/**
根据pattern获取对应的路由，找不到就返回defaultHandler路由
*/
func (r *RouterManager) MatchHandler(
	methods string,
	path string,
	defaultHandler *NekoHandlerFunc,
) (*NekoHandlerFunc, *map[string]string) {
	patternList := splitPath(path)
	pathParams := make(map[string]string)
	h := matchHandlerAndGetPathParams(methods, patternList, r.rootRouteNode, &pathParams)
	if h == nil {
		// 没有找到handler
		return defaultHandler, nil
	} else {
		return h, &pathParams
	}
}

func matchHandlerAndGetPathParams(
	methods string,
	patternList []string,
	currencyNode *routeNode,
	pathParams *map[string]string,
) *NekoHandlerFunc {
	if len(patternList) == 0 {
		return currencyNode.handler.Handler
	} else {

		c, in := currencyNode.CheckMatchChildren(methods, patternList, pathParams)
		if in == true {
			return matchHandlerAndGetPathParams(methods, patternList[1:], c, pathParams)
		}
	}
	return nil

}

var DefaultRouter = Router{
	"/": CreateMethodsHandler(GetMethodsHandler, DefaultIndexHandler),
}
