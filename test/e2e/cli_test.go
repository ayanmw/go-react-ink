package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestCLI 测试命令行工具
func TestCLI(t *testing.T) {
	// 构建 CLI 工具
	binPath := filepath.Join(t.TempDir(), "gox")
	if filepath.Ext(binPath) == "" {
		binPath += ".exe" // Windows
	}
	buildCmd := exec.Command("go", "build", "-o", binPath, "github.com/ayanmw/go-react-ink/cmd/gox")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\n%s", err, output)
	}

	// 创建测试文件
	testDir := t.TempDir()
	inputPath := filepath.Join(testDir, "test.gox")
	source := `package main

func App() Element {
	return <Box><Text>Hello CLI</Text></Box>
}
`
	if err := os.WriteFile(inputPath, []byte(source), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// 运行编译
	outputPath := filepath.Join(testDir, "test.go")
	cmd := exec.Command(binPath, "-o", outputPath, inputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("CLI failed: %v\n%s", err, output)
	}

	// 检查输出文件
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)
	expectedStrings := []string{
		"package main",
		"func App()",
		"CreateElement",
		"Hello CLI",
	}

	for _, expected := range expectedStrings {
		if !contains(contentStr, expected) {
			t.Errorf("Output missing expected string %q", expected)
		}
	}
}

// TestCLIHelp 测试帮助命令
func TestCLIHelp(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "gox")
	if filepath.Ext(binPath) == "" {
		binPath += ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", binPath, "github.com/ayanmw/go-react-ink/cmd/gox")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\n%s", err, output)
	}

	cmd := exec.Command(binPath, "-help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Help command failed: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "gox") || !contains(outputStr, "JSX") {
		t.Errorf("Help output should contain usage info, got: %s", outputStr)
	}
}

// TestCLIVersion 测试版本命令
func TestCLIVersion(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "gox")
	if filepath.Ext(binPath) == "" {
		binPath += ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", binPath, "github.com/ayanmw/go-react-ink/cmd/gox")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\n%s", err, output)
	}

	cmd := exec.Command(binPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "gox") {
		t.Errorf("Version output should contain 'gox', got: %s", outputStr)
	}
}

// TestDirectoryCompilation 测试目录编译
func TestDirectoryCompilation(t *testing.T) {
	binPath := filepath.Join(t.TempDir(), "gox")
	if filepath.Ext(binPath) == "" {
		binPath += ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", binPath, "github.com/ayanmw/go-react-ink/cmd/gox")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build CLI: %v\n%s", err, output)
	}

	// 创建测试目录
	testDir := t.TempDir()

	// 创建多个文件
	files := map[string]string{
		"a.gox": `package main

func A() Element {
	return <Text>A</Text>
}
`,
		"b.gox": `package main

func B() Element {
	return <Text>B</Text>
}
`,
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(testDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", name, err)
		}
	}

	// 运行目录编译
	cmd := exec.Command(binPath, testDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Directory compilation failed: %v\n%s", err, output)
	}

	// 检查生成的文件
	for name := range files {
		goPath := filepath.Join(testDir, name[:len(name)-1]+"o")
		if _, err := os.Stat(goPath); os.IsNotExist(err) {
			t.Errorf("Generated file not found: %s", goPath)
		}
	}
}

// 辅助函数
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func Example() {
	fmt.Println("See TestCLI for CLI usage example")
	// Output: See TestCLI for CLI usage example
}
