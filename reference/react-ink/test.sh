#!/usr/bin/env bash
# React Ink 测试脚本
# 作为 Makefile 的替代，在 Windows 上可直接运行

set -e

# 颜色输出
if [ -t 1 ] && [ -z "$NO_COLOR" ]; then
    COLOR_RESET='\033[0m'
    COLOR_BOLD='\033[1m'
    COLOR_GREEN='\033[32m'
    COLOR_BLUE='\033[34m'
    COLOR_YELLOW='\033[33m'
    COLOR_CYAN='\033[36m'
    COLOR_RED='\033[31m'
else
    COLOR_RESET=''
    COLOR_BOLD=''
    COLOR_GREEN=''
    COLOR_BLUE=''
    COLOR_YELLOW=''
    COLOR_CYAN=''
    COLOR_RED=''
fi

# 项目目录
EXAMPLES_DIR="./examples"
TEST_DIR="./test"
TSX="npx tsx"

# 帮助信息
show_help() {
    echo "${COLOR_BOLD}React Ink 测试脚本${COLOR_RESET}"
    echo ""
    echo "用法: $0 <command> [args]"
    echo ""
    echo "${COLOR_BOLD}可用命令:${COLOR_RESET}"
    echo "  help                  显示此帮助信息"
    echo "  build                 编译 TypeScript"
    echo "  clean                 清理构建输出"
    echo "  typecheck             类型检查"
    echo "  lint                  代码检查"
    echo "  format                格式化代码"
    echo "  check                 运行所有检查"
    echo "  test                  运行测试套件"
    echo ""
    echo "  demo                  运行计数器示例（最简单）"
    echo "  run <name>            运行指定示例（如: run borders）"
    echo "  list                  列出所有示例"
    echo ""
    echo "  run-all               依次运行所有示例"
    echo "  info                  显示项目信息"
    echo "  watch                 监听模式编译"
    echo ""
    echo "${COLOR_YELLOW}注意: 带 * 的命令会启动交互式程序，需要 Ctrl+C 退出${COLOR_RESET}"
    exit 0
}

# 编译脚本
cmd_build() {
    echo "${COLOR_BLUE}编译项目...${COLOR_RESET}"
    npm run build
    echo "${COLOR_GREEN}✓ 编译完成${COLOR_RESET}"
}

# 清理构建输出
cmd_clean() {
    echo "${COLOR_BLUE}清理构建输出...${COLOR_RESET}"
    rm -rf build
    echo "${COLOR_GREEN}✓ 清理完成${COLOR_RESET}"
}

# 类型检查
cmd_typecheck() {
    echo "${COLOR_BLUE}运行类型检查...${COLOR_RESET}"
    npm run typecheck
    echo "${COLOR_GREEN}✓ 类型检查通过${COLOR_RESET}"
}

# Lint
cmd_lint() {
    echo "${COLOR_BLUE}运行代码检查...${COLOR_RESET}"
    npm run lint
    echo "${COLOR_GREEN}✓ Lint 检查完成${COLOR_RESET}"
}

# 格式化
cmd_format() {
    echo "${COLOR_BLUE}格式化代码...${COLOR_RESET}"
    npm run format
    echo "${COLOR_GREEN}✓ 格式化完成${COLOR_RESET}"
}

# 运行测试套件
cmd_test() {
    echo "${COLOR_BLUE}运行测试套件...${COLOR_RESET}"
    npm test
}

# 运行所有检查
cmd_check() {
    echo "${COLOR_BLUE}运行所有检查...${COLOR_RESET}"
    cmd_typecheck
    cmd_lint
    echo "${COLOR_GREEN}✓ 所有检查通过${COLOR_RESET}"
}

# 运行示例
cmd_run() {
    local name="$1"
    if [ -z "$name" ]; then
        echo "${COLOR_RED}错误: 请指定示例名称${COLOR_RESET}"
        cmd_list
        exit 1
    fi

    local example_file="${EXAMPLES_DIR}/${name}/${name}.tsx"
    if [ ! -f "$example_file" ]; then
        echo "${COLOR_RED}错误: 示例 '$name' 不存在${COLOR_RESET}"
        cmd_list
        exit 1
    fi

    echo "${COLOR_YELLOW}运行示例: $name${COLOR_RESET}"
    echo "${COLOR_BLUE}按 Ctrl+C 退出${COLOR_RESET}"
    $TSX "$example_file"
}

# 运行 demo
cmd_demo() {
    cmd_run counter
}

# 列出所有示例
cmd_list() {
    echo "${COLOR_BOLD}可用示例:$(COLOR_RESET)"
    for dir in "$EXAMPLES_DIR"/*/; do
        name=$(basename "$dir")
        example_file="${dir}${name}.tsx"
        if [ -f "$example_file" ]; then
            echo "  - $name"
        fi
    done
}

# 依次运行所有示例
cmd_run_all() {
    echo "${COLOR_YELLOW}注意: 需要为每个示例按 Ctrl+C 退出${COLOR_RESET}"
    for dir in "$EXAMPLES_DIR"/*/; do
        name=$(basename "$dir")
        example_file="${dir}${name}.tsx"
        if [ -f "$example_file" ]; then
            echo ""
            echo "${COLOR_BOLD}=== 运行: $name ===${COLOR_RESET}"
            $TSX "$example_file" || true
            sleep 1
        fi
    done
    echo ""
    echo "${COLOR_GREEN}✓ 所有示例运行完成${COLOR_RESET}"
}

# 显示项目信息
cmd_info() {
    echo "${COLOR_BOLD}项目信息:$(COLOR_RESET)"
    echo "  Node 版本:     $(node --version)"
    echo "  npm 版本:      $(npm --version)"
    echo "  tsx 版本:      $($TSX --version 2>&1 | head -1)"
    echo ""
    echo "${COLOR_BOLD}文件统计:$(COLOR_RESET)"
    echo "  源文件数:     $(find src -name "*.ts" -o -name "*.tsx" 2>/dev/null | wc -l)"
    echo "  示例数:       $(find examples -name "*.tsx" 2>/dev/null | wc -l)"
    echo "  测试文件数:   $(find test -name "*.ts" -o -name "*.tsx" 2>/dev/null | wc -l)"
}

# 监听模式
cmd_watch() {
    echo "${COLOR_BLUE}启动监听模式...${COLOR_RESET}"
    npm run dev
}

# Dev 模式
cmd_dev() {
    echo "${COLOR_BOLD}进入开发工具模式（每5秒执行一次:typecheck lint format）$(COLOR_RESET)"
    while true; do
        cmd_typecheck
        cmd_lint
        cmd_format
        echo "${COLOR_GREEN}完成！等待5秒...${COLOR_RESET}"
        sleep 5
    done
}

# 主函数
main() {
    case "${1:-}" in
        help|"-h"|"--help"|"")
            show_help
            ;;
        build)
            cmd_build
            ;;
        clean)
            cmd_clean
            ;;
        typecheck)
            cmd_typecheck
            ;;
        lint)
            cmd_lint
            ;;
        format)
            cmd_format
            ;;
        check)
            cmd_check
            ;;
        test)
            cmd_test
            ;;
        run)
            cmd_run "$2"
            ;;
        run-all)
            cmd_run_all
            ;;
        demo)
            cmd_demo
            ;;
        list)
            cmd_list
            ;;
        info)
            cmd_info
            ;;
        watch)
            cmd_watch
            ;;
        dev)
            cmd_dev
            ;;
        *)
            echo "${COLOR_RED}错误: 未知命令 '$1'${COLOR_RESET}"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

main "$@"
