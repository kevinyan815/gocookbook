> 以下内容为几年前学Go语言时记录的，发布在这里让这个项目对更多人能起到参考作用。 更多时下流行技术的应用和实战教程，可通过我公众号「网管叨bi叨」每周的推文来学习

## 方法

### 方法声明

在函数声明时，在其名字之前放上一个变量，即是一个方法。这个附加的参数会将该函数附加到这种类型上，即相当于为这种类型定义了一个独占的方法。

```go
package geometry

import "math"

type Point struct{ X, Y float64 }

// traditional function
func Distance(p, q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
    return math.Hypot(q.X-p.X, q.Y-p.Y)
}
```

上面的代码里那个附加的参数p，叫做方法的接收器(receiver)。在Go语言中，我们并不会像其它语言那样用this或者self作为接收器；我们可以任意的选择接收器的名字。建议是可以使用其类型的第一个字母，比如这里使用了Point的首字母p。

在方法调用过程中，接收器参数一般会在方法名之前出现。这和方法声明是一样的，都是接收器参数在方法名字之前。下面是例子：

```Go
p := Point{1, 2}
q := Point{4, 6}
fmt.Println(Distance(p, q)) // "5", function call
fmt.Println(p.Distance(q))  // "5", method call
```

可以看到，上面的两个函数调用都是Distance，但是却没有发生冲突。第一个Distance的调用实际上用的是包级别的函数geometry.Distance，而第二个则是使用刚刚声明的Point，调用的是Point类下声明的Point.Distance方法。这种`p.Distance`的表达式叫做选择器，因为他会选择合适的对应p这个对象的`Distance`方法来执行。

因为每种类型都有其方法的命名空间，我们在用Distance这个名字的时候，不同的Distance调用指向了不同类型里的Distance方法。

```Go
// A Path is a journey connecting the points with straight lines.
type Path []Point
// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
    sum := 0.0
    for i := range path {
        if i > 0 {
            sum += path[i-1].Distance(path[i])
        }
    }
    return sum
}
```

Path是一个命名的slice类型，而不是Point那样的struct类型，然而我们依然可以为它定义方法。两个Distance方法有不同的类型。他们两个方法之间没有任何关系，尽管Path的Distance方法会在内部调用`Point.Distance`方法来计算每个连接邻接点的线段的长度。

Go和很多其它的面向对象的语言不太一样。在Go语言里，我们可以为一些简单的数值、字符串、slice、map来定义一些附加行为很方便。方法可以被声明到任意类型，只要不是一个指针或者一个interface（接收者不能是一个指针类型，但是它可以是任何其他允许类型的指针）。

对于一个给定的类型，其内部的方法都必须有唯一的方法名，但是不同的类型却可以有同样的方法名，比如我们这里Point和Path就都有Distance这个名字的方法；所以我们没有必要非在方法名之前加类型名来消除歧义，比如PathDistance。在上面两个对Distance名字的方法的调用中，编译器会根据方法的名字以及接收器来决定具体调用的是哪一个函数。

### 指针对象的方法

当调用一个函数时，会对其每一个参数值进行拷贝，如果一个函数需要更新一个变量，或者函数的其中一个参数实在太大我们希望能够避免进行这种默认的拷贝，这种情况下我们就需要用到指针了。对应到我们这里用来更新接收器的对象的方法，当这个接受者变量本身比较大时，我们就可以用其指针而不是对象来声明方法，如下：

```go
func (p *Point) ScaleBy(factor float64) {
    p.X *= factor
    p.Y *= factor
}
```

这个方法的名字是`(*Point).ScaleBy`。这里的括号是必须的；没有括号的话这个表达式可能会被理解为`*(Point.ScaleBy)`。

- 在现实的程序里，一般会约定如果Point这个类有一个指针作为接收器的方法，那么所有Point的方法都必须有一个指针接收器，即使是那些并不需要这个指针接收器的函数。我们在这里打破了这个约定只是为了展示一下两种方法的异同而已。

- 不管你的method的receiver是指针类型还是非指针类型，都是可以通过指针/非指针类型进行调用的，编译器会帮你做类型转换。

```go
  p := Point{1, 2}
  pptr := &p
  p.ScaleBy(2) // implicit (&p)
  pptr.Distance(q) // implicit (*pptr)
```

