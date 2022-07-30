> 以下内容为几年前学Go语言时记录的，发布在这里让这个项目对更多人能起到参考作用。 更多时下流行技术的应用和实战教程，可通过我公众号「网管叨bi叨」每周的推文来学习

### 接口概述
- 一个具体的类型可以准确的描述它所代表的值并且展示出对类型本身的一些操作方式就像数字类型的算术操作，切片类型的索引、附加和取范围操作。总的来说，当你拿到一个具体的类型时你就知道它的本身是什么和你可以用它来做什么。

- 在Go语言中还存在着另外一种类型：接口类型。接口类型是一种抽象的类型。它不会暴露出它所代表的对象的内部结构和这个对象支持的基础操作的集合；它只会展示出自己的方法。也就是说当你有看到一个接口类型的值时，你不知道它是什么，唯一知道的就是可以通过它的方法来做什么。

- fmt.Printf它会把结果写到标准输出和fmt.Sprintf它会把结果以字符串的形式返回，实际上，这两个函数都使用了另一个函数fmt.Fprintf来进行封装。fmt.Fprintf这个函数对它的计算结果会被怎么使用是完全不知道的。

```go
package fmt

func Fprintf(w io.Writer, format string, args ...interface{}) (int, error)
func Printf(format string, args ...interface{}) (int, error) {
    return Fprintf(os.Stdout, format, args...)
}
func Sprintf(format string, args ...interface{}) string {
    var buf bytes.Buffer
    Fprintf(&buf, format, args...)
    return buf.String()
}
```

​	 Fprintf函数中的第一个参数也不是一个文件类型。它是io.Writer类型这是一个接口类型定义如下：

```go
package io

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

io.Writer类型定义了函数Fprintf和这个函数调用者之间的约定，只要是实现了`io.Writer`接口的类型都可以作为 Fprintf 函数的第一个参数。

- 一个类型可以自由的使用另一个满足相同接口的类型来进行替换被称作可替换性(LSP里氏替换)。这是一个面向对象的特征。

### 接口定义
- io.Writer类型是用的最广泛的接口之一，因为它提供了所有的类型写入bytes的抽象，包括文件类型，内存缓冲区，网络链接，HTTP客户端，压缩工具，哈希等等。io包中定义了很多其它有用的接口类型。Reader可以代表任意可以读取bytes的类型，Closer可以是任意可以关闭的值，例如一个文件或是网络链接。

```go
package io
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}
```

-  可以通过组合已有接口类型来定义新的接口类型，比如 io 包中的

```go
  type ReadWriter interface {
      Reader
      Writer
  }
  type ReadWriteCloser interface {
      Reader
      Writer
      Closer
  }
```

上面用到的语法和结构内嵌相似，我们可以用这种方式命名另一个接口，而不用声明它所有的方法。这种方式称为接口内嵌，我们可以像下面这样，不使用内嵌来声明`io.ReadWriter`接口。

```go
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Write(p []byte) (n int, err error)
}
```

或者甚至使用种混合的风格：

```go
type ReadWriter interface {
    Read(p []byte) (n int, err error)
    Writer
}
```

这三种方式定义的`io.ReadWriter`是完全一样的。

### 接口实现

- 一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。例如，*os.File类型实现了io.Reader，Writer，Closer，和ReadWriter接口。*bytes.Buffer实现了Reader，Writer，和ReadWriter这些接口，但是它没有实现Closer接口因为它不具有Close方法。Go的程序员经常会简要的把一个具体的类型描述成一个特定的接口类型。举个例子，*bytes.Buffer是io.Writer；*os.Files是io.ReadWriter。
- 接口实现的规则非常简单：表达一个类型属于某个接口只要这个类型实现这个接口。

```go
var w io.Writer
w = os.Stdout           // OK: *os.File has Write method
w = new(bytes.Buffer)   // OK: *bytes.Buffer has Write method
w = time.Second         // compile error: time.Duration lacks Write method

