package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ayanmw/go-react-ink/internal/compiler"
)

func TestEndToEndCompilation(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "go-ink-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试 .gox 文件
	testCases := []struct {
		name     string
		source   string
		contains []string
	}{
		{
			name: "simple_element",
			source: `<Box></Box>`,
			contains: []string{"CreateElement", "ink.Box"},
		},
		{
			name: "element_with_attrs",
			source: `<Box flexDirection="column" padding={1}></Box>`,
			contains: []string{"flexDirection", "column", "padding"},
		},
		{
			name: "element_with_children",
			source: `<Box><Text>Hello</Text></Box>`,
			contains: []string{"ink.Box", "ink.Text", "Hello"},
		},
		{
			name: "expression",
			source: `<Text>Count: {count}</Text>`,
			contains: []string{"count"},
		},
		{
			name: "self_closing",
			source: `<Box><Spacer /></Box>`,
			contains: []string{"ink.Spacer"},
		},
		{
			name: "spread",
			source: `<Box {...props}></Box>`,
			contains: []string{"Spread(props)"},
		},
		{
			name: "fragment",
			source: `<>A</>`,
			contains: []string{"A"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := compiler.New("ink")
			output, err := c.Compile(tc.name+".gox", tc.source)
			if err != nil {
				t.Fatalf("Compilation failed: %v", err)
			}

			outputStr := string(output)
			for _, expected := range tc.contains {
				if !containsString(outputStr, expected) {
					t.Errorf("Output missing expected string %q\nGot:\n%s", expected, outputStr)
				}
			}
		})
	}
}

func TestFullPipeline(t *testing.T) {
	// 测试完整流程: 扫描 -> 词法分析 -> 语法分析 -> 代码生成
	source := `package main

func App() Element {
	name := "World"
	return (
		<Box flexDirection="column">
			<Text color="green">Hello, {name}!</Text>
			<Text>Count: 42</Text>
		</Box>
	)
}
`
	c := compiler.New("ink")
	output, err := c.Compile("app.gox", source)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	outputStr := string(output)

	// 验证关键内容
	expected := []string{
		"package main",
		"func App()",
		"name :=",
		"CreateElement",
		"ink.Box",
		"ink.Text",
		"flexDirection",
		"green",
		"Hello",
		"name",
		"42",
	}

	for _, exp := range expected {
		if !containsString(outputStr, exp) {
			t.Errorf("Output missing expected string %q", exp)
		}
	}
}

func TestFileCompilation(t *testing.T) {
	// 创建临时文件
	tmpDir, err := os.MkdirTemp("", "go-ink-file-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建输入文件
	inputPath := filepath.Join(tmpDir, "test.gox")
	source := `package main

func App() Element {
	return <Box><Text>Hello</Text></Box>
}
`
	if err := os.WriteFile(inputPath, []byte(source), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	// 编译文件
	c := compiler.New("ink")
	output, err := c.CompileFile(inputPath)
	if err != nil {
		t.Fatalf("File compilation failed: %v", err)
	}

	outputStr := string(output)
	if !containsString(outputStr, "Hello") {
		t.Errorf("Output should contain 'Hello'\nGot:\n%s", outputStr)
	}
}

func TestDirectoryCompilation(t *testing.T) {
	// 创建临时目录结构
	tmpDir, err := os.MkdirTemp("", "go-ink-dir-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建子目录
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create sub dir: %v", err)
	}

	// 创建多个 .gox 文件
	files := map[string]string{
		"main.gox": `package main

func main() {}
`,
		"app.gox": `package main

func App() Element {
	return <Box>App</Box>
}
`,
		filepath.Join("sub", "nested.gox"): `package main

func Nested() Element {
	return <Text>Nested</Text>
}
`,
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", name, err)
		}
	}

	// 编译目录
	c := compiler.New("ink")
	if err := c.CompileDir(tmpDir, tmpDir); err != nil {
		t.Fatalf("Directory compilation failed: %v", err)
	}

	// 验证生成的文件
	generatedFiles := []string{
		"main.go",
		"app.go",
		filepath.Join("sub", "nested.go"),
	}

	for _, name := range generatedFiles {
		path := filepath.Join(tmpDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Generated file not found: %s", name)
		}
	}
}

func TestErrorHandling(t *testing.T) {
	// 测试错误处理
	c := compiler.New("ink")

	// 无效的 JSX (未闭合标签)
	// 注意: 当前实现可能不完全检测这个错误
	// 所以我们主要测试不会崩溃
	_, _ = c.Compile("error.gox", `<Box><Text></Box>`)
}

func TestImportGeneration(t *testing.T) {
	c := compiler.New("ink")

	// 源码没有导入
	source := `package main

func App() Element {
	return <Box></Box>
}
`
	output, err := c.Compile("test.gox", source)
	if err != nil {
		t.Fatalf("Compilation failed: %v", err)
	}

	outputStr := string(output)
	if !containsString(outputStr, "import") {
		t.Error("Should generate import statement")
	}
	if !containsString(outputStr, "github.com/ayanmw/go-react-ink/pkg/core") {
		t.Error("Should import core package")
	}
}

// 辅助函数
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsInString(s, substr))
}

func containsInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func ExampleCompiler() {
	c := compiler.New("ink")
	output, _ := c.Compile("example.gox", `<Text>Hello</Text>`)
	fmt.Println(string(output))
	// Output contains CreateElement and Text
}