- 在声明一个method的receiver是指针还是非指针类型时，你需要考虑两方面的内部，第一方面是这个对象本身是不是特别大，如果声明为非指针变量时，调用会产生一次拷贝；第二方面是如果你用指针类型作为receiver，那么你一定要注意，这种指针类型指向的始终是一块内存地址，就算你对其进行了拷贝（指针调用时也是值拷贝，只不过指针的值是一个内存地址，所以在函数里的指针与调用方的指针变量是两个不同的指针但是指向了相同的内存地址）。

### Nil也是一个合法的接收器类型

- 就像一些函数允许nil指针作为参数一样，方法理论上也可以用nil指针作为其接收器，尤其当nil对于对象来说是合法的零值时，比如map或者slice。在下面的简单int链表的例子里，nil代表的是空链表：

```go
  // An IntList is a linked list of integers.
  // A nil *IntList represents the empty list.
  type IntList struct {
      Value int
      Tail  *IntList
  }
  // Sum returns the sum of the list elements.
  func (list *IntList) Sum() int {
      if list == nil {
          return 0
      }
      return list.Value + list.Tail.Sum()
  }
```

  当你定义一个允许nil作为接收器的方法的类型时，在类型前面的注释中指出nil变量代表的意义是很有必要的，就像我们上面例子里做的这样。

### 通过嵌入结构体来扩展类型

- 下面的ColoredPoint类型

```go
  import "image/color"
  
  type Point struct{ X, Y float64 }
  
  type ColoredPoint struct {
      Point
      Color color.RGBA
  }
```

  内嵌可以使我们在定义ColoredPoint时得到一种句法上的简写形式，并使其包含Point类型所具有的一切字段和方法。

```go
  var cp ColoredPoint
  cp.X = 1
  fmt.Println(cp.Point.X) // "1"
  cp.Point.Y = 2
  fmt.Println(cp.Y) // "2"
  
  red := color.RGBA{255, 0, 0, 255}
  blue := color.RGBA{0, 0, 255, 255}
  var p = ColoredPoint{Point{1, 1}, red}
  var q = ColoredPoint{Point{5, 4}, blue}
  fmt.Println(p.Distance(q.Point)) // "5"
  p.ScaleBy(2)
  q.ScaleBy(2)
  fmt.Println(p.Distance(q.Point)) // "10"
```

  通过内嵌结构体可以使我们定义字段特别多的复杂类型，我们可以将字段先按小类型分组，然后定义小类型的方法，之后再把它们组合起来。

- 内嵌字段会指导编译器去生成额外的包装方法来委托已经声明好的方法，和下面的形式是等价的：

```go
  func (p ColoredPoint) Distance(q Point) float64 {
      return p.Point.Distance(q)
  }
  
  func (p *ColoredPoint) ScaleBy(factor float64) {
      p.Point.ScaleBy(factor)
  }
```

  当Point.Distance被第一个包装方法调用时，它的接收器值是p.Point，而不是p，当然了，在Point类的方法里，你是访问不到ColoredPoint的任何字段的。

- 方法只能在命名类型(像Point)或者指向类型的指针上定义，但是多亏了内嵌，我们给匿名struct类型来定义方法也有了手段。这个例子中我们为变量起了一个更具表达性的名字：cache。因为sync.Mutex类型被嵌入到了这个struct里，其Lock和Unlock方法也就都被引入到了这个匿名结构中了，这让我们能够以一个简单明了的语法来对其进行加锁解锁操作。

```go
  var cache = struct {
      sync.Mutex
      mapping map[string]string
  }{
      mapping: make(map[string]string),
  }
  
  
  func Lookup(key string) string {
      cache.Lock()
      v := cache.mapping[key]
      cache.Unlock()
      return v
  }
```

### 方法值和方法表达式

- 我们经常选择一个方法，并且在同一个表达式里执行，比如常见的p.Distance()形式，实际上将其分成两步来执行也是可能的。p.Distance叫作“选择器”，选择器会返回一个方法"值"->一个将方法(Point.Distance)绑定到特定接收器变量的函数。因为已经在前文中指定过了，这个函数可以不通过指定其接收器即可被调用，只要传入函数的参数即可：

```go
  p := Point{1, 2}
  q := Point{4, 6}
  
  distanceFromP := p.Distance        // method value
  fmt.Println(distanceFromP(q))      // "5"
  var origin Point                   // {0, 0}
  fmt.Println(distanceFromP(origin)) // "2.23606797749979", sqrt(5)
  
  scaleP := p.ScaleBy // method value
  scaleP(2)           // p becomes (2, 4)
  scaleP(3)           //      then (6, 12)
  scaleP(10)          //      then (60, 120)
```

