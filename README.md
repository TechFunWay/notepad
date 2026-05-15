# 记事本 (Notepad)

一款轻量、精美的多用户记事本 Web 应用，支持富文本编辑、标签分类管理、暗色模式、多用户数据隔离，可部署到飞牛NAS、Docker 或直接运行。

## 技术栈

- **前端**: Vue 3 + Element Plus + Vite
- **后端**: Go + Gin + SQLite (modernc.org/sqlite, 纯 Go 无 CGO)
- **认证**: JWT + bcrypt

## 功能特性

- 笔记的创建、编辑、删除、搜索
- 多用户注册登录，数据完全隔离
- 管理员管理用户和系统配置
- 安全问题找回密码
- 终端命令恢复管理员账号
- 一条数据一个配置项，管理员可控制是否允许注册等
- 首个注册用户自动成为超级管理员（唯一）
- 数据库自动迁移，支持版本升级

## 快速开始

### 直接运行

```bash
# 从 release 目录下载对应平台的二进制文件
./notepad

# 访问 http://localhost:8904
# 第一个注册的用户自动成为管理员
```

### Docker 运行

```bash
docker run -d \
  --name notepad \
  -p 8904:8904 \
  -v ./data:/app/data \
  -e JWT_SECRET=your-secret-key \
  wycto/notepad:latest
```

### Docker Compose

```bash
# 编辑 docker-compose.yaml 中的 JWT_SECRET
docker compose up -d
```

### 飞牛NAS 安装

1. 下载 `notepad_<version>_fpk.tar.gz`
2. 在飞牛应用中心选择"手动安装"
3. 上传 FPK 包
4. 按向导提示配置端口和 JWT 密钥

## 终端命令

```bash
# 查看管理员用户名
./notepad find-admin

# 重置管理员密码
./notepad recover-admin

# 列出所有用户
./notepad list-users
```

## 环境变量

| 变量 | 默认值 | 说明 |
|------|--------|------|
| PORT | 8904 | 服务端口 |
| DB_PATH | ./data/notepad.db | SQLite 数据库路径 |
| JWT_SECRET | 随机生成 | JWT 签名密钥（建议设置固定值） |
| DATA_DIR | ./data | 数据目录 |

## 从源码构建

```bash
# 前置要求: Go 1.22+, Node.js 20+

# 完整构建（前端+后端+FPK）
make build-fpk

# 或使用 Makefile
make build          # 当前平台
make cross-compile  # 多平台交叉编译
make build-fpk      # 飞牛 FPK 包
```

构建产物输出到 `release/<version>/` 目录。

## 项目结构

```
├── server/              # Go 后端
│   ├── main.go          # 入口（服务器/CLI）
│   ├── cmd/             # 服务启动 + CLI 命令
│   ├── config/          # 配置加载
│   ├── database/        # SQLite + 迁移
│   ├── model/           # 数据模型
│   ├── handler/         # API 处理器
│   ├── middleware/       # 认证/CORS
│   ├── auth/            # JWT 工具
│   ├── router/          # 路由注册
│   └── static/          # 嵌入前端资源
├── web/                 # Vue 3 前端
│   └── src/
│       ├── views/       # 页面视图
│       ├── components/  # 布局组件
│       ├── stores/      # Pinia 状态
│       ├── api/         # API 调用
│       └── router/      # 前端路由
├── scripts/             # 构建脚本
│   ├── build-all.sh     # 多平台编译
│   ├── build-fnpack.sh  # 飞牛 FPK 打包
│   └── build-docker.sh  # Docker 镜像
├── fnpack/              # 飞牛 FPK 模板
├── VERSION              # 版本号
├── Makefile             # 构建入口
└── Dockerfile           # Docker 多阶段构建
```

## 升级

数据库支持自动迁移。升级时只需替换二进制文件或更新 Docker 镜像，数据库会自动升级到最新版本。

## License

MIT
