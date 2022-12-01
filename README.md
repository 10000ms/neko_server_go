# 基于Golang实现的MVC Web应用框架：neko server go

项目说明

1. 轻依赖的完整 MVC 功能的 IO 多路复用技术的 Web 应用框架
2. 功能完整封装并暴露合适的API，如浏览器请求会封装成 Request，返回的 Response会转成为对应的响应返回浏览器，路由规则会被解析生成 Router
3. 模板管理器：可以读取模板文件目录，处理对应的模板函数，渲染成正确的 HTML 文件
