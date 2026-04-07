# 终端渲染器设计

> 复刻 Ink 核心魔法：增量渲染 + 最小输出

## 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                  Terminal Renderer 架构                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   输入: 布局树 (来自 Layout Engine)                           │
│                                                             │
│   处理流程:                                                  │
│                                                             │
│   ┌─────────────┐                                          │
│   │ Layout Tree │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐                                          │
│   │  Build      │  构建输出缓冲                             │
│   │  Buffer     │                                          │
│   └─────────────┘                                          │
│         │                                                   │
│         ▼                                                   │
│   ┌─────────────┐    ┌─────────────┐                      │
│   │  Previous   │───▶│    Diff     │                      │
│   │  Buffer     │    │  (差异计算)  │                      │
│   └─────────────┘    └─────────────┘                      │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Patches    │                     │
│                        │  (补丁列表)  │                     │
│                        └─────────────┘                     │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Generate   │                     │
│                        │  ANSI       │                     │
│                        └─────────────┘                     │
│                              │                              │
│                              ▼                              │
│                        ┌─────────────┐                     │
│                        │  Output     │                     │
│                        │  to Stdout  │                     │
│                        └─────────────┘                     │
│                                                             │
│   核心: 最小化输出，只更新变化的部分                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 输出缓冲结构

```go
// 输出缓冲 (二维字符网格)
type OutputBuffer struct {
    // 尺寸
    Width  int
    Height int
    
    // 内容 (行 → 列 → 单元格)
    Lines [][]Cell
    
    // 光标位置
    CursorX int
    CursorY int
}

// 单元格
type Cell struct {
    Char    rune      // 字符
    Style   CellStyle // 样式
    
    // 元数据
    ZIndex  int       // 层级 (用于覆盖)
    Dirty   bool      // 是否变化
}

// 单元格样式
type CellStyle struct {
    // 前景色
    Foreground Color
    FgBold     bool
    FgDim      bool
    
    // 背景色
    Background Color
    
    // 装饰
    Underline  bool
    Strikethrough bool
    Inverse    bool  // 反色
}

// 颜色
type Color struct {
    Type ColorType
    
    // 基础色 (0-15)
    Basic int
    
    // 256色 (0-255)
    Extended int
    
    // 真彩色 (RGB)
    R, G, B uint8
}

type ColorType int

const (
    ColorDefault ColorType = iota
    ColorBasic
    ColorExtended
    ColorTrueColor
)
```

---

## 缓冲构建