- 当T是一个类型时，方法表达式可能会写作T.f或者(*T).f，会返回一个函数"值"，这种函数会将其第一个参数用作接收器，所以可以用通常(译注：不写选择器)的方式来对其进行调用：

```go
  p := Point{1, 2}
  q := Point{4, 6}
  
  distance := Point.Distance   // method expression
  fmt.Println(distance(p, q))  // "5"
  fmt.Printf("%T\n", distance) // "func(Point, Point) float64"
  
  scale := (*Point).ScaleBy
  scale(&p, 2)
  fmt.Println(p)            // "{2 4}"
  fmt.Printf("%T\n", scale) // "func(*Point, float64)"
  // 译注：这个Distance实际上是指定了Point对象为接收器的一个方法func (p Point) Distance()，
  // 但通过Point.Distance得到的函数需要比实际的Distance方法多一个参数，
  // 即其需要用第一个额外参数指定接收器，后面排列Distance方法的参数。
```

- 当你根据一个变量来决定调用同一个类型的哪个函数时，方法表达式就显得很有用了。你可以根据选择来调用接收器各不相同的方法。下面的例子，变量op代表Point类型的addition或者subtraction方法，Path.TranslateBy方法会为其Path数组中的每一个Point来调用对应的方法：

```go
  type Point struct{ X, Y float64 }
  
  func (p Point) Add(q Point) Point { return Point{p.X + q.X, p.Y + q.Y} }
  func (p Point) Sub(q Point) Point { return Point{p.X - q.X, p.Y - q.Y} }
  
  type Path []Point
  
  func (path Path) TranslateBy(offset Point, add bool) {
      var op func(p, q Point) Point
      if add {
          op = Point.Add
      } else {
          op = Point.Sub
      }
      for i := range path {
          // Call either path[i].Add(offset) or path[i].Sub(offset).
          path[i] = op(path[i], offset)
      }
  }
```

### 封装

- 一个对象的变量或者方法如果对调用方是不可见的话，一般就被定义为“封装”。封装有时候也被叫做信息隐藏，同时也是面向对象编程最关键的一个方面。
- Go语言只有一种控制可见性的手段：大写首字母的标识符会从定义它们的包中被导出，小写字母的则不会。
- 这种基于名字的手段使得在语言中最小的封装单元是package，而不是像其它语言一样的`Class`。一个struct类型的字段对同一个包的所有代码都有可见性，无论你的代码是写在一个函数还是一个方法里。

- 封装提供了三方面的优点。首先，因为调用方不能直接修改对象的变量值，其只需要关注少量的语句并且只要弄懂少量变量的可能的值即可。

  第二，隐藏实现的细节，可以防止调用方依赖那些可能变化的具体实现，这样使设计包的程序员在不破坏对外的api情况下能得到更大的自由。

  封装的第三个优点也是最重要的优点，是阻止了外部调用方对对象内部的值任意地进行修改。因为对象内部变量只可以被同一个包内的函数修改，所以包的作者可以让这些函数确保对象内部的一些值的不变性。比如下面的Counter类型允许调用方来增加counter变量的值，并且允许将这个值reset为0，但是不允许随便设置这个值(译注：因为压根就访问不到)：

```go
  type Counter struct { n int }
  func (c *Counter) N() int     { return c.n }
  func (c *Counter) Increment() { c.n++ }
  func (c *Counter) Reset()     { c.n = 0 }
```

- 只用来访问或修改内部变量的函数被称为setter或者getter，例子如下，比如log包里的Logger类型对应的一些函数。在命名一个getter方法时，我们通常会省略掉前面的Get前缀。这种简洁上的偏好也可以推广到各种类型的前缀比如Fetch，Find或者Lookup。

```go
  package log
  type Logger struct {
      flags  int
      prefix string
      // ...
  }
  func (l *Logger) Flags() int
  func (l *Logger) SetFlags(flag int)
  func (l *Logger) Prefix() string
  func (l *Logger) SetPrefix(prefix string)
```

- Go的编码风格不禁止直接导出字段。当然，一旦进行了导出，就没有办法在保证API兼容的情况下去除对其的导出，所以在一开始的选择一定要经过深思熟虑并且要考虑到包内部的一些不变量的保证，还有未来可能的变化。
