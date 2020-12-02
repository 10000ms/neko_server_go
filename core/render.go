package core

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type node int32

const (
	notTypeNode node = iota
	htmlNode
	ifNode
	forNode
	endNode
	elseNode
	valueNode
)

/*
模版的节点结构体

目前支持的节点：
1. 纯html节点
2. if 判断开始节点： 语法为 {{ if xxx }}，if和xxx中间只可以有一个空格，其中xxx一定要为environment中存在的key
3. else 判断else节点：语法为 {{ else }}
4. for 循环开始节点：语法为 {{ for yyy in xxx }}，for和xxx中间只可以有一个空格, 其中xxx一定要为environment中存在的key
5. end 结束判断和循环的节点：语法为 {{ end }}
6. value 变量取值节点：语法为 {{ xxx }} 或 {{ xxx.yyy }}等，可以从字典或者是对象中取值，
  xxx一定要为environment中存在的key或者是循环开始节点中增加定义的yyy

处理流程：
1. 先确定节点类型
2. 根据节点类型进行预处理
3. 传入数据的时候根据数据渲染不同的节点
4. 返回html
*/
type TemplateNode struct {
	content     string
	contentList []string
	htmlOnly    bool
	nodeType    node
	childNodes  []*TemplateNode
	fatherNode  *TemplateNode
	environment map[string]interface{}
	done        bool
}

func (t *TemplateNode) Init() {
	t.typeForNode()
	t.analyse()
}

// 分析确实自身节点类型
func (t *TemplateNode) typeForNode() {
	if t.htmlOnly == true {
		t.nodeType = htmlNode
	} else {
		tempContent := strings.TrimSpace(t.content)
		if len(tempContent) > 3 && tempContent[:3] == "if " {
			t.nodeType = ifNode
		} else if len(tempContent) > 4 && tempContent[:4] == "for " {
			t.nodeType = forNode
		} else if tempContent == "else" {
			t.nodeType = elseNode
		} else if tempContent == "end" {
			t.nodeType = endNode
		} else {
			t.nodeType = valueNode
		}
	}
}

// 根据自身不同的节点类型，使用不同的方法进行分析和处理

func (t *TemplateNode) analyse() {
	typeAnalyseMethods := map[node]func(){
		htmlNode:  t.analyseHtmlNode,
		ifNode:    t.analyseIfNode,
		forNode:   t.analyseForNode,
		endNode:   t.analyseEndNode,
		elseNode:  t.analyseElseNode,
		valueNode: t.analyseValueNode,
	}
	m := typeAnalyseMethods[t.nodeType]
	m()
}

func (t *TemplateNode) analyseHtmlNode() {
	t.done = true
}

func (t *TemplateNode) analyseIfNode() {
	tempContent := strings.TrimSpace(t.content)
	t.content = tempContent[3:]
}

func (t *TemplateNode) analyseForNode() {
	tempContent := strings.TrimSpace(t.content)
	tempContent = tempContent[4:]
	t.contentList = strings.Split(tempContent, " ")
}

func (t *TemplateNode) analyseEndNode() {
	t.content = ""
	t.fatherNode.done = true
	t.done = true
}

func (t *TemplateNode) analyseElseNode() {
	t.content = ""
	t.done = true
}

func (t *TemplateNode) analyseValueNode() {
	valueName := strings.TrimSpace(t.content)
	value := strings.Split(valueName, ".")
	t.contentList = value
	t.done = true
}

func (t *TemplateNode) _addContent(content string, htmlOnly bool) *TemplateNode {
	if t.nodeType == ifNode || t.nodeType == forNode {
		newNode := TemplateNode{
			content,
			nil,
			htmlOnly,
			notTypeNode,
			nil,
			t,
			make(map[string]interface{}),
			false,
		}
		newNode.Init()
		t.childNodes = append(t.childNodes, &newNode)
		var returnNode *TemplateNode
		if newNode.done == false {
			returnNode = &newNode
		} else {
			returnNode = t
		}
		return returnNode
	} else {
		panic("当前node类型不正确")
	}
}

func (t *TemplateNode) addContent(content string, htmlOnly bool) *TemplateNode {
	// 根节点自己没办法导到父节点
	var trueNode *TemplateNode
	if t.fatherNode != nil {
		trueNode = t.fatherNode
	} else {
		trueNode = t
	}
	targetDict := map[bool]*TemplateNode{
		true:  trueNode,
		false: t,
	}
	r := targetDict[t.done]._addContent(content, htmlOnly)
	return r
}

func (t *TemplateNode) UpdateEnvironment(environment map[string]interface{}) {
	for k, v := range environment {
		t.environment[k] = v
	}
}

func (t *TemplateNode) htmlFromHtmlNode() string {
	return t.content
}

func (t *TemplateNode) htmlFromIfNode() string {
	statement := t.environment[t.content].(bool)
	r := ""
	currentPart := true
	changed := false
	for _, n := range t.childNodes {
		if n.nodeType == elseNode && changed == false {
			currentPart = false
			changed = true
		} else if n.nodeType == elseNode && changed == true {
			panic("模板if语句有多个else！")
		}
		if statement == currentPart {
			r += n.htmlFromSelf(t.environment)
		}
	}
	return r
}