var rwc io.ReadWriteCloser
rwc = os.Stdout         // OK: *os.File has Read, Write, Close methods
rwc = new(bytes.Buffer) // compile error: *bytes.Buffer lacks Close method
```

- 这个规则甚至适用于等式右边本身也是一个接口类型

```go
w = rwc                 // OK: io.ReadWriteCloser has Write method
rwc = w                 // compile error: io.Writer lacks Close method
```

- 因为ReadWriter和ReadWriteCloser包含Writer的方法，所以任何实现了ReadWriter和ReadWriteCloser的类型必定也实现了Writer接口
- 对于一些命名的具体类型T；它一些方法的接收者是类型T本身然而另一些则是一个`*T`的指针。在`T`类型的变量上调用一个`*T`的方法是合法的，编译器隐式的获取了它的地址。但这仅仅是一个语法糖：T类型的值不拥有所有*T指针的方法。
- interface{}类型，它没有任何方法，但实际上interface{}被称为空接口类型是不可或缺的。因为空接口类型对实现它的类型没有要求，所以所有类型都实现了interface{}，我们可以将任意一个值赋给空接口类型。

```go
var any interface{}
any = true
any = 12.34
any = "hello"
any = map[string]int{"one": 1}
any = new(bytes.Buffer)
```

### 接口值

- 接口值由两个部分组成，一个具体的类型和那个类型的值。它们被称为接口的动态类型和动态值。

- 像Go语言这种静态类型的语言，类型是编译期的概念；因此一个类型不是一个值，提供每个类型信息的值被称为类型描述符。

- 在Go语言中，变量总是被一个定义明确的值初始化，一个接口的零值就是它的类型和值的部分都是nil。

  ![img](https://cdn.learnku.com/uploads/images/201912/28/6964/TQyayLLPRk.png!large)

- 在你非常确定接口值的动态类型是可比较类型时（比如基本类型）才可以使用`==`和`!=`对两个接口值进行比较。如果两个接口值的动态类型相同，但是这个动态类型是不可比较的（比如切片），将它们进行比较就会失败并且panic:

```go
var x interface{} = []int{1, 2, 3}
fmt.Println(x == x) // panic: comparing uncomparable type []int
```

- 下面4个语句中，变量w得到了3个不同的值。（开始和最后的值是相同的）

```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = nil
```

第一个语句定义了变量w:

```go
var w io.Writer
```

在Go语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值就是它的类型和值的部分都是nil，如图 7.1。

一个接口值基于它的动态类型被描述为空或非空，所以这是一个空的接口值。你可以通过使用w==nil或者w!=nil来判读接口值是否为空。调用一个空接口值上的任意方法都会产生panic:

```go
w.Write([]byte("hello")) // panic: nil pointer dereference
```

第二个语句将一个*os.File类型的值赋给变量w:

```go
w = os.Stdout
```

这个赋值过程调用了一个具体类型到接口类型的隐式转换，这和显式的使用io.Writer(os.Stdout)是等价的。这类转换不管是显式的还是隐式的，都会刻画出操作到的类型和值。这个接口值的动态类型被设为`*os.File`指针的类型描述符（os.Stdout 是指向 os.File 的指针），它的动态值持有os.Stdout的拷贝；这是一个指向处理标准输出的os.File类型变量的指针。

![img](https://cdn.learnku.com/uploads/images/201912/28/6964/COeLNISHi3.png!large)

调用一个包含`*os.File`类型指针的接口值的Write方法，使得`(*os.File).Write`方法被调用。这个调用输出“hello”。

```go
w.Write([]byte("hello")) // "hello"
```

第三个语句给接口值赋了一个*bytes.Buffer类型的值

```go
w = new(bytes.Buffer)
```

现在动态类型是*bytes.Buffer并且动态值是一个指向新分配的缓冲区的指针（图7.3）。

![img](https://cdn.learnku.com/uploads/images/201912/28/6964/jJf6DrVptY.png!large)

Write方法的调用也使用了和之前一样的机制：

```go
w.Write([]byte("hello")) // writes "hello" to the bytes.Buffers
```

这次类型描述符是`*bytes.Buffer`，所以调用了`(*bytes.Buffer).Write`方法，并且接收者是该缓冲区的地址。这个调用把字符串“hello”添加到缓冲区中。

最后，第四个语句将nil赋给了接口值：

```go
w = nil
```

这个重置将它所有的部分都设为nil值，把变量w恢复到和它之前定义时相同的状态图，在图7.1中可以看到。

### 一个包含nil指针的接口不是nil接口

一个不包含任何值的nil接口值和一个刚好包含nil指针的接口值是不同的。这个细微区别产生了一个容易绊倒每个Go程序员的陷阱。

思考下面的程序。当debug变量设置为true时，main函数会将f函数的输出收集到一个bytes.Buffer类型中。

```go
const debug = true

