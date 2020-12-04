package core

import (
	"neko_server_go/utils"
	"regexp"
	"strings"
)

type routeNode struct {
	part            string
	originTotalPath string
	pattern         *regexp.Regexp
	handler         *NekoHandlerFunc
	children        []*routeNode
}

type Router map[string]NekoHandlerFunc

type RouterManager struct {
	OriginRouter  *Router
	rootRouteNode *routeNode
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
*/
func (r *routeNode) CheckMatchChildren(part string, pathParams *map[string]string) (*routeNode, bool) {
	for _, c := range r.children {
		if c != nil && c.pattern.MatchString(part) {
			utils.LogDebug("Match Child: ", part, " node: ", c.originTotalPath)

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

func (r RouterManager) addRouteNode(patternList []string, handler NekoHandlerFunc, currencyNode *routeNode, originPath string) {
	if len(patternList) == 0 {
		if handler != nil {
			utils.LogDebug("path: ", originPath, ", add Node handler to ", currencyNode.part, ", ", currencyNode, ", ", handler)
			currencyNode.handler = &handler
			currencyNode.originTotalPath = originPath
		}
	} else {
		// 长度不为0代表有子节点
		c, in := currencyNode.CheckHaveChildren(patternList[0])
		if in == true {
			//utils.LogDebug("path： ", originPath, ", ", patternList, ", in Node: ", c.originTotalPath, ", currencyNode: ", currencyNode.originTotalPath)
			r.addRouteNode(patternList[1:], handler, c, originPath)
		} else {
			// 不在就创建一个子节点
			child := routeNode{
				part:            patternList[0],
				pattern:         regexp.MustCompile(pathToRegexp(patternList[0])),
			}
			//utils.LogDebug("path： ", originPath, ", ", patternList, ", create new Node: ", child.part)
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
		rootRouteNode: &root,
	}
	for pattern, handler := range *router {
		patternList := splitPath(pattern)
		utils.LogDebug("sasdasdsdad", pattern, handler)
		r.addRouteNode(patternList, &handler, &root, pattern)
	}
	return &r
}

/**
根据pattern获取对应的路由，找不到就返回defaultHandler路由
*/
func (r *RouterManager) MatchHandler(path string, defaultHandler *NekoHandlerFunc) (*NekoHandlerFunc, *map[string]string) {
	patternList := splitPath(path)
	pathParams := make(map[string]string)
	h := matchHandlerAndGetPathParams(patternList, r.rootRouteNode, &pathParams)
	if h == nil {
		// 没有找到handler
		return defaultHandler, nil
	} else {
		return h, &pathParams
	}
}

func matchHandlerAndGetPathParams(patternList []string, currencyNode *routeNode, pathParams *map[string]string) *NekoHandlerFunc {
	if len(patternList) == 0 {
		utils.LogDebug("match handler: ", currencyNode.originTotalPath, ", ", currencyNode)
		return currencyNode.handler
	} else {
		c, in := currencyNode.CheckMatchChildren(patternList[0], pathParams)
		if in == true {
			return matchHandlerAndGetPathParams(patternList[1:], c, pathParams)
		}
	}
	return nil

}

var DefaultRouter = Router{
	"/": DefaultIndexHandler,
}