func (t *TemplateNode) htmlFromForNode() string {
	iterItemKey := t.contentList[len(t.contentList)-1]
	if _, ok := t.environment[iterItemKey]; !ok {
		panic(errors.New("don't have for item: " + iterItemKey))
	}
	iterItem := t.environment[iterItemKey]
	r := ""
	keyName := t.contentList[0]

	// 复制一份老的环境参数
	newEnvironment := make(map[string]interface{})
	for k, v := range t.environment {
		newEnvironment[k] = v
	}

	iterItemReflect := reflect.ValueOf(iterItem)
	iterItemKind := iterItemReflect.Kind()
	if iterItemKind == reflect.Array || iterItemKind == reflect.Slice {
		for i := 0; i < iterItemReflect.Len(); i++ {
			for _, n := range t.childNodes {
				newEnvironment[keyName] = iterItemReflect.Index(i).Interface()
				r += n.htmlFromSelf(newEnvironment)
			}
		}
	} else {
		panic(errors.New("unexpect type of for"))
	}

	return r
}

func (t *TemplateNode) htmlFromEndNode() string {
	t.fatherNode.done = true
	return t.content
}

func (t *TemplateNode) htmlFromElseNode() string {
	return t.content
}

// 传入的environment必须为map[string]interface{}
func (t *TemplateNode) htmlFromValueNode() string {
	var temp interface{}
	temp = t.environment
	for _, c := range t.contentList {
		tReflect := reflect.ValueOf(temp)
		tKind := tReflect.Kind()
		if tKind == reflect.Map {
			temp = tReflect.MapIndex(reflect.ValueOf(c)).Interface()
		} else if tKind == reflect.Struct {
			temp = tReflect.FieldByName(c).Interface()
		} else if tKind == reflect.Ptr {
			temp = tReflect.Elem().FieldByName(c).Interface()
		} else {
			panic(errors.New("unexpect type"))
		}
	}
	// 处理一下数字
	vReflect := reflect.ValueOf(temp)
	vKind := vReflect.Kind()
	switch vKind {
	case reflect.Int:
		temp = strconv.Itoa(temp.(int))
	case reflect.Int64:
		temp = strconv.FormatInt(temp.(int64), 10)
	case reflect.Float64:
		temp = strconv.FormatFloat(temp.(float64), 'E', -1, 64)
	case reflect.Float32:
		temp = strconv.FormatFloat(float64(temp.(float32)), 'E', -1, 64)
	}
	return fmt.Sprintf("%s", temp)
}

// 根据自身不同的类型，使用对应方法，输出html
func (t *TemplateNode) htmlFromSelf(environment map[string]interface{}) string {
	htmlMethods := map[node]func() string{
		htmlNode:  t.htmlFromHtmlNode,
		ifNode:    t.htmlFromIfNode,
		forNode:   t.htmlFromForNode,
		endNode:   t.htmlFromEndNode,
		elseNode:  t.htmlFromElseNode,
		valueNode: t.htmlFromValueNode,
	}
	t.environment = environment
	m := htmlMethods[t.nodeType]
	return m()
}

// 节点管理器，存放主节点
type TemplateNodeManage struct {
	nodeList    []*TemplateNode
	currentNode *TemplateNode
}

func (t *TemplateNodeManage) htmlFromNode(environment map[string]interface{}) string {
	var html string
	for _, i := range t.nodeList {
		html += i.htmlFromSelf(environment)
	}
	return html
}

func (t *TemplateNodeManage) addContent(content string, htmlOnly bool) *TemplateNode {
	if len(t.nodeList) > 0 && t.nodeList[len(t.nodeList)-1].done == false {
		r := t.currentNode.addContent(content, htmlOnly)
		t.currentNode = r
		return r
	} else {
		templateNode := TemplateNode{
			content,
			nil,
			htmlOnly,
			notTypeNode,
			nil,
			nil,
			make(map[string]interface{}),
			false,
		}
		//templateNode.fatherNode = &templateNode
		templateNode.Init()
		t.nodeList = append(t.nodeList, &templateNode)
		t.currentNode = &templateNode
		return &templateNode
	}
}

// 渲染器
type Render struct {
	TemplateLoader TemplateLoader
}

// 获取搜索模板语法的标签和状态
func (r *Render) renewSearchState(newItem string) (string, bool) {
	if newItem == "{" {
		return "}", false
	} else {
		return "{", true
	}
}

// 将文本内容转化为分析好的node管理器
func (r *Render) contentToNodeManage(content string) TemplateNodeManage {
	manage := TemplateNodeManage{}
	lastIndex := 0
	index := 0
	// html用来确定当前搜索内容是否为html only
	searchItem, html := r.renewSearchState("}")
	runeContent := []rune(content)
	for index < len(runeContent) {
		if string(runeContent[index]) == searchItem && string(runeContent[index+1]) == searchItem {
			manage.addContent(string(runeContent[lastIndex:index]), html)
			searchItem, html = r.renewSearchState(searchItem)
			// 跳过搜索的标识符号
			index += 1
			lastIndex = index + 1
		}
		index += 1
	}
	// 把剩下部分也加入到数列当中
	manage.addContent(string(runeContent[lastIndex:index+1]), html)
	return manage
}

// 渲染出 html
func (r *Render) Render(template string, environment map[string]interface{}) string {
	content := r.TemplateLoader.GetSource(template)
	templateManage := r.contentToNodeManage(content)
	html := templateManage.htmlFromNode(environment)
	return html
}
