# go-cli-template

沿用当前目录结构的 Go CLI 模板，使用 Cobra（命令与参数）、Zap（结构化日志）、Oops（错误包装）和 GORM（数据库接入）。提供帮助、版本注入、示例子命令、单元测试、VS Code 调试与容器运行配置。

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

CLI 不传参数时显示帮助；`greet` 默认输出 `Hello, world!`。`greet --help` 显示子命令帮助。成功退出码为 `0`，命令或参数错误为 `2`，应用通过 Zap 将错误以 JSON 写入标准错误流。

## 配置与分层

`internal/config` 定义并加载日志与数据库配置：

- `APP_LOG_LEVEL`：默认 `info`，支持 `debug`、`info`、`warn`、`error`、`dpanic`、`panic`、`fatal`。无效级别会在启动时返回错误。
- `APP_DB_DSN`：SQLite DSN，默认 `:memory:`，每次运行使用独立内存数据库。设置为文件路径可跨进程持久化；父目录须存在且当前用户可写。

```bash
APP_LOG_LEVEL=debug go run ./cmd/go-cli-template greet --name Go
```

`internal/app` 负责初始化配置、全局日志及依赖，并定义 Cobra 命令，不再保留独立 CLI 目录。`internal/pkg/log.Init` 根据配置初始化 Zap 全局对象，各层直接使用 `zap.L()` 获取日志；初始化只需在启动时执行。日志输出到标准错误流，正常命令结果输出到标准输出流。

示例调用链：`app → service.Greeting → repository.Greeting → GORM → SQLite`，数据结构为 `model.Greeting`。Service 校验姓名并拼接问候语；Repository 使用注入的数据库查询问候前缀，不再返回硬编码数据。

`internal/db.Open` 接收 `gorm.Dialector` 和可选 GORM 配置，返回数据库、清理函数和错误，不绑定数据库类型，也不限制连接池。GORM 不内置默认数据库驱动；App 当前显式选择纯 Go SQLite 驱动 `github.com/glebarez/sqlite`，兼容 `CGO_ENABLED=0` 容器构建。更换数据库时替换 App 中的驱动即可。

App 打开数据库后使用 `defer cleanup()` 释放连接，无需接触底层 `sql.DB`。Repository 初始化表并插入默认 `Hello` 记录，重复初始化不会覆盖已有数据。数据库错误通过 Oops 包装并保留错误链。当前 CLI 顺序执行查询，默认内存数据库无需额外连接池设置；扩展为并发访问时，应按所选数据库配置连接池及 DSN。

```bash
APP_DB_DSN=./app.db go run ./cmd/go-cli-template greet --name Go
```

容器默认使用内存数据库，无需写入文件。需要持久化时，将 `APP_DB_DSN` 指向挂载目录中的文件，并确保 UID 65532 有写入权限。

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
├── internal/
│   ├── app/                   # Cobra 命令、应用装配与测试
│   ├── config/                # 配置定义与环境变量加载
│   ├── db/                    # GORM 初始化与错误包装
│   ├── model/                 # 问候数据模型
│   ├── pkg/
│   │   ├── client/            # 外部客户端（预留）
│   │   └── log/               # 根据配置初始化 Zap 全局日志
│   ├── repository/            # GORM 表初始化与问候数据查询
│   └── service/               # 姓名校验与问候逻辑
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
3. 在 `internal/app` 中添加子命令及对应测试，业务逻辑放在 `internal/service`，数据访问放在 `internal/repository`，数据结构放在 `internal/model`；保持入口仅负责进程退出和标准流接入。
