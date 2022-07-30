> 以下内容为几年前学Go语言时记录的，发布在这里让这个项目对更多人能起到参考作用。
更多时下流行技术的应用和实战教程，可通过我公众号「网管叨bi叨」每周的推文来学习

## 函数

### 函数声明

- 函数声明包括函数名、形式参数列表、返回值列表（可省略）以及函数体。

```Go
	func name(parameter-list) (result-list) {
  	  body
	}
```
  形式参数列表描述了函数的参数名以及参数类型。这些参数作为局部变量，其值由参数调用者提供。返回值也可以像形式参数一样被命名，在这种情况下，每个返回值被声明成一个局部变量，并初始化为其类型的零值。

- 用 _ 符号作为形参名可以强调某个参数未被使用。

```Go
  func first(x int, _ int) int { return x }
```

- 函数的类型被称为函数的标识符。如果两个函数形式参数列表和返回值列表中的变量类型一一对应，那么这两个函数被认为有相同的类型和标识符。
- 在函数调用时，Go语言没有默认参数值，也没有任何方法可以通过参数名指定形参，因此形参和返回值的变量名对于函数调用者而言没有意义。
- 实参通过值的方式传递，因此函数的形参是实参的拷贝。对形参进行修改不会影响实参。但是，如果实参包括引用类型，如指针，slice(切片)、map、function、channel等类型，实参可能会由于函数的引用而被修改。
- golang.org/x/... 目录下存储了一些由Go团队设计、维护，对网络编程、国际化文件处理、移动平台、图像处理、加密解密、开发者工具提供支持的扩展包。未将这些扩展包加入到标准库原因有二，一是部分包仍在开发中，二是对大多数Go语言的开发者而言，扩展包提供的功能很少被使用。

### 递归调用

- 大部分编程语言使用固定大小的函数调用栈，常见的大小从64KB到2MB不等。固定大小栈会限制递归的深度，当你用递归处理大量数据时，需要避免栈溢出；除此之外，还会导致安全性问题。与相反,Go语言使用可变栈，栈的大小按需增加(初始时很小)。这使得我们使用递归时不必考虑溢出和安全问题
- 虽然Go的垃圾回收机制会回收不被使用的内存，但是这不包括操作系统层面的资源，比如打开的文件、网络连接。因此我们必须显式的释放这些资源。

### 多返回值函数

- 调用多返回值函数时，返回给调用者的是一组值，调用者必须显式的将这些值分配给变量:

```Go
  links, err := findLinks(url)
```

  如果某个值不被使用，可以将其分配给blank identifier:

```Go
  links, _ := findLinks(url) // errors ignored
```

- 如果一个函数将所有的返回值都显示的变量名，那么该函数的return语句可以省略操作数。这称之为bare return。

```Go
  // CountWordsAndImages does an HTTP GET request for the HTML
  // document url and returns the number of words and images in it.
  func CountWordsAndImages(url string) (words, images int, err error) {
      resp, err := http.Get(url)
      if err != nil {
          return
      }
      doc, err := html.Parse(resp.Body)
      resp.Body.Close()
      if err != nil {
          err = fmt.Errorf("parsing HTML: %s", err)
      return
      }
      words, images = countWordsAndImages(doc)
      return
  }
  func countWordsAndImages(n *html.Node) (words, images int) { /* ... */ }
```

  按照函数声明中返回值列表的次序，返回所有的返回值，在上面的例子中，每一个return语句等价于：

```Go
  return words, images, err
```

- 当一个函数有多处return语句以及许多返回值时，bare return 可以减少代码的重复，但是使得代码难以被理解。如果你没有仔细的审查上面的代码，很难发现前2处return等价于 `return 0,0,err`（Go会将返回值 words和images在函数体的开始处，根据它们的类型，将其初始化为0），最后一处return等价于 `return words，image，nil`。基于以上原因，不宜过度使用bare return。

### 错误

- 在Go的错误处理中，错误是软件包API和应用程序用户界面的一个重要组成部分，程序运行失败仅被认为是几个预期的结果之一。