func main() {
    var buf *bytes.Buffer
    if debug {
        buf = new(bytes.Buffer) // enable collection of output
    }
    f(buf) // NOTE: subtly incorrect!
    if debug {
        // ...use buf...
    }
}

// If out is non-nil, output will be written to it.
func f(out io.Writer) {
    // ...do something...
    if out != nil {
        out.Write([]byte("done!\n"))
    }
}
```

我们可能会预计当把变量debug设置为false时可以禁止对输出的收集，但是实际上在out.Write方法调用时程序发生了panic：

```go
if out != nil {
    out.Write([]byte("done!\n")) // panic: nil pointer dereference
}
```

当main函数调用函数f时，它给f函数的out参数赋了一个`*bytes.Buffer`的空指针，所以out的动值是nil。然而，它的动态类型是`*bytes.Buffer`，意思就是out变量是一个包含空指针值的非空接口（如图7.5），所以防御性检查out!=nil的结果依然是true。

![img](https://cdn.learnku.com/uploads/images/201912/28/6964/WlEzgXBzjA.png!large)

动态分配机制依然决定`(*bytes.Buffer).Write`的方法会被调用，但是这次的接收者的值是nil。对于一些如`*os.File`的类型，nil是一个有效的接收者(§6.2.1)，但是`*bytes.Buffer`类型不在这些类型中。这个方法会被调用，但是当它尝试去获取缓冲区时会发生panic。

问题在于尽管一个nil的`*bytes.Buffer`指针有实现这个接口的方法，它也不满足这个接口具体的行为上的要求。特别是这个调用违反了`(*bytes.Buffer).Write`方法的接收者非空的隐含先觉条件，所以将nil指针赋给这个接口是错误的。解决方案就是将main函数中的变量buf声明的类型改为io.Writer，（它的零值动态类型和动态值都为 nil）因此可以避免一开始就将一个不完全的值赋值给这个接口：

```go
var buf io.Writer
if debug {
    buf = new(bytes.Buffer) // enable collection of output
}
f(buf) // OK
```



### error 接口

- 预定义的error类型实际上就是interface类型，这个类型有一个返回错误信息的单一方法：

 ```go
  type error interface {
      Error() string
  }
 ```

- 创建一个error最简单的方法就是调用errors.New函数，它会根据传入的错误信息返回一个新的error。整个errors包仅只有4行：

```go
package errors

func New(text string) error { return &errorString{text} }

type errorString struct { text string }

func (e *errorString) Error() string { return e.text }
```

每个New函数的调用都分配了一个独特的和其他错误不相同的实例。我们也不想要重要的`error`例如`io.EOF`和一个刚好有相同错误消息的`error`比较后相等。

```go
fmt.Println(errors.New("EOF") == errors.New("EOF")) // "false"
```

调用`errors.New`函数是非常稀少的，因为有一个方便的封装函数`fmt.Errorf`，它还会处理字符串格式化。

```go
package fmt

import "errors"

