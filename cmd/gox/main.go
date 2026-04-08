// Package main 提供 gox 编译命令行工具
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayanmw/go-react-ink/internal/compiler"
)

var (
	// 版本信息
	version = "dev"

	// 命令行参数
	flagOut     = flag.String("o", "", "Output file (default: stdout)")
	flagPkg     = flag.String("pkg", "ink", "Component package name")
	flagWatch   = flag.Bool("watch", false, "Watch for file changes")
	flagVersion = flag.Bool("version", false, "Show version")
	flagHelp    = flag.Bool("help", false, "Show help")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *flagVersion {
		fmt.Printf("gox %s\n", version)
		os.Exit(0)
	}

	if *flagHelp {
		usage()
		os.Exit(0)
	}

	// 过滤掉空参数（可能是 flag 已处理的）
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input files")
		usage()
		os.Exit(1)
	}

	// 检查是否有未解析的 flag（用户可能把 flag 放在了输入文件后面）
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && arg != "-" {
			fmt.Fprintf(os.Stderr, "Error: flags must come before input files\n")
			fmt.Fprintf(os.Stderr, "  Example: gox -o output.go input.gox\n")
			os.Exit(1)
		}
	}

	// 创建编译器
	c := compiler.New(*flagPkg)

	// 处理输入
	if len(args) == 1 && !isDir(args[0]) {
		// 单文件编译
		compileFile(c, args[0])
	} else {
		// 多文件或目录编译
		compileInputs(c, args)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `gox - JSX to Go compiler for go-ink

Usage:
  gox [options] <input.gox...>
  gox [options] <input-dir>

Options:`)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, `
Examples:
  gox app.gox                    # Compile single file to stdout
  gox -o app.go app.gox          # Compile to file
  gox ./src                      # Compile all .gox files in directory
  gox -watch ./src               # Watch and recompile on changes
  gox -pkg ink app.gox           # Use 'ink' as component package`)
}

func compileFile(c *compiler.Compiler, input string) {
	output, err := c.CompileFile(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *flagOut != "" {
		if err := os.WriteFile(*flagOut, output, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error: write %s: %v\n", *flagOut, err)
			os.Exit(1)
		}
	} else {
		os.Stdout.Write(output)
	}
}

func compileInputs(c *compiler.Compiler, inputs []string) {
	for _, input := range inputs {
		if isDir(input) {
			// 编译目录
			if err := compileDir(c, input); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else {
			// 编译文件
			output, err := c.CompileFile(input)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s: %v\n", input, err)
				continue
			}

			// 确定输出路径
			outPath := *flagOut
			if outPath == "" {
				// 默认: .gox -> .go
				outPath = input[:len(input)-1] + "o"
			}

			if err := os.WriteFile(outPath, output, 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error: write %s: %v\n", outPath, err)
				continue
			}

			fmt.Printf("%s -> %s\n", input, outPath)
		}
	}
}

func compileDir(c *compiler.Compiler, dir string) error {
	// 遍历目录
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// 只处理 .gox 文件
		if !strings.HasSuffix(path, ".gox") {
			return nil
		}

		// 编译文件
		output, err := c.CompileFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s: %v\n", path, err)
			return nil
		}

		// 输出路径: .gox -> .go
		outPath := path[:len(path)-1] + "o"

		if err := os.WriteFile(outPath, output, 0644); err != nil {
			return fmt.Errorf("write %s: %w", outPath, err)
		}

		fmt.Printf("%s -> %s\n", path, outPath)
		return nil
	})
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