```go
// 渲染器
type Renderer struct {
    // 当前缓冲
    current  *OutputBuffer
    
    // 上一次缓冲 (用于 Diff)
    previous *OutputBuffer
    
    // 终端接口
    terminal Terminal
    
    // 配置
    config   RenderConfig
}

// 构建输出缓冲
func (r *Renderer) BuildBuffer(root *LayoutNode) *OutputBuffer {
    // 1. 获取终端尺寸
    width, height := r.terminal.Size()
    
    // 2. 创建缓冲
    buffer := NewOutputBuffer(width, height)
    
    // 3. 递归渲染节点
    r.renderNode(buffer, root, 0, 0)
    
    return buffer
}

// 渲染节点到缓冲
func (r *Renderer) renderNode(buffer *OutputBuffer, node *LayoutNode, offsetX, offsetY int) {
    // 计算绝对位置
    x := int(node.Layout.X) + offsetX
    y := int(node.Layout.Y) + offsetY
    w := int(node.Layout.Width)
    h := int(node.Layout.Height)
    
    // 根据节点类型渲染
    switch node.Type {
    case NodeTypeBox:
        r.renderBox(buffer, node, x, y, w, h)
    case NodeTypeText:
        r.renderText(buffer, node, x, y)
    }
    
    // 递归渲染子节点
    for _, child := range node.Children {
        r.renderNode(buffer, child, x, y)
    }
}

// 渲染 Box
func (r *Renderer) renderBox(buffer *OutputBuffer, node *LayoutNode, x, y, w, h int) {
    style := node.Style
    
    // 绘制背景
    if style.BackgroundColor != ColorDefault {
        for row := y; row < y+h; row++ {
            for col := x; col < x+w; col++ {
                cell := buffer.GetCell(col, row)
                cell.Style.Background = style.BackgroundColor
                cell.Dirty = true
            }
        }
    }
    
    // 绘制边框
    if style.BorderStyle != BorderNone {
        r.renderBorder(buffer, node, x, y, w, h)
    }
}

// 渲染文本
func (r *Renderer) renderText(buffer *OutputBuffer, node *LayoutNode, x, y int) {
    text := node.Text
    style := node.Style
    
    col := x
    row := y
    
    for _, r := range text {
        // 处理换行
        if r == '\n' {
            row++
            col = x
            continue
        }
        
        // 写入单元格
        cell := buffer.GetCell(col, row)
        cell.Char = r
        cell.Style.Foreground = style.Color
        cell.Style.FgBold = style.Bold
        cell.Style.FgDim = style.Dim
        cell.Style.Underline = style.Underline
        cell.Style.Strikethrough = style.Strikethrough
        cell.Dirty = true
        
        // 更新列位置
        if r < 128 {
            col += 1
        } else {
            col += 2  // 宽字符
        }
    }
}
```

---

## 差异计算 (Diff)

```go
// 计算差异
func (r *Renderer) Diff(old, new *OutputBuffer) []Patch {
    patches := []Patch{}
    
    // 尺寸变化
    if old.Width != new.Width || old.Height != new.Height {
        patches = append(patches, Patch{
            Type: PatchResize,
            Width: new.Width,
            Height: new.Height,
        })
    }
    
    // 逐行比较
    for y := 0; y < new.Height; y++ {
        // 跳过未变化的行
        if y < old.Height && r.linesEqual(old, new, y) {
            continue
        }
        
        // 行变化，逐单元格比较
        for x := 0; x < new.Width; x++ {
            oldCell := r.getCell(old, x, y)
            newCell := new.Lines[y][x]
            
            if !r.cellsEqual(oldCell, newCell) {
                patches = append(patches, Patch{
                    Type: PatchCell,
                    X:    x,
                    Y:    y,
                    Cell: newCell,
                })
            }
        }
    }
    
    return patches
}

// 比较行是否相同
func (r *Renderer) linesEqual(old, new *OutputBuffer, y int) bool {
    if y >= len(old.Lines) || y >= len(new.Lines) {
        return false
    }
    
    oldLine := old.Lines[y]
    newLine := new.Lines[y]
    
    for x := 0; x < len(oldLine) && x < len(newLine); x++ {
        if !r.cellsEqual(oldLine[x], newLine[x]) {
            return false
        }
    }
    
    return true
}

// 比较单元格是否相同
func (r *Renderer) cellsEqual(a, b Cell) bool {
    return a.Char == b.Char &&
           a.Style.Foreground == b.Style.Foreground &&
           a.Style.Background == b.Style.Background &&
           a.Style.FgBold == b.Style.FgBold &&
           a.Style.FgDim == b.Style.FgDim &&
           a.Style.Underline == b.Style.Underline
}
```

---

## ANSI 生成