func Errorf(format string, args ...interface{}) error {
    return errors.New(Sprintf(format, args...))
}
```

### 类型断言

- 类型断言是一个使用在接口值上的操作。语法上它看起来像x.(T)被称为断言类型。这里x表示一个接口值，T表示一个类型（接口类型或者具体类型）。一个类型断言会检查操作对象的动态类型是否和断言类型匹配。

- x.(T)中如果断言的类型T是一个具体类型，类型断言检查x的动态类型是否和T相同。如果是，类型断言的结果是x的动态值，当然它的类型是T。换句话说，具体类型的类型断言从它的操作对象中获得具体的值。如果x 的动态类型与 T 不相同，会抛出panic。

```go
var w io.Writer
w = os.Stdout
f := w.(*os.File)      // success: f == os.Stdout
c := w.(*bytes.Buffer) // panic: interface holds *os.File, not *bytes.Buffer
```

- 相反断言的类型T是一个接口类型，然后类型断言检查是否x的动态类型满足T。如果这个检查成功了，这个结果仍然是一个有相同类型和值部分的接口值，但是结果接口值的动态类型为T。换句话说，对一个接口类型的类型断言改变了类型的表述方式，改变了可以获取的方法集合（通常更大），但是它保护了接口值内部的动态类型和值的部分。

- 在下面的第一个类型断言后，w和rw都持有os.Stdout因为它们每个值的动态类型都是`*os.File`，但是变量的类型是io.Writer只对外公开出文件的Write方法，变量rw的类型为 io.ReadWriter，只对外公开文件的Read方法。

```go
var w io.Writer
w = os.Stdout
rw := w.(io.ReadWriter) // success: *os.File has both Read and Write
w = new(ByteCounter)
rw = w.(io.ReadWriter) // panic: *ByteCounter has no Read method
```

- 如果断言操作的对象是一个nil接口值，那么不论被断言的类型是什么这个类型断言都会失败。
- 经常地我们对一个接口值的动态类型是不确定的，并且我们更愿意去检验它是否是一些特定的类型。如果类型断言出现在一个有两个结果的赋值表达式中，例如如下的定义，这个类型断言不会在失败的时候发生panic，代替地返回的第二个返回值是一个标识类型断言是否成功的布尔值：

```go
var w io.Writer = os.Stdout
f, ok := w.(*os.File)      // success:  ok, f == os.Stdout
b, ok := w.(*bytes.Buffer) // failure: !ok, b == nil
```

### type switch

接口被以两种不同的方式使用。在第一个方式中，以io.Reader，io.Writer，fmt.Stringer，sort.Interface，http.Handler，和error为典型，一个接口的方法表达了实现这个接口的具体类型间的相似性，但是隐藏了代表的细节和这些具体类型本身的操作。重点在于方法上，而不是具体的类型上。

第二个方式利用一个接口值可以持有各种具体类型值的能力并且将这个接口认为是这些类型的union（联合）。类型断言用来动态地区别这些类型。在这个方式中，重点在于具体的类型满足这个接口，而不是在于接口的方法（如果它确实有一些的话），并且没有任何的信息隐藏。我们将以这种方式使用的接口描述为discriminated unions（可辨识联合）。

一个类型开关像普通的switch语句一样，它的运算对象是x.(type)－它使用了关键词字面量type－并且每个case有一到多个类型。一个类型开关基于这个接口值的动态类型使一个多路分支有效。这个nil的case和if x == nil匹配，并且这个default的case和如果其它case都不匹配的情况匹配。一个对sqlQuote的类型开关可能会有这些case

```go
switch x.(type) {
    case nil:       // ...
    case int, uint: // ...
    case bool:      // ...
    case string:    // ...
    default:        // ...
}
```

类型开关语句有一个扩展的形式，它可以将提取的值绑定到一个在每个case范围内的新变量上。

```go
switch x := x.(type) { /* ... */ }
```

使用类型开关的扩展形式来重写sqlQuote函数会让这个函数更加的清晰：

```go
func sqlQuote(x interface{}) string {
    switch x := x.(type) {
    case nil:
        return "NULL"
    case int, uint:
        return fmt.Sprintf("%d", x) // x has type interface{} here.
    case bool:
        if x {
            return "TRUE"
        }
        return "FALSE"
    case string:
        return sqlQuoteString(x) // (not shown)
    default:
        panic(fmt.Sprintf("unexpected type %T: %v", x, x))
    }
}
```

尽管sqlQuote接受一个任意类型的参数，但是这个函数只会在它的参数匹配类型开关中的一个case时运行到结束；其它情况的它会panic出“unexpected type”消息。虽然x的类型是interface{}，但是我们把它认为是一个int，uint，bool，string，和nil值的discriminated union（可识别联合）

### 使用建议

- 接口只有当有两个或两个以上的具体类型必须以相同的方式进行处理时才需要。
- 当一个接口只被一个单一的具体类型实现时有一个例外，就是由于它的依赖，这个具体类型不能和这个接口存在在一个相同的包中。这种情况下，一个接口是解耦这两个包的一个好的方式。

- 因为在Go语言中只有当两个或更多的类型须以相同的方式进行处理时才有必要使用接口，它们必定会从任意特定的实现细节中抽象出来。结果就是有更少和更简单方法（经常和io.Writer或 fmt.Stringer一样只有一个）的更小的接口。当新的类型出现时，小的接口更容易满足。对于接口设计的一个好的标准就是 ask only for what you need（只考虑你需要的东西）。