- 对于那些将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，通常是最后一个，来传递错误信息。

```go
  resp, err := http.Get(url)
```

- 内置的error是接口类型。nil意味着函数运行成功，non-nil表示失败。对于non-nil的error类型,我们可以通过调用error的`Error`函数或者输出函数获得字符串类型的错误信息。

```Go
  fmt.Println(err)
  fmt.Printf("%v", err)
```

### 函数值

- 在Go中，函数被看作第一类值（first-class values）：函数像其他值一样，拥有类型，可以被赋值给其他变量，传递给函数，从函数返回。对函数值（function value）的调用类似函数调用。例子如下：

```Go
  func square(n int) int { return n * n }
  func negative(n int) int { return -n }
  func product(m, n int) int { return m * n }
  
  f := square
  fmt.Println(f(3)) // "9"
  
  f = negative
  fmt.Println(f(3))     // "-3"
  fmt.Printf("%T\n", f) // "func(int) int"
  
  f = product // compile error: can't assign func(int, int) int to func(int) int
```

- 函数类型的零值是nil。调用值为nil的函数值会引起panic错误：

```Go
  var f func(int) int
  f(3) // 此处f的值为nil, 会引起panic错误
```

- 函数值可以与nil比较：

```Go
  var f func(int) int
  if f != nil {
     f(3)
  }
```

  但是函数值之间是不可比较的，也不能用函数值作为map的key。



### 匿名函数

- 拥有函数名的函数只能在包级语法块中被声明，通过函数字面量（function literal），我们可绕过这一限制，在任何表达式中表示一个函数值。函数字面量的语法和函数声明相似，区别在于func关键字后没有函数名。函数值字面量是一种表达式，它的值被称为匿名函数（anonymous function）。

  函数字面量允许我们在使用函数时，再定义它。通过这种技巧，我们可以改写之前对strings.Map的调用：

```Go
  strings.Map(func(r rune) rune { return r + 1 }, "HAL-9000")
```

  更为重要的是，通过这种方式定义的函数可以访问完整的词法环境（lexical environment），这意味着在函数中定义的内部函数可以引用该函数的变量。

```Go
  // squares返回一个匿名函数。
  // 该匿名函数每次被调用时都会返回下一个数的平方。
  func squares() func() int {
      var x int
      return func() int {
          x++
          return x * x
      }
  }
```

  通过这个例子，我们看到变量的生命周期不由它的作用域决定：squares返回后，变量x仍然隐式的存在于f中。

- 当匿名函数需要被递归调用时，我们必须首先声明一个变量，再将匿名函数赋值给这个变量。如果不分成两步，函数字面量无法与变量绑定，我们也无法递归调用该匿名函数，比如:

```Go
  var visitAll func(items []string)
  visitAll = func(items []string) {
      ......  
      visitAll(m[item])
      ......
  }
  
```
	否则会出现编译错误
```Go
  visitAll := func(items []string) {
      // ...
      visitAll(m[item]) // compile error: undefined:   visitAll
      // ...
  }
```

### 可变参数

- 参数数量可变的函数称为为可变参数函数。典型的例子就是fmt.Printf和类似函数。Printf首先接收一个的必备参数，之后接收任意个数的后续参数。

  在声明可变参数函数时，需要在参数列表的最后一个参数类型之前加上省略符号“...”，这表示该函数会接收任意数量的该类型参数。

```Go
  func sum(vals...int) int {
      total := 0
      for _, val := range vals {
          total += val
      }
      return total
  }
```

  sum函数返回任意个int型参数的和。在函数体中,vals被看作是类型为`[] int`的切片。sum可以接收任意数量的int型参数：

```Go
  fmt.Println(sum())           // "0"
  fmt.Println(sum(3))          // "3"
  fmt.Println(sum(1, 2, 3, 4)) // "10"
```

