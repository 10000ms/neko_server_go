package core

import (
    "fmt"
    "strings"
)

type node int32

const (
    notTypeNode node = 0
    htmlNode    node = 1
    ifNode      node = 2
    forNode     node = 3
    endNode     node = 4
    elseNode    node = 5
    valueNode   node = 6
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
    childNodes  []TemplateNode
    fatherNode  *TemplateNode
    environment map[string]interface{}
    done        bool
}

func (self *TemplateNode) Init() {
    self.typeForNode()
    self.analyse()
}

// 分析确实自身节点类型
func (self *TemplateNode) typeForNode() {
    if self.htmlOnly == true {
        self.nodeType = htmlNode
    } else {
        tempContent := strings.TrimSpace(self.content)
        if tempContent[:3] == "if " {
            self.nodeType = ifNode
        } else if tempContent[:4] == "for " {
            self.nodeType = forNode
        } else if tempContent == "else" {
            self.nodeType = elseNode
        } else if tempContent == "end" {
            self.nodeType = endNode
        } else {
            self.nodeType = valueNode
        }
    }
}

// 根据自身不同的节点类型，使用不同的方法进行分析和处理

func (self *TemplateNode) analyse() {
    typeAnalyseMethods := map[node]func(){
        htmlNode:  self.analyseHtmlNode,
        ifNode:    self.analyseIfNode,
        forNode:   self.analyseForNode,
        endNode:   self.analyseEndNode,
        elseNode:  self.analyseElseNode,
        valueNode: self.analyseValueNode,
    }
    m := typeAnalyseMethods[self.nodeType]
    m()
}

func (self *TemplateNode) analyseHtmlNode() {
    self.done = true
}

func (self *TemplateNode) analyseIfNode() {
    tempContent := strings.TrimSpace(self.content)
    self.content = tempContent[3:]
}

func (self *TemplateNode) analyseForNode() {
    tempContent := strings.TrimSpace(self.content)
    tempContent = tempContent[4:]
    self.contentList = strings.Split(tempContent, " ")
}

func (self *TemplateNode) analyseEndNode() {
    self.content = ""
    self.fatherNode.done = true
    self.done = true
}

func (self *TemplateNode) analyseElseNode() {
    self.content = ""
    self.done = true
}

func (self *TemplateNode) analyseValueNode() {
    valueName := strings.TrimSpace(self.content)
    value := strings.Split(valueName, ".")
    self.contentList = value
    self.done = true
}


func (self *TemplateNode) _addContent(content string, htmlOnly bool) TemplateNode {
    if self.nodeType == ifNode || self.nodeType == forNode {
        newNode := TemplateNode{
            content,
            nil,
            htmlOnly,
            notTypeNode,
            nil,
            self,
            make(map[string]interface{}),
            false,
        }
        newNode.Init()
        self.childNodes = append(self.childNodes, newNode)
        var returnNode TemplateNode
        if newNode.done == false {
            returnNode = newNode
        } else {
            returnNode = *self
        }
        return returnNode
    } else {
        panic("当前node类型不正确")
    }
}


func (self *TemplateNode) addContent(content string, htmlOnly bool) TemplateNode {
    targetDict := map[bool]*TemplateNode{
        true:  self.fatherNode,
        false: self,
    }
    r := targetDict[self.done]._addContent(content, htmlOnly)
    return r
}

func (self *TemplateNode) UpdateEnvironment(environment map[string]interface{}) {
    for k, v := range environment {
        self.environment[k] = v
    }
}


func (self *TemplateNode) htmlFromHtmlNode() string {
    return self.content
}


func (self *TemplateNode) htmlFromIfNode() string {
    statement := self.environment[self.content].(bool)
    r := ""
    currentPart := true
    changed := false
    for _, n := range self.childNodes {
        if n.nodeType == elseNode && changed == false {
            currentPart = false
            changed = true
        } else if n.nodeType == elseNode && changed == true {
            panic("模板if语句有多个else！")
        }
        if statement == currentPart {
            r += n.htmlFromSelf(self.environment)
        }
    }
    return r
}


func (self *TemplateNode) htmlFromForNode() string {
    iterItem := self.environment[self.contentList[len(self.contentList)-1]].([]interface{})
    r := ""
    keyName := self.contentList[0]
    for _, i := range iterItem {
        newEnvironment := make(map[string]interface{})
        for k, v := range self.environment {
            newEnvironment[k] = v
        }
        newEnvironment[keyName] = i
        for _, n := range self.childNodes {
            r += n.htmlFromSelf(newEnvironment)
        }
    }
    return r
}


func (self *TemplateNode) htmlFromEndNode() string {
    return self.content
}


func (self *TemplateNode) htmlFromElseNode() string {
    return self.content
}


// 传入的environment必须为map[string]interface{}
func (self *TemplateNode) htmlFromValueNode() string {
    var t interface{}
    t = self.environment
    for _, c := range self.contentList {
        t = t.(map[string]interface{})[c]
    }
    return fmt.Sprintf("%s", t)
}


// 根据自身不同的类型，使用对应方法，输出html
func (self *TemplateNode) htmlFromSelf(environment map[string]interface{}) string {
    htmlMethods := map[node]func() string {
        htmlNode:  self.htmlFromHtmlNode,
        ifNode:    self.htmlFromIfNode,
        forNode:   self.htmlFromForNode,
        endNode:   self.htmlFromEndNode,
        elseNode:  self.htmlFromElseNode,
        valueNode: self.htmlFromValueNode,
    }
    self.environment = environment
    m := htmlMethods[self.nodeType]
    return m()
}


// 节点管理器，存放主节点
type TemplateNodeManage struct {
    nodeList    []TemplateNode
    currentNode TemplateNode
}

func (self *TemplateNodeManage) htmlFromNode(environment map[string]interface{}) string {
    var html string
    for _, i := range environment {
        html += fmt.Sprintf("%s", i)
    }
    return html
}

func (self *TemplateNodeManage) addContent(content string, htmlOnly bool) TemplateNode {
    if len(self.nodeList) > 0 && self.nodeList[len(self.nodeList)-1].done == false {
        r := self.currentNode.addContent(content, htmlOnly)
        self.currentNode = r
        return r
    } else {
        t := TemplateNode{
            content,
            nil,
            htmlOnly,
            notTypeNode,
            nil,
            nil,
            make(map[string]interface{}),
            false,
        }
        t.Init()
        self.nodeList = append(self.nodeList, t)
        self.currentNode = t
        return t
    }
}

// 渲染器
type Render struct {
    templateLoader TemplateLoader
}

// 获取搜索模板语法的标签和状态
func (self *Render) renewSearchState(newItem string) (string, bool) {
    if newItem == "{" {
        return "}", false
    } else {
        return "{", true
    }
}

// 将文本内容转化为分析好的node管理器
func (self *Render) contentToNodeManage(content string) TemplateNodeManage {
    manage := TemplateNodeManage{}
    lastIndex := 0
    index := 0
    // html用来确定当前搜索内容是否为html only
    searchItem, html := self.renewSearchState("}")
    runeContent := []rune(content)
    for index < len(runeContent) {
        if string(runeContent[index]) == searchItem && string(runeContent[index+1]) == searchItem {
            manage.addContent(string(runeContent[lastIndex:index]), html)
            searchItem, html = self.renewSearchState(searchItem)
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
func (self *Render) Render(template string, environment map[string]interface{}) string {
    content := self.templateLoader.GetSource(template)
    templateManage := self.contentToNodeManage(content)
    html := templateManage.htmlFromNode(environment)
    return html
}
