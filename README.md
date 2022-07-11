# gee
web framework in go
使用Go语言写一个迷你Web框架
gee框架的设计以及API均参考了gin。
1. 使用New()创建 gee 的实例
2. 使用 GET()方法添加路由
3. 使用Run()启动Web服务
4. 使用 Trie 树实现动态路由(dynamic route)解析。
5. 支持两种模式:name和*filepath