- 在上面的代码中，调用者隐式的创建一个数组，并将原始参数复制到数组中，再把数组的一个切片作为参数传给被调函数。如果原始参数已经是切片类型，我们该如何传递给sum？只需在最后一个参数后加上省略符。下面的代码功能与上个例子中最后一条语句相同。

```Go
  values := []int{1, 2, 3, 4}
  fmt.Println(sum(values...)) // "10"
  // fmt.Println(sum(1, 2, 3, 4))
```

- 虽然在可变参数函数内部，`...int` 型参数的行为看起来很像切片类型，但实际上，可变参数函数和以切片作为参数的函数是不同的。

```Go
  func f(...int) {}
  func g([]int) {}
  fmt.Printf("%T\n", f) // "func(...int)"
  fmt.Printf("%T\n", g) // "func([]int)"
```

- 可变参数函数经常被用于格式化字符串。下面的errorf函数构造了一个以行号开头的，经过格式化的错误信息。函数名的后缀f是一种通用的命名规范，代表该可变参数函数可以接收Printf风格的格式化字符串。

```Go
  func errorf(linenum int, format string, args ...interface{}) {
      fmt.Fprintf(os.Stderr, "Line %d: ", linenum)
      fmt.Fprintf(os.Stderr, format, args...)
      fmt.Fprintln(os.Stderr)
  }
  linenum, name := 12, "count"
  errorf(linenum, "undefined: %s", name) // "Line 12: undefined: count"
```

  `...interfac{}`表示函数在`format`参数后可以接收任意个任意类型的参数。`interface{}`会在后面介绍。

  

###  Deferred 函数

- 你只需要在调用普通函数或方法前加上关键字defer，就完成了defer所需要的语法。当defer语句被执行时，跟在defer后面的函数会被延迟执行。直到包含该defer语句的函数执行完毕时，defer后的函数才会被执行，不论包含defer语句的函数是通过return正常结束，还是由于panic导致的异常结束。你可以在一个函数中执行多条defer语句，它们的执行顺序与声明顺序相反。

- defer语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。通过defer机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。释放资源的defer应该直接跟在请求资源的语句后。

- 对文件的操作

```Go
  package ioutil
  func ReadFile(filename string) ([]byte, error) {
      f, err := os.Open(filename)
      if err != nil {
          return nil, err
      }
      defer f.Close()
      return ReadAll(f)
  }
```

- 处理互斥锁

```Go
  var mu sync.Mutex
  var m = make(map[string]int)
  func lookup(key string) int {
      mu.Lock()
      defer mu.Unlock()
      return m[key]
  }
```

- 调试复杂程序时，defer机制也常被用于记录何时进入和退出函数。下例中的bigSlowOperation函数，直接调用trace记录函数的被调情况。bigSlowOperation被调时，trace会返回一个函数值，该函数值会在bigSlowOperation退出时被调用。通过这种方式， 我们可以只通过一条语句控制函数的入口和所有的出口，甚至可以记录函数的运行时间，如例子中的start。需要注意一点：不要忘记defer语句后的圆括号，否则本该在进入时执行的操作会在退出时执行，而本该在退出时执行的，永远不会被执行。

  *gopl.io/ch5/trace*

```Go
  func bigSlowOperation() {
      defer trace("bigSlowOperation")() // don't forget the extra parentheses
      // ...lots of work…
      time.Sleep(10 * time.Second) // simulate slow
      operation by sleeping
  }
  func trace(msg string) func() {
      start := time.Now()
      log.Printf("enter %s", msg)
      return func() { 
          log.Printf("exit %s (%s)", msg,time.Since(start)) 
      }
  }
```

  每一次bigSlowOperation被调用，程序都会记录函数的进入，退出，持续时间。（我们用time.Sleep模拟一个耗时的操作）

```
  $ go build gopl.io/ch5/trace
  $ ./trace
  2015/11/18 09:53:26 enter bigSlowOperation
  2015/11/18 09:53:36 exit bigSlowOperation (10.000589217s)
```

- 用 defer 函数记录返回值(需要是命名返回值才能记录)

