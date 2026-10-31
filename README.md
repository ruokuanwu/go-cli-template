# go-cli-template

基于 Go 标准库的轻量 CLI 模板，无第三方运行依赖。提供命令分发、帮助、版本注入、示例子命令、单元测试、VS Code 调试与容器运行配置。

## 环境要求

- Go 1.25 或更高版本，`go` 已加入 `PATH`。
- Make（可选）；执行 `make test` 的竞态检测需要 CGO 和 C 编译器。
- VS Code Go 扩展及 Delve（仅调试需要）。
- Docker / Docker Compose v2（仅容器运行需要）。

## 快速开始

```bash
go run ./cmd/go-cli-template --help
go run ./cmd/go-cli-template greet --name Go
make build VERSION=v0.1.0
./bin/go-cli-template --version
```

不使用 Make 时：

```bash
go test ./...
go build -o bin/go-cli-template ./cmd/go-cli-template
```

CLI 不传参数时显示帮助；`greet` 默认输出 `Hello, world!`。`greet --help` 显示子命令帮助。成功退出码为 `0`，命令或参数错误为 `2`，错误写入标准错误流。子命令选项需放在位置参数之前。

## 常用命令

| 命令 | 用途 |
| --- | --- |
| `make help` | 列出可用目标 |
| `make run ARGS="greet --name Go"` | 运行示例命令 |
| `make test` | 竞态检测及覆盖率测试 |
| `make vet` | 静态检查 |
| `make fmt` | 格式化代码 |
| `make check` | 静态检查、测试及格式检查；格式不符时修正并失败 |
| `make build VERSION=v0.1.0` | 构建并注入版本号，默认版本为 `dev` |
| `make clean` | 清理构建产物 |
| `make docker-build IMAGE=go-cli-template:local` | 构建非 root 容器镜像 |

VS Code 打开项目后，使用 `Debug CLI` 配置启动调试。

## 容器运行

CLI 是一次性任务，不启动常驻服务。

```bash
make docker-build
docker run --rm go-cli-template:local greet --name Go

# 测试环境：本地构建并执行
docker compose -f deploy/test/compose.yaml run --build --rm cli

# 生产环境：使用已发布且固定版本的镜像
export CLI_IMAGE=your-registry/go-cli-template:v0.1.0
docker compose -f deploy/prod/compose.yaml run --rm cli
```

生产配置仅引用镜像，不负责发布；请先自行构建并推送镜像。

## 目录结构

```text
.
├── cmd/go-cli-template/main.go  # 入口及版本变量
├── internal/cli/               # 命令实现与测试
├── .vscode/launch.json         # Go 调试配置
├── docker/Dockerfile           # 多阶段镜像构建
├── deploy/test/compose.yaml    # 测试环境一次性任务
├── deploy/prod/compose.yaml    # 生产环境一次性任务
├── go.mod
├── Makefile
└── README.md
```

## 使用为新项目

1. 将 `go.mod` 中的 `example.com/go-cli-template` 改为实际模块路径，同步修改入口文件的导入路径。
2. 按需重命名 `cmd/go-cli-template`，同步修改 CLI 帮助、测试、Makefile、调试配置、Dockerfile 和部署配置中的名称或路径。
3. 在 `internal/cli` 中添加子命令及对应测试，保持入口仅负责进程退出和标准流接入。