```go
// 补丁类型
type PatchType int

const (
    PatchResize PatchType = iota
    PatchCell
    PatchClear
)

// 补丁
type Patch struct {
    Type   PatchType
    X, Y   int
    Cell   Cell
    Width  int
    Height int
}

// ANSI 生成器
type ANSIGenerator struct {
    buf bytes.Buffer
    
    // 当前样式状态 (避免重复输出)
    currentStyle CellStyle
}

// 生成 ANSI 输出
func (g *ANSIGenerator) Generate(patches []Patch) []byte {
    g.buf.Reset()
    
    // 清屏 (如果需要)
    for _, patch := range patches {
        if patch.Type == PatchResize {
            g.writeClearScreen()
            break
        }
    }
    
    // 优化: 合并连续的同一行补丁
    optimized := g.optimizePatches(patches)
    
    // 生成每个补丁的 ANSI
    for _, patch := range optimized {
        switch patch.Type {
        case PatchCell:
            g.generateCell(patch)
        }
    }
    
    // 重置样式
    g.writeResetStyle()
    
    return g.buf.Bytes()
}

// 生成单元格 ANSI
func (g *ANSIGenerator) generateCell(patch Patch) {
    // 移动光标
    g.writeCursorMove(patch.X, patch.Y)
    
    // 设置样式 (只输出变化的部分)
    g.writeStyleDiff(patch.Cell.Style)
    
    // 输出字符
    g.buf.WriteRune(patch.Cell.Char)
}

// 写入样式差异
func (g *ANSIGenerator) writeStyleDiff(style CellStyle) {
    // 前景色
    if style.Foreground != g.currentStyle.Foreground {
        g.writeForegroundColor(style.Foreground)
    }
    
    // 背景色
    if style.Background != g.currentStyle.Background {
        g.writeBackgroundColor(style.Background)
    }
    
    // 粗体
    if style.FgBold != g.currentStyle.FgBold {
        if style.FgBold {
            g.buf.WriteString("\x1b[1m")
        } else {
            g.buf.WriteString("\x1b[22m")
        }
    }
    
    // 下划线
    if style.Underline != g.currentStyle.Underline {
        if style.Underline {
            g.buf.WriteString("\x1b[4m")
        } else {
            g.buf.WriteString("\x1b[24m")
        }
    }
    
    g.currentStyle = style
}

// 写入前景色
func (g *ANSIGenerator) writeForegroundColor(color Color) {
    switch color.Type {
    case ColorDefault:
        g.buf.WriteString("\x1b[39m")
    case ColorBasic:
        fmt.Fprintf(&g.buf, "\x1b[%dm", 30+color.Basic)
    case ColorExtended:
        fmt.Fprintf(&g.buf, "\x1b[38;5;%dm", color.Extended)
    case ColorTrueColor:
        fmt.Fprintf(&g.buf, "\x1b[38;2;%d;%d;%dm", color.R, color.G, color.B)
    }
}

// 写入光标移动
func (g *ANSIGenerator) writeCursorMove(x, y int) {
    // CSI H: \x1b[row;colH (1-indexed)
    fmt.Fprintf(&g.buf, "\x1b[%d;%dH", y+1, x+1)
}
```

---

## 补丁优化

```go
// 优化补丁: 合并连续单元格
func (g *ANSIGenerator) optimizePatches(patches []Patch) []Patch {
    if len(patches) <= 1 {
        return patches
    }
    
    // 按位置排序
    sort.Slice(patches, func(i, j int) bool {
        if patches[i].Y != patches[j].Y {
            return patches[i].Y < patches[j].Y
        }
        return patches[i].X < patches[j].X
    })
    
    optimized := []Patch{}
    
    for i := 0; i < len(patches); i++ {
        patch := patches[i]
        
        // 尝试合并连续单元格
        if len(optimized) > 0 {
            last := &optimized[len(optimized)-1]
            
            // 同一行且连续
            if last.Y == patch.Y && last.X+1 == patch.X &&
               last.Cell.Style == patch.Cell.Style {
                // 合并: 追加字符
                last.Cell.Char = last.Cell.Char + patch.Cell.Char
                continue
            }
        }
        
        optimized = append(optimized, patch)
    }
    
    return optimized
}
```

---

## 终端抽象

