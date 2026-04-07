# 表达式转换规则

> 将 JSX 表达式 `{expr}` 转换为有效的 Go 代码

## 核心挑战

```
┌─────────────────────────────────────────────────────────────┐
│              表达式转换的本质问题                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX 表达式: {expr}                                         │
│                                                             │
│  JavaScript 上下文:                                         │
│  ├── expr 是 JS 表达式                                      │
│  ├── 运行时求值                                             │
│  ├── 类型动态                                               │
│  └── 可以是任意 JS 值                                       │
│                                                             │
│  Go 上下文:                                                 │
│  ├── 需要生成有效的 Go 表达式                               │
│  ├── 编译时类型检查                                         │
│  ├── 静态类型                                               │
│  └── 需要明确的类型转换                                     │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 表达式分类

### 类别 A: 简单值引用

```
┌─────────────────────────────────────────────────────────────┐
│ 变量引用                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {name}          → Go:      name                    │
│ JSX:     {count}         → Go:      count                   │
│                                                             │
│ 属性访问                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {props.title}   → Go:      props.Title             │
│ JSX:     {user.profile.name} → Go:  user.Profile.Name       │
│ JSX:     {items.length}  → Go:      len(items)              │
│                                                             │
│ 转换规则:                                                   │
│ - .property → .Property (首字母大写)                        │
│ - .length → len() (特殊处理)                                │
│                                                             │
│ 索引访问                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {items[0]}      → Go:      items[0]                │
│ JSX:     {items[index]}  → Go:      items[index]            │
│ JSX:     {obj["key"]}    → Go:      obj["key"]              │
└─────────────────────────────────────────────────────────────┘
```

### 类别 B: 字面量

```
┌─────────────────────────────────────────────────────────────┐
│ 数字字面量                                                  │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {42}            → Go:      42                      │
│ JSX:     {3.14}          → Go:      3.14                    │
│ JSX:     {0xFF}          → Go:      0xFF                    │
│                                                             │
│ 字符串字面量                                                │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {"hello"}       → Go:      "hello"                 │
│ JSX:     {'a'}           → Go:      'a' (rune)              │
│                                                             │
│ 布尔字面量                                                  │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {true}          → Go:      true                    │
│ JSX:     {false}         → Go:      false                   │
│                                                             │
│ 空值                                                        │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {null}          → Go:      nil                     │
│ JSX:     {undefined}     → Go:      nil                     │
└─────────────────────────────────────────────────────────────┘
```

### 类别 C: 运算表达式

```
┌─────────────────────────────────────────────────────────────┐
│ 算术运算                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {a + b}         → Go:      a + b                   │
│ JSX:     {a - b}         → Go:      a - b                   │
│ JSX:     {a * b}         → Go:      a * b                   │
│ JSX:     {a / b}         → Go:      a / b                   │
│ JSX:     {a % b}         → Go:      a % b                   │
│                                                             │
│ 比较运算                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {a === b}       → Go:      a == b                  │
│ JSX:     {a !== b}       → Go:      a != b                  │
│ JSX:     {a > b}         → Go:      a > b                   │
│ JSX:     {a >= b}        → Go:      a >= b                  │
│                                                             │
│ 逻辑运算                                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {a && b}        → Go:      a && b                  │
│ JSX:     {a || b}        → Go:      a || b                  │
│ JSX:     {!a}            → Go:      !a                      │
│                                                             │
│ 三元运算 (Go 无原生支持)                                    │
├─────────────────────────────────────────────────────────────┤
│ JSX:     {cond ? a : b}  → Go:      ink.Ternary(cond, a, b) │
└─────────────────────────────────────────────────────────────┘
```

### 类别 D: 函数调用

```
┌─────────────────────────────────────────────────────────────┐
│ 内置函数映射表                                              │
├──────────────────────┬──────────────────────────────────────┤
│ JSX/JS               │ Go                                   │
├──────────────────────┼──────────────────────────────────────┤
│ Math.max(a, b)       │ max(a, b)                            │
│ Math.min(a, b)       │ min(a, b)                            │
│ Math.abs(a)          │ abs(a)                               │
│ Math.floor(a)        │ int(a)                               │
│ Math.ceil(a)         │ int(math.Ceil(a))                    │
│ Math.round(a)        │ int(math.Round(a))                   │
│ parseInt(s)          │ strconv.Atoi(s)                      │
│ parseFloat(s)        │ strconv.ParseFloat(s, 64)            │
│ String(n)            │ strconv.Itoa(n)                      │
│ isNaN(n)             │ math.IsNaN(n)                        │
└──────────────────────┴──────────────────────────────────────┘
```

### 字符串方法映射

```
┌──────────────────────┬──────────────────────────────────────┐
│ JSX/JS               │ Go                                   │
├──────────────────────┼──────────────────────────────────────┤
│ s.toUpperCase()      │ strings.ToUpper(s)                   │
│ s.toLowerCase()      │ strings.ToLower(s)                   │
│ s.trim()             │ strings.TrimSpace(s)                 │
│ s.trimStart()        │ strings.TrimLeft(s, " \t")           │
│ s.trimEnd()          │ strings.TrimRight(s, " \t")          │
│ s.startsWith(x)      │ strings.HasPrefix(s, x)              │
│ s.endsWith(x)        │ strings.HasSuffix(s, x)              │
│ s.includes(x)        │ strings.Contains(s, x)               │
│ s.indexOf(x)         │ strings.Index(s, x)                  │
│ s.replace(a, b)      │ strings.Replace(s, a, b, 1)          │
│ s.replaceAll(a, b)   │ strings.ReplaceAll(s, a, b)          │
│ s.split(sep)         │ strings.Split(s, sep)                │
│ s.repeat(n)          │ strings.Repeat(s, n)                 │
│ s.length             │ len(s)                               │
└──────────────────────┴──────────────────────────────────────┘
```

### 数组方法映射

```
┌──────────────────────┬──────────────────────────────────────┐
│ JSX/JS               │ Go                                   │
├──────────────────────┼──────────────────────────────────────┤
│ arr.length           │ len(arr)                             │
│ arr.push(x)          │ arr = append(arr, x)                 │
│ arr.slice(i, j)      │ arr[i:j]                             │
│ arr.concat(b)        │ append(arr, b...)                    │
│ arr.includes(x)      │ slices.Contains(arr, x)              │
│ arr.indexOf(x)       │ slices.Index(arr, x)                 │
│ arr.reverse()        │ slices.Reverse(arr)                  │
│ arr.sort()           │ slices.Sort(arr)                     │
│ arr.join(sep)        │ strings.Join(arr, sep)               │
├──────────────────────┴──────────────────────────────────────┤
│ 需要特殊处理 (转换为 Go 循环):                               │
│ arr.map(fn)          │ ink.Map(arr, fn)                     │
│ arr.filter(fn)       │ 需要展开循环                         │
│ arr.reduce(fn, init) │ 需要展开循环                         │
│ arr.forEach(fn)      │ 需要展开循环                         │
│ arr.find(fn)         │ 需要展开循环                         │
│ arr.every(fn)        │ 需要展开循环                         │
│ arr.some(fn)         │ 需要展开循环                         │
└──────────────────────┴──────────────────────────────────────┘
```

---

## 类别 E: JSX 生成表达式 (核心重点)

### 条件渲染: && 运算符

```
┌─────────────────────────────────────────────────────────────┐
│  模式: {condition && <Element/>}                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX 源码:                                                  │
│  <Box>                                                      │
│      {show && <Text color="green">Visible</Text>}           │
│  </Box>                                                     │
│                                                             │
│  Go 输出:                                                   │
│  ink.Box(ink.Props{},                                       │
│      ink.Cond(show,                                         │
│          ink.Text(ink.Props{Color: ink.Green},              │
│              ink.TextContent("Visible"),                    │
│          ),                                                 │
│          nil,                                               │
│      ),                                                     │
│  )                                                          │
│                                                             │
│  ink.Cond 辅助函数:                                         │
│  func Cond(cond bool, ifTrue, ifFalse Element) Element {    │
│      if cond {                                              │
│          return ifTrue                                      │
│      }                                                      │
│      return ifFalse                                         │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

### 条件渲染: 三元运算符

```
┌─────────────────────────────────────────────────────────────┐
│  模式: {condition ? <A/> : <B/>}                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX 源码:                                                  │
│  {status === "loading" ? (                                  │
│      <Text>Loading...</Text>                                │
│  ) : (                                                      │
│      <Text>Ready</Text>                                     │
│  )}                                                         │
│                                                             │
│  Go 输出:                                                   │
│  ink.Cond(status == "loading",                              │
│      ink.Text(ink.Props{},                                  │
│          ink.TextContent("Loading..."),                     │
│      ),                                                     │
│      ink.Text(ink.Props{},                                  │
│          ink.TextContent("Ready"),                          │
│      ),                                                     │
│  )                                                          │
└─────────────────────────────────────────────────────────────┘
```

### 列表渲染: map

```
┌─────────────────────────────────────────────────────────────┐
│  模式: {items.map(item => <Item key={item.id} {...item}/>)} │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX 源码:                                                  │
│  <Box flexDirection="column">                               │
│      {items.map((item, index) => (                          │
│          <Text key={index} color={item.color}>              │
│              {item.name}                                    │
│          </Text>                                            │
│      ))}                                                    │
│  </Box>                                                     │
│                                                             │
│  Go 输出:                                                   │
│  ink.Box(ink.Props{FlexDirection: ink.Column},              │
│      ink.Map(items, func(item Item, index int) ink.Element {│
│          return ink.Text(ink.Props{                         │
│              Key:   index,                                  │
│              Color: item.Color,                             │
│          },                                                 │
│              ink.TextContent(item.Name),                    │
│          )                                                  │
│      }),                                                    │
│  )                                                          │
│                                                             │
│  ink.Map 辅助函数:                                          │
│  func Map[T any](items []T, fn func(T, int) Element) []Element {│
│      result := make([]Element, len(items))                  │
│      for i, item := range items {                           │
│          result[i] = fn(item, i)                            │
│      }                                                      │
│      return result                                          │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

### 列表渲染: filter + map

```
┌─────────────────────────────────────────────────────────────┐
│  模式: {items.filter(fn).map(fn)}                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX 源码:                                                  │
│  {items                                                     │
│      .filter(item => item.active)                           │
│      .map(item => <Text>{item.name}</Text>)                 │
│  }                                                          │
│                                                             │
│  Go 输出:                                                   │
│  ink.ChildrenFunc(func() []ink.Element {                    │
│      var children []ink.Element                             │
│      for _, item := range items {                           │
│          if !item.Active {                                  │
│              continue  // filter                            │
│          }                                                  │
│          children = append(children,                        │
│              ink.Text(ink.Props{},                          │
│                  ink.TextContent(item.Name),                │
│              ),                                             │
│          )                                                  │
│      }                                                      │
│      return children                                        │
│  })                                                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 类别 F: 复合表达式

### 模板字符串

```
┌─────────────────────────────────────────────────────────────┐
│  模板字符串转换                                              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX:     {`Hello ${name}, count: ${count}`}                │
│                                                             │
│  Go:      fmt.Sprintf("Hello %s, count: %d", name, count)   │
│                                                             │
│  转换规则:                                                  │
│  - `text ${expr}` → fmt.Sprintf                             │
│  - ${expr} → %v (通用) 或根据类型推断                       │
│    - string → %s                                            │
│    - int → %d                                               │
│    - float → %f                                             │
│    - bool → %t                                              │
│    - other → %v                                             │
└─────────────────────────────────────────────────────────────┘
```

### 对象展开

```
┌─────────────────────────────────────────────────────────────┐
│  属性展开转换                                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  JSX:     <Box {...baseProps} padding={10}/>                │
│                                                             │
│  Go:      ink.Box(ink.MergeProps(baseProps, ink.Props{      │
│              Padding: 10,                                   │
│          }),                                                │
│      )                                                      │
│                                                             │
│  ink.MergeProps 辅助函数:                                   │
│  func MergeProps(base Props, overrides Props) Props         │
└─────────────────────────────────────────────────────────────┘
```

---

## 上下文驱动转换

### 转换策略矩阵

```
┌─────────────────────────────────────────────────────────────┐
│  上下文决定转换方式                                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  上下文 1: 属性值                                           │
│  <Box padding={10}>                                         │
│          ↑                                                  │
│      需要类型匹配 Props.Padding (int)                       │
│      转换: 直接使用表达式                                   │
│                                                             │
│  上下文 2: 文本内容                                         │
│  <Text>Hello {name}</Text>                                  │
│             ↑                                               │
│      需要转换为 string                                      │
│      转换: ink.ToString(expr) 或 fmt.Sprintf                │
│                                                             │
│  上下文 3: 子节点                                           │
│  {show && <Text>A</Text>}                                   │
│  ↑                                                          │
│      需要返回 ink.Element 或 []ink.Element                  │
│      转换: 条件表达式 → Go if 或 ink.Cond                   │
└─────────────────────────────────────────────────────────────┘
```

---

## 类型推断策略

```
┌─────────────────────────────────────────────────────────────┐
│  类型推断                                                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  策略 A: 从上下文推断                                       │
│  type BoxProps struct {                                     │
│      Padding    int      // 已知类型                        │
│      Color      string   // 已知类型                        │
│  }                                                          │
│  <Box padding={expr}>  // expr 必须是 int                   │
│                                                             │
│  策略 B: 从字面量推断                                       │
│  {42}      → int                                           │
│  {3.14}    → float64                                       │
│  {"hello"} → string                                        │
│  {true}    → bool                                          │
│                                                             │
│  策略 C: 文本内容统一转 string                              │
│  <Text>{count}</Text>  → ink.Textf(ink.Props{}, "%d", count)│
│                                                             │
│  ink.ToString 辅助函数:                                     │
│  func ToString(v any) string {                              │
│      switch val := v.(type) {                               │
│      case string:  return val                               │
│      case int:     return strconv.Itoa(val)                 │
│      case float64: return strconv.FormatFloat(val, 'f', -1, 64)│
│      case bool:    return strconv.FormatBool(val)           │
│      default:      return fmt.Sprintf("%v", val)            │
│      }                                                      │
│  }                                                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 辅助函数 API

```go
// 条件渲染
func Cond(cond bool, ifTrue, ifFalse Element) Element

// 列表渲染
func Map[T any](items []T, fn func(T, int) Element) []Element
func MapFilter[T any](items []T, fn func(T, int) Element) []Element

// 属性合并
func MergeProps(base Props, overrides Props) Props

// 文本内容
func TextContent(s string) Content
func ExprContent(v any) Content
func ToString(v any) string

// 子节点
func ChildrenFunc(fn func() []Element) Children

// Fragment
func Fragment(children ...Element) Element
```

---

## 总结

### 核心要点

1. **上下文驱动转换**
   - 属性值: 直接转换，类型由 Props 定义决定
   - 文本内容: 转换为 string，使用 ToString 或 Textf
   - 子节点: 转换为 Element 或 []Element

2. **条件渲染**
   - && → ink.Cond(expr, element, nil)
   - ? : → ink.Cond(cond, ifTrue, ifFalse)

3. **列表渲染**
   - .map() → ink.Map(items, fn)
   - .filter().map() → ink.MapFilter 或展开循环
   - 复杂链 → 展开为 Go 循环

4. **JS → Go 映射**
   - 运算符: === → ==, !== → !=
   - 方法: .toUpperCase() → strings.ToUpper()
   - 属性: .length → len()
   - 模板字符串 → fmt.Sprintf

5. **类型处理**
   - 从上下文推断
   - 从字面量推断
   - 文本内容统一转 string