```Go
  func double(x int) (result int) {
      defer func() { fmt.Printf("double(%d) = %d\n", x,result) }()
      return x + x
  }
  _ = double(4)
  // Output:
  // "double(4) = 8"
```

- 被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值：

```Go
  func triple(x int) (result int) {
      defer func() { result += x }()
      return double(x)
  }
  fmt.Println(triple(4)) // "12"
```

- 在循环体中的defer语句需要特别注意，因为只有在函数执行完毕后，这些被延迟的函数才会执行。下面的代码会导致系统的文件描述符耗尽，因为在所有文件都被处理之前，没有文件会被关闭。

```Go
  for _, filename := range filenames {
      f, err := os.Open(filename)
      if err != nil {
          return err
      }
      defer f.Close() // NOTE: risky; could run out of file
      descriptors
      // ...process f…
  }
```

  一种解决方法是将循环体中的文件操作和defer语句移至另外一个函数。在每次循环时，调用这个函数。

```Go
  for _, filename := range filenames {
      if err := doFile(filename); err != nil {
          return err
      }
  }
  func doFile(filename string) error {
      f, err := os.Open(filename)
      if err != nil {
          return err
      }
      defer f.Close()
      // ...process f…
  }
```

### Panic 和 Recover

- Go的类型系统会在编译时捕获很多错误，但有些错误只能在运行时检查，如数组访问越界、空指针引用等。这些运行时错误会引起painc异常。

- 当panic异常发生时，程序会中断运行，并立即执行在该goroutine（可以先理解成线程，在第8章会详细介绍）中被延迟的函数（defer 机制）。随后，程序崩溃并输出日志信息。日志信息包括panic value和函数调用的堆栈跟踪信息。

- 虽然Go的panic机制类似于其他语言的异常，但panic的适用场景有一些不同。由于panic会引起程序的崩溃，因此panic一般用于严重错误，如程序内部的逻辑不一致。

- 通常来说，不应该对panic异常做任何处理，但有时，也许我们可以从异常中恢复，至少我们可以在程序崩溃前，做一些操作。举个例子，当web服务器遇到不可预料的严重问题时，在崩溃前应该将所有的连接关闭；如果不做任何处理，会使得客户端一直处于等待状态。

- 如果在deferred函数中调用了内置函数recover，并且定义该defer语句的函数发生了panic异常，recover会使程序从panic中恢复，并返回panic value。导致panic异常的函数不会继续运行，但能正常返回。在未发生panic时调用recover，recover会返回nil。

- 例子中deferred函数帮助Parse从panic中恢复。在deferred函数内部，panic value被附加到错误信息中；并用err变量接收错误信息，返回给调用者。

```Go
  func Parse(input string) (s *Syntax, err error) {
      defer func() {
          if p := recover(); p != nil {
              err = fmt.Errorf("internal error: %v", p)
          }
      }()
      // ...parser...
  }
```

- 不加区分的恢复所有的panic异常，不是可取的做法。

- 只恢复应该被恢复的panic异常，此外，这些异常所占的比例应该尽可能的低。为了标识某个panic是否应该被恢复，我们可以将panic value设置成特殊类型。在recover时对panic value进行检查，如果发现panic value是特殊类型，就将这个panic作为errror处理，如果不是，则按照正常的panic进行处理

```Go
  func soleTitle(doc *html.Node) (title string, err error) {
      type bailout struct{}
      defer func() {
          switch p := recover(); p {
          case nil:       // no panic
          case bailout{}: // "expected" panic
              err = fmt.Errorf("multiple title elements")
          default:
              panic(p)
          }
      }()
      forEachNode(doc, func(n *html.Node) {
          if n.Type == html.ElementNode && n.Data == "title" &&
              n.FirstChild != nil {
              if title != "" {
                  panic(bailout{}) // multiple titleelements
              }
              title = n.FirstChild.Data
          }
      }, nil)
      if title == "" {
          return "", fmt.Errorf("no title element")
      }
      return title, nil
  }
```
