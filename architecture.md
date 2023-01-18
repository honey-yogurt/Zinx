# v0.1 
server.Start() 直接监听地址，acceptTCP，对conn只进行回显逻辑。


# v0.2
封装了conn，Conn 中嵌入了 handleAPI ， 提供了一个 func（回调逻辑） 进行绑定业务逻辑，处理客户端的请求数据。


# v0.3
将 v0.2 的 handleAPI 替换成一个固定路由，在初始化时候注册路由，conn 中增加路由，根据模板模式，依次执行路由中的方法。