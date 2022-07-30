### 前言

以下内容为几年前学Go语言时记录的，发布在这里让这个项目对更多人能起到参考作用。
更多时下流行技术的应用和实战教程，可通过我公众号「网管叨bi叨」每周的推文来学习

--------
最近在读《Go 语言程序设计》这本书想通过看书巩固一下自己的基础知识，把已经积累的点通过看书学习再编织成一个网，这样看别人写的优秀代码时才能更好理解。当初工作中需要使用 Go开发项目时看了网上不少教程，比如 uknown 翻译的《the way to go》看完基本上每次使用遇到不会的时候还会再去翻阅，这次把书中的重点还有一些平时容易忽视的Go语言中各种内部结构（类型、函数、方法）的一些行为整理成读书笔记。

因为《Go 语言程序设计》不是针对初学者的，所以我只摘选最重要的部分并适当补充和调换描述顺序力求用最少的篇幅描述清楚每个知识点。

《Go 语言程序设计》在线阅读地址：https://yar999.gitbooks.io/gopl-zh/content/

如果刚接触 Go建议先去读 《the-way-to-go》在线阅读地址：https://github.com/unknwon/the-way-to-go_ZH_CN/blob/master/eBook/directory.md

### 命名：

- 函数名、变量名、常量名、类型名、包名等所有的命名，都遵循一个简单的命名规则：一个名字必须以一个字母（Unicode字母）或下划线开头，后面可以跟任意数量的字母、数字或下划线。

- 大写字母和小写字母是不同的：heapSort和Heapsort是两个不同的名字。

- 关键字不可用于命名

```
  break      default       func     interface   select
  case       defer         go       map         struct
  chan       else          goto     package     switch
  const      fallthrough   if       range       type
  continue   for           import   return      var
```

- 推荐驼峰式命名

- 名字的开头字母的大小写决定了名字在包外的可见性。如果一个名字是大写字母开头的，那么它可以被外部的包访问，包本身的名字一般总是用小写字母。

### 声明：

- Go语言主要有四种类型的声明语句：var、const、type和func，分别对应变量、常量、类型和函数。

### 变量：

- var声明语句可以创建一个特定类型的变量，然后给变量附加一个名字，并且设置变量的初始值。变量声明的一般语法如下：

```Go
  var 变量名字 类型 = 表达式
```

  其中“*类型*”或“*= 表达式*”两个部分可以省略其中的一个。如果省略的是类型信息，那么将  根据初始化表达式来推导变量的类型信息。如果初始化表达式被省略，那么将用零值初始化该变量。 数值类型变量对应的零值是0，布尔类型变量对应的零值是false，字符串类型对应的零值是空字符串，接口或引用类型（包括slice、map、chan和函数）变量对应的零值是nil。数组或结构体等聚合类型对应的零值是每个元素或字段都是对应该类型的零值。

  零值初始化机制可以确保每个声明的变量总是有一个良好定义的值，因此在Go语言中不存在未初始化的变量。这个特性可以简化很多代码，而且可以在没有增加额外工作的前提下确保边界条件下的合理行为。例如：

```Go
  var s string
  fmt.Println(s) // ""
```



### 字符串：

- 文本字符串通常被解释为采用UTF8编码的Unicode码点（rune）序列。

- 内置的len函数可以返回一个字符串中的字节数目（不是rune字符数目），索引操作s[i]返回第i个字节的字节值，i必须满足0 ≤ i< len(s)条件约束。

- 字符串的值是不可变的：一个字符串包含的字节序列永远不会被改变，当然我们也可以给一个字符串变量分配一个新字符串值。可以像下面这样将一个字符串追加到另一个字符串：

```Go
  s := "left foot"
  t := s
  s += ", right foot"
```

  这并不会导致原始的字符串值被改变，但是变量s将因为+=语句持有一个新的字符串值，但是t依然是包含原先的字符串值。

  因为字符串是不可修改的，因此尝试修改字符串内部数据的操作也是被禁止的：

```Go
  s[0] = 'L' // compile error: cannot assign to s[0]
```

- 每一个UTF8字符解码，不管是显式地调用utf8.DecodeRuneInString解码或是在range循环中隐式地解码，如果遇到一个错误的UTF8编码输入，将生成一个特别的Unicode字符'\uFFFD'，在印刷中这个符号通常是一个黑色六角或钻石形状，里面包含一个白色的问号"�"。当程序遇到这样的一个字符，通常是一个危险信号，说明输入并不是一个完美没有错误的UTF8字符串。

- 字符串的各种转换：

  string接受到[]rune的类型转换，可以将一个UTF8编码的字符串解码为Unicode字符序列：

```Go
  // "program" in Japanese katakana
  s := "プログラム"
  fmt.Printf("% x\n", s) // "e3 83 97 e3 83 ad e3 82 b0 e3 83 a9 e3 83 a0"
  r := []rune(s)
  fmt.Printf("%x\n", r)  // "[30d7 30ed 30b0 30e9 30e0]"
```

  （在第一个Printf中的`% x`参数用于在每个十六进制数字前插入一个空格。）

  如果是将一个[]rune类型的Unicode字符slice或数组转为string，则对它们进行UTF8编码：

