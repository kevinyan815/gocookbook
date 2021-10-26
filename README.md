# Golang开发手记


用Go语言做开发，在这个Repository里整理一些常用的案例，计划慢慢积累作为以后开发的CookBook。

仓库里所有知识点对应的代码示例都可正常运行，拿来直接应用到生产项目上也不会有问题。因为目的是积累Go语言开发的案头书，所以并不会讲源码分析之类的东西，如果想更多了解 Go 语言各种内部原理和源代码解读欢迎关注我的公众号 **「网管叨bi叨」** ，那里除了应用还会用大量的原理分析。

![#公众号：网管叨bi叨](https://cdn.learnku.com/uploads/images/202109/24/6964/ZXgD1fAlOU.png!large)

### 课程推荐

[视频课--Go语言从入门到实战](http://gk.link/a/10tim)

[实战课--Go并发编程实战20讲](http://gk.link/a/10nB6)

## 目录
- 初始化
  - [Go应用初始化工作的执行顺序](https://github.com/kevinyan815/gocookbook/issues/24)

- 项目工程
  - [依赖管理工具GOMODULE](https://mp.weixin.qq.com/s/xtvTUl2IZFQ79dSR_m-b7A)
  - [GoModules 管理私有依赖模块](https://mp.weixin.qq.com/s/8E1PwnglrS18hZsUEvE-Qw)
  - [常用编码规范](https://github.com/kevinyan815/gocookbook/issues/61)
- 字符串
  - [看透Go语言的字符串](https://github.com/kevinyan815/gocookbook/issues/40)
  - [操作中文字符串](https://github.com/kevinyan815/gocookbook/issues/11)
  - [常用字符串操作](https://yourbasic.org/golang/string-functions-reference-cheat-sheet/)
  - [高性能地拼接字符串](https://github.com/kevinyan815/gocookbook/issues/68)
- 数组
  - [数组的上限推导和越界检查](https://github.com/kevinyan815/gocookbook/issues/37)
- Slice切片
  - [声明和初始化](https://github.com/kevinyan815/gocookbook/issues/3)
  - [追加和删除元素](https://github.com/kevinyan815/gocookbook/issues/4)
  - [过滤重复元素](https://github.com/kevinyan815/gocookbook/issues/5)
  - [排序结构体切片](https://github.com/kevinyan815/gocookbook/issues/12)
  - [切片并非引用类型](https://github.com/kevinyan815/gocookbook/issues/38)
- Map
  - [(通识概念)哈希表的设计原理](https://github.com/kevinyan815/gocookbook/issues/39)
  - [声明和初始化](https://github.com/kevinyan815/gocookbook/issues/6)
  - [不要向nil map写入键值](https://github.com/kevinyan815/gocookbook/issues/7)
  - [修改map](https://github.com/kevinyan815/gocookbook/issues/8)
  - [遍历map](https://github.com/kevinyan815/gocookbook/issues/15)
  - [make 和 new](https://github.com/kevinyan815/gocookbook/issues/53)
- 读写数据
  - [编码JSON](https://github.com/kevinyan815/gocookbook/issues/2)
  - [解码JSON](https://github.com/kevinyan815/gocookbook/issues/1)
  - [逐行读取文件](https://github.com/kevinyan815/gocookbook/issues/13)
  - [Go语言IO库使用方法汇总 (进行IO操作时到底应该用哪个库)](https://github.com/kevinyan815/gocookbook/issues/62)
  - [Go语言文件操作大全](https://mp.weixin.qq.com/s/dQUEq0lJekEUH4CHEMwANw)
  - [字节序：大端序和小端序](https://mp.weixin.qq.com/s/ri2tt4nvEJub-wEsh0WPPA)
  - [用Golang读写HTTP请求(附Options设计模式实现)](https://github.com/kevinyan815/gocookbook/issues/64)
  
- 指针
  - [用法和使用限制](https://github.com/kevinyan815/gocookbook/issues/41)
  - [uintptr 和 unsafer.Pointer](https://github.com/kevinyan815/gocookbook/issues/42)
  - [扩展阅读：内存对齐](https://github.com/kevinyan815/gocookbook/issues/43)
- 接口 
  - [认识Go的接口](https://github.com/kevinyan815/gocookbook/issues/45)
  - [Go接口的类型和方法的接收者](https://github.com/kevinyan815/gocookbook/issues/46)
  - [接口的类型转换和断言](https://github.com/kevinyan815/gocookbook/issues/47)
  - [接口调用时的动态派发](https://github.com/kevinyan815/gocookbook/issues/67)
- [Range 迭代](https://github.com/kevinyan815/gocookbook/issues/15)
- 函数
  - [调用惯例和参数传递](https://github.com/kevinyan815/gocookbook/issues/44)
  - [defer的用法和行为分析](https://github.com/kevinyan815/gocookbook/issues/51)
  - [panic和recover](https://github.com/kevinyan815/gocookbook/issues/52)

- 错误处理
  - [关于Golang错误处理的一些建议](https://github.com/kevinyan815/gocookbook/issues/66)
- 包
  - [内部包](https://github.com/kevinyan815/gocookbook/issues/58)
- 标准库
  - [正则表达式](https://github.com/kevinyan815/gocookbook/issues/9)
  - [time](https://github.com/kevinyan815/gocookbook/issues/14)
- 系统编程
  - [命令行flag](https://github.com/kevinyan815/gocookbook/issues/36)
  - [监听系统信号](https://github.com/kevinyan815/gocookbook/issues/55)
- 并发编程
  - [Context上下文](https://github.com/kevinyan815/gocookbook/issues/50)
  - [Channel 基本概念和用法](https://github.com/kevinyan815/gocookbook/issues/54)
  - [用WaitGroup进行协同等待](https://github.com/kevinyan815/gocookbook/issues/34)
  - [ErrorGroup 兼顾协同等待和错误传递](https://github.com/kevinyan815/gocookbook/issues/35)
  - [Reset计时器的正确姿势](https://github.com/kevinyan815/gocookbook/issues/17)
  - [结合cancelCtx, Timer, Goroutine, Channel的一个例子](https://github.com/kevinyan815/gocookbook/issues/18)
  - [使用WaitGroup, Channel和Context打造一个并发用户标签查询器](https://github.com/kevinyan815/gocookbook/issues/21)
  - [使用sync.Cond实现一个有限容量的队列](https://github.com/kevinyan815/gocookbook/issues/22)
  - [使用信号量控制有限资源的并发访问](https://github.com/kevinyan815/gocookbook/issues/30)
  - [使用Chan扩展互斥锁的功能](https://github.com/kevinyan815/gocookbook/issues/25)
  - [用SingleFlight合并重复请求](https://github.com/kevinyan815/gocookbook/issues/31)
  - [CyclicBarrier 循环栅栏](https://github.com/kevinyan815/gocookbook/issues/32)
  - [原子操作的用法详解](https://github.com/kevinyan815/gocookbook/issues/65)
- 反射
  - [Go反射的使用教程](https://github.com/kevinyan815/gocookbook/issues/69)
- 数据库访问
  - [使用标准库 database/sql 访问数据库](https://mp.weixin.qq.com/s/bhsFCXTZ_TBP0EvyRM-bdA)
  - [使用ORM库 gorm 访问数据库](https://mp.weixin.qq.com/s/N-ZAgRrEu2FJBlApIhuVsg)
  - [GORM 指南](https://gorm.io/zh_CN/docs/index.html)
- 线上问题解决实录
  - [重定向运行时panic到日志文件](https://github.com/kevinyan815/gocookbook/issues/19)
  - [用Go的交叉编译和条件编译让自己的软件包运行在多平台上](https://github.com/kevinyan815/gocookbook/issues/20)
  - [在容器里怎么设置GOMAXPRCS](https://github.com/kevinyan815/gocookbook/issues/57)
  - [预防并发搞垮友军的几个方法](https://github.com/kevinyan815/gocookbook/issues/63)
- 编译原理
  - [Go程序的编译原理](https://github.com/kevinyan815/gocookbook/issues/56)
- 一些有意思的小程序
  - [一个简单的概率抽奖工具](https://github.com/kevinyan815/gocookbook/issues/23)
  - [限流算法之计数器](https://github.com/kevinyan815/gocookbook/issues/29)
  - [限流算法之滑动窗口](https://github.com/kevinyan815/gocookbook/issues/26)
  - [限流算法之漏桶](https://github.com/kevinyan815/gocookbook/issues/28)
  - [限流算法之令牌桶](https://github.com/kevinyan815/gocookbook/issues/27)
  - [并发趣题--H2O制造工厂](https://github.com/kevinyan815/gocookbook/issues/33)
  - [可以自解释的Token生成算法](https://github.com/kevinyan815/gocookbook/blob/master/codes/gen_token/main.go)
  - [生成分布式链路追踪traceid和spanid的算法](https://github.com/kevinyan815/gocookbook/blob/master/codes/trace_span/main.go)
  - [一个带阻塞限流器的HTTP客户端](https://github.com/kevinyan815/gocookbook/blob/master/codes/http_client_with_rate/http_rl_client.go)
- gRPC应用实践
  - [interceptor拦截器--gRPC的Middleware](https://github.com/kevinyan815/gocookbook/issues/60)
