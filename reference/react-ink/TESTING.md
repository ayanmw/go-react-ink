# React Ink 测试脚本使用指南

本目录包含两个测试工具：

## 1. Makefile（推荐用于 Unix-like 系统）
如果系统安装了 `make`，可以使用以下命令：

```bash
make help           # 显示所有可用命令
make build         # 编译 TypeScript
make typecheck     # 类型检查
make lint          # 代码检查
make test          # 运行测试套件
make check         # 运行所有检查

# 运行示例
make demo          # 运行计数器示例（最简单）
make example-bords  # 运行边框示例
make example-use-animation  # 运行动画示例

# 列出所有示例
make list-examples

# 运行指定示例
make run-example-by-name NAME=borders
```

## 2. test.sh（跨平台，推荐）
适用于 Windows、Linux 和 macOS 的 Bash 测试脚本：

```bash
bash test.sh help    # 显示所有可用命令
bash test.sh build   # 编译 TypeScript
bash test.sh typecheck   # 类型检查
bash test.sh lint    # 代码检查
bash test.sh test    # 运行测试套件
bash test.sh check   # 运行所有检查

# 运行示例
bash test.sh demo    # 运行计数器示例（最简单）
bash test.sh run borders    # 运行边框示例
bash test.sh run use-animation  # 运行动画示例

# 列出所有示例
bash test.sh list

# 运行所有示例
bash test.sh run-all

# 显示项目信息
bash test.sh info

# 监听模式
bash test.sh watch

# 开发模式（循环执行检查）
bash test.sh dev
```

注意：交互式示例（带动画、输入功能的）需要按 **Ctrl+C** 退出。