```Go
  fmt.Println(string(r)) // "プログラム"
```

  将一个整数转型为字符串意思是生成以只包含对应Unicode码点字符的UTF8字符串：

```Go
  fmt.Println(string(65))     // "A", not "65"
  fmt.Println(string(0x4eac)) // "京"
```

  如果对应码点的字符是无效的，则用'\uFFFD'无效字符作为替换：

```Go
  fmt.Println(string(1234567)) // "�"
```

### 复合数据类型：

- 基本数据类型，它们可以用于构建程序中数据结构，是Go语言的世界的原子。以不同的方式组合基本类型可以构造出复合数据类型。我们主要讨论四种类型——数组、slice、map和结构体，数组和结构体都是有固定内存大小的数据结构。相比之下，slice和map则是动态的数据结构，它们将根据需要动态增长。

### 数组：

- 数组的长度是数组类型的一个组成部分，因此[3]int和[4]int是两种不同的数组类型。数组的长度必须是常量表达式，因为数组的长度需要在编译阶段确定。

```Go
  q := [3]int{1, 2, 3}
  q = [4]int{1, 2, 3, 4} // compile error: cannot assign [4]int to [3]int
```

### Slice:

- 长度对应slice中元素的数目；长度不能超过容量，容量一般是从slice的开始位置到底层数据的结尾位置。内置的len和cap函数分别返回slice的长度和容量。

- x[m:n]切片操作对于字符串则生成一个新字符串，如果x是[]byte的话则生成一个新的[]byte。

- slice并不是一个纯粹的引用类型，它实际上是一个类似下面结构体的聚合类型：

```Go
  type IntSlice struct {
      ptr      *int
      len, cap int
  }
```

### Map:

- 在Go语言中，一个map就是一个哈希表的引用，map类型可以写为map[K]V，其中K和V分别对应key和value。map中所有的key都有相同的类型，所有的value也有着相同的类型，但是key和value之间可以是不同的数据类型。

- map中的元素并不是一个变量，因此我们不能对map的元素进行取址操作：

```Go
  _ = &ages["bob"] // compile error: cannot take address of map element
```

  禁止对map元素取址的原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而可能导致之前的地址无效。

- map上的大部分操作，包括查找、删除、len和range循环都可以安全工作在nil值的map上，它们的行为和一个空的map类似。但是向一个nil值的map存入元素将导致一个panic异常：

```Go
  ages["carol"] = 21 // panic: assignment to entry in nil map
```

  在向map存数据前必须先创建map。

- 和slice一样，map之间也不能进行相等比较；唯一的例外是和nil进行比较。要判断两个map是否包含相同的key和value，我们必须通过一个循环实现。

### 结构体：

- 下面两个语句声明了一个叫Employee的命名的结构体类型，并且声明了一个Employee类型的变量dilbert：

```Go
  type Employee struct {
      ID        int
      Name      string
      Address   string
      DoB       time.Time
      Position  string
      Salary    int
      ManagerID int
  }
  
  var dilbert Employee
```

  dilbert结构体变量的成员可以通过点操作符访问，比如dilbert.Name和dilbert.DoB。因为dilbert是一个变量，它所有的成员也同样是变量，我们可以直接对每个成员赋值：

```Go
  dilbert.Salary -= 5000 // demoted, for writing too few lines of code
```

  或者是对成员取地址，然后通过指针访问：

```Go
  position := &dilbert.Position
  *position = "Senior " + *position // promoted, for outsourcing to Elbonia
```

- 如果结构体成员名字是以大写字母开头的，那么该成员就是导出的；这是Go语言导出规则决定的。一个结构体可能同时包含导出和未导出的成员。未导出的成员只能在包内部访问，在外部包不可访问。

- 结构体类型的零值中每个成员其类型的是零值。通常会将零值作为最合理的默认值。例如，对于bytes.Buffer类型，结构体初始值就是一个随时可用的空缓存，还有sync.Mutex的零值也是有效的未锁定状态。有时候这种零值可用的特性是自然获得的，但是也有些类型需要一些额外的工作。

- 因为结构体通常通过指针处理，可以用下面的写法来创建并初始化一个结构体变量，并返回结构体的地址：

```Go
  pp := &Point{1, 2}
```

- Go语言有一个特性让我们只声明一个成员对应的数据类型而不指名成员的名字；这类成员就叫匿名成员。匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。下面的代码中，Circle和Wheel各自都有一个匿名成员。我们可以说Point类型被嵌入到了Circle结构体，同时Circle类型被嵌入到了Wheel结构体。

```Go
  type Circle struct {
      Point
      Radius int
  }
  
  type Wheel struct {
      Circle
      Spokes int
  }
```

  得益于匿名嵌入的特性，我们可以直接访问叶子属性而不需要给出完整的路径：

```Go
  var w Wheel
  w.X = 8            // equivalent to w.Circle.Point.X = 8
  w.Y = 8            // equivalent to w.Circle.Point.Y = 8
  w.Radius = 5       // equivalent to w.Circle.Radius = 5
  w.Spokes = 20
```

- 外层的结构体不仅仅是获得了匿名成员类型的所有成员，而且也获得了该类型导出的全部的方法。这个机制可以用于将一个有简单行为的对象组合成有复杂行为的对象。