```go
// 终端接口
type Terminal interface {
    // 初始化
    Init() error
    Close() error
    
    // 尺寸
    Size() (width, height int)
    
    // 输出
    Write(data []byte) error
    
    // 事件
    Events() <-chan Event
    
    // 原始模式
    EnterRawMode() error
    ExitRawMode() error
}

// 事件
type Event interface {
    eventType()
}

type KeyEvent struct {
    Key  Key
    Rune rune
    Mod  Modifier
}

type ResizeEvent struct {
    Width  int
    Height int
}

type MouseEvent struct {
    X, Y int
    Button MouseButton
    Action MouseAction
}
```

---

## tcell 实现

```go
// tcell 终端实现
type TcellTerminal struct {
    screen tcell.Screen
    events chan Event
}

func (t *TcellTerminal) Init() error {
    screen, err := tcell.NewScreen()
    if err != nil {
        return err
    }
    
    if err := screen.Init(); err != nil {
        return err
    }
    
    t.screen = screen
    t.events = make(chan Event, 100)
    
    // 启动事件循环
    go t.eventLoop()
    
    return nil
}

func (t *TcellTerminal) Size() (int, int) {
    return t.screen.Size()
}

func (t *TcellTerminal) Write(data []byte) error {
    // tcell 直接写入屏幕
    // 这里我们使用 ANSI 序列
    return nil
}

func (t *TcellTerminal) eventLoop() {
    for {
        ev := t.screen.PollEvent()
        
        switch e := ev.(type) {
        case *tcell.EventKey:
            t.events <- KeyEvent{
                Key:  convertKey(e.Key()),
                Rune: e.Rune(),
            }
        case *tcell.EventResize:
            w, h := e.Size()
            t.events <- ResizeEvent{Width: w, Height: h}
        }
    }
}
```

---

## 渲染循环

```go
// 渲染循环
func (r *Renderer) Run(root *LayoutNode) {
    // 初始渲染
    buffer := r.BuildBuffer(root)
    patches := r.Diff(nil, buffer)
    ansi := r.ansiGen.Generate(patches)
    r.terminal.Write(ansi)
    r.previous = buffer
    
    // 事件循环
    for {
        select {
        case ev := <-r.terminal.Events():
            switch e := ev.(type) {
            case ResizeEvent:
                // 终端尺寸变化
                r.handleResize(e)
            }
            
        case <-r.renderTick:
            // 定时渲染
            r.render(root)
        }
    }
}

// 渲染一帧
func (r *Renderer) render(root *LayoutNode) {
    // 1. 构建缓冲
    buffer := r.BuildBuffer(root)
    
    // 2. 计算差异
    patches := r.Diff(r.previous, buffer)
    
    if len(patches) == 0 {
        return  // 无变化
    }
    
    // 3. 生成 ANSI
    ansi := r.ansiGen.Generate(patches)
    
    // 4. 输出
    r.terminal.Write(ansi)
    
    // 5. 保存当前缓冲
    r.previous = buffer
}
```

---

## 模块结构

```
renderer/
├── renderer.go        # 渲染器主入口
├── buffer.go          # 输出缓冲
├── cell.go            # 单元格定义
├── build.go           # 缓冲构建
├── diff.go            # 差异计算
├── ansi.go            # ANSI 生成
├── optimize.go        # 补丁优化
├── terminal.go        # 终端接口
├── tcell.go           # tcell 实现
├── events.go          # 事件定义
└── colors.go          # 颜色定义
```

---

## 总结

### 核心要点

1. **增量渲染**
   - 双缓冲对比
   - 只输出变化部分
   - 最小化 ANSI 序列

2. **差异算法**
   - 逐行比较
   - 逐单元格比较
   - 跳过未变化区域

3. **ANSI 优化**
   - 样式状态缓存
   - 连续单元格合并
   - 避免冗余序列

4. **终端抽象**
   - 统一接口
   - tcell 跨平台实现
   - 事件驱动

5. **性能优化**
   - 脏标记
   - 补丁排序
   - 批量输出
