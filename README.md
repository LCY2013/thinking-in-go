# Thinking in Go | Go语言百科全书

[![Go Version](https://img.shields.io/badge/Go-1.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](https://github.com/LCY2013/thinking-in-go)

## 项目简介 | Project Description

这是一个雄心勃勃的开源书籍项目，旨在创建最全面的Go编程参考，涵盖所有基础概念、高级模式和云原生实现。这本百科全书既可作为新手的入门路径，也可作为经验丰富的Gopher在现代基础设施中工作的权威参考。

This ambitious open-source book project aims to create the most comprehensive Go programming reference covering all fundamental concepts, advanced patterns, and cloud-native implementations. The encyclopedia will serve as both a learning path for newcomers and a definitive reference for experienced Gophers working in modern infrastructure.

## 核心特性 | Core Features

### 🚀 Go语言基础 | Go Language Fundamentals
- **深入语法覆盖** | In-depth coverage of Go syntax, types, and core constructs
- **内存管理机制** | Memory management and runtime mechanics  
- **并发原语** | Concurrency primitives (goroutines, channels, sync patterns)
- **接口机制** | Interface mechanics and type system
- **性能优化** | Performance optimization and profiling

### ☁️ 云原生生态 | Cloud Native Ecosystem
- **Kubernetes Operator开发** | Kubernetes operator development with client-go
- **服务网格实现** | Service mesh implementations (Istio, Linkerd integration)
- **云提供商SDK** | Cloud provider SDKs and infrastructure-as-code patterns
- **分布式系统原语** | Distributed systems primitives (raft, gossip protocols)
- **可观测性栈** | Observability stack (OpenTelemetry, Prometheus exporters)

### 🔧 中间件实现 | Middleware Implementations
- **Redis协议实现** | Redis protocol implementation (RESP parser, command handling)
- **高性能HTTP/2和gRPC服务器** | High-performance HTTP/2 and gRPC servers
- **自定义数据库引擎** | Custom database engines (LSM tree, B+ tree implementations)
- **消息队列系统** | Message queue systems (AMQP, MQTT brokers)
- **缓存层** | Caching layers with expiration strategies

### 🎯 Go设计模式 | Design Patterns in Go
- **惯用Go模式** | Idiomatic Go patterns (functional options, decorators)
- **并发模式** | Concurrency patterns (worker pools, fan-out/fan-in)
- **架构模式** | Architectural patterns (CQRS, event sourcing)
- **微服务模式** | Microservices patterns (circuit breakers, bulkheads)
- **反模式和常见陷阱** | Anti-patterns and common pitfalls

### 🏭 生产级实践 | Production-Grade Practices
- **安全编码指南** | Secure coding guidelines
- **性能基准测试套件** | Performance benchmarking suites
- **跨平台编译技术** | Cross-platform compilation techniques
- **插件系统和动态加载** | Plugin systems and dynamic loading
- **WASM编译目标** | WASM compilation targets

## 项目结构 | Project Structure

```
thinking-in-go/
├── action/              # 基础教程和实战案例 | Basic tutorials and practical cases
│   ├── 01_introduce/    # Go语言介绍 | Go language introduction
│   ├── 02_quickstart/   # 快速开始 | Quick start guide
│   ├── 03_package_tools/# 包管理工具 | Package management tools
│   ├── 04_array_slice_map/ # 数组、切片、映射 | Arrays, slices, maps
│   ├── 05_go_types_system/ # Go类型系统 | Go type system
│   ├── 06_concurrent/   # 并发编程 | Concurrent programming
│   ├── 07_concurrent_mode/ # 并发模式 | Concurrency patterns
│   ├── 08_standard_doc/ # 标准库文档 | Standard library documentation
│   ├── 09_test/         # 测试 | Testing
│   ├── 10_struct/       # 结构体 | Structs
│   ├── 11_func/         # 函数 | Functions
│   ├── 12_method/       # 方法 | Methods
│   ├── 13_defer/        # defer语句 | Defer statements
│   ├── 14_interface/    # 接口 | Interfaces
│   ├── 15_concurrent/   # 并发进阶 | Advanced concurrency
│   ├── 16_action/       # 实战项目 | Practical projects
│   ├── 17_timer/        # 定时器 | Timers
│   ├── 18_pprof/        # 性能分析 | Performance profiling
│   └── 19_dynamic_param_lua/ # 动态参数和Lua | Dynamic parameters and Lua
├── advanced/            # 高级主题 | Advanced topics
│   ├── async/           # 异步编程 | Asynchronous programming
│   ├── concurrent/      # 高级并发 | Advanced concurrency
│   ├── errors/          # 错误处理 | Error handling
│   ├── reordering/      # 内存重排序 | Memory reordering
│   └── tools/           # 开发工具 | Development tools
├── ai/                  # AI集成 | AI integration
│   ├── ollama-chat/     # Ollama聊天 | Ollama chat
│   ├── ollama-completion/ # Ollama完成 | Ollama completion
│   └── nomic-embed-text/ # 文本嵌入 | Text embedding
├── base/                # 基础示例 | Basic examples
│   ├── generics/        # 泛型 | Generics
│   ├── network/         # 网络编程 | Network programming
│   └── web/             # Web开发 | Web development
├── cicd/                # CI/CD实践 | CI/CD practices
│   ├── docker/          # Docker容器化 | Docker containerization
│   ├── jenkins/         # Jenkins自动化 | Jenkins automation
│   ├── k8s/             # Kubernetes部署 | Kubernetes deployment
│   └── istio/           # Istio服务网格 | Istio service mesh
├── cloudnative/         # 云原生技术 | Cloud native technologies
├── concurrent/          # 并发编程专题 | Concurrent programming topics
├── container/           # 容器化技术 | Containerization
├── crontab/             # 定时任务 | Cron jobs
├── custom-web/          # 自定义Web框架 | Custom web framework
├── effective/           # Effective Go实践 | Effective Go practices
├── gedis/               # Redis实现 | Redis implementation
├── gin/                 # Gin Web框架 | Gin web framework
├── micservices/         # 微服务架构 | Microservices architecture
├── middleware/          # 中间件开发 | Middleware development
├── micro-with-containerization/ # 微服务容器化 | Microservice containerization
├── playground/          # 代码实验场 | Code playground
└── thewaytogo/          # Go学习路径 | Go learning path
```

## 技术栈 | Technical Stack

- **主要语言** | Primary Language: Go 1.24
- **文档格式** | Documentation: Markdown with automated code examples
- **CI/CD** | GitHub Actions with build verification
- **测试** | Testing: Extensive unit/benchmark tests with coverage reporting
- **依赖管理** | Dependency Management: Go Modules with vulnerability scanning

## 快速开始 | Quick Start

### 环境要求 | Prerequisites

- Go 1.24 或更高版本 | Go 1.24 or higher
- Git
- Docker (可选 | Optional)

### 安装 | Installation

```bash
# 克隆仓库 | Clone repository
git clone https://github.com/LCY2013/thinking-in-go.git
cd thinking-in-go

# 安装依赖 | Install dependencies
go mod tidy

# 运行设置脚本 | Run setup script
make setup

# 运行测试 | Run tests
make ut

# 代码格式化 | Format code
make fmt

# 代码检查 | Lint code
make lint
```

### 运行示例 | Running Examples

```bash
# 运行基础示例 | Run basic examples
cd action/02_quickstart/sample
go run main.go

# 运行Gin Web服务器 | Run Gin web server
cd gin/basic/v1/main
go run main.go

# 运行Redis服务器 | Run Redis server
cd gedis
go run main.go
```

## 主要模块详解 | Detailed Module Description

### 🎯 Action - 基础教程 | Basic Tutorials
包含从Go语言基础到高级特性的完整学习路径，每个模块都有详细的代码示例和说明。

Contains a complete learning path from Go language basics to advanced features, with detailed code examples and explanations for each module.

### 🚀 Advanced - 高级主题 | Advanced Topics
涵盖异步编程、高级并发、错误处理、内存重排序等高级主题。

Covers advanced topics such as asynchronous programming, advanced concurrency, error handling, and memory reordering.

### 🤖 AI - AI集成 | AI Integration
集成Ollama等AI工具，展示Go语言在AI领域的应用。

Integrates AI tools like Ollama, demonstrating Go language applications in the AI field.

### 🌐 Gin - Web框架 | Web Framework
完整的Gin Web框架使用示例，包括路由、中间件、认证等。

Complete Gin web framework usage examples, including routing, middleware, authentication, etc.

### 🗄️ Gedis - Redis实现 | Redis Implementation
用Go语言实现的Redis服务器，包含RESP协议解析、命令处理等核心功能。

Redis server implemented in Go, including RESP protocol parsing, command handling, and other core features.

### 🏗️ Microservices - 微服务架构 | Microservices Architecture
完整的微服务架构示例，包括服务发现、API网关、RPC通信等。

Complete microservices architecture examples, including service discovery, API gateway, RPC communication, etc.

### 🔧 Middleware - 中间件开发 | Middleware Development
各种中间件的实现，包括并发控制、反射、模板生成等。

Implementation of various middleware, including concurrency control, reflection, template generation, etc.

### 🐳 CI/CD - 持续集成部署 | Continuous Integration and Deployment
完整的CI/CD实践，包括Docker、Kubernetes、Jenkins、Istio等。

Complete CI/CD practices, including Docker, Kubernetes, Jenkins, Istio, etc.

## 学习路径 | Learning Path

### 初学者路径 | Beginner Path
1. **action/01_introduce** - Go语言介绍
2. **action/02_quickstart** - 快速开始
3. **action/04_array_slice_map** - 基础数据类型
4. **action/05_go_types_system** - 类型系统
5. **action/10_struct** - 结构体
6. **action/11_func** - 函数
7. **action/12_method** - 方法

### 进阶路径 | Intermediate Path
1. **action/06_concurrent** - 并发编程
2. **action/07_concurrent_mode** - 并发模式
3. **action/14_interface** - 接口
4. **action/15_concurrent** - 并发进阶
5. **advanced/** - 高级主题

### 实战路径 | Practical Path
1. **gin/** - Web开发
2. **gedis/** - 系统编程
3. **micservices/** - 微服务架构
4. **middleware/** - 中间件开发
5. **cicd/** - DevOps实践

## 贡献指南 | Contributing

我们欢迎各种形式的贡献！| We welcome contributions in various forms!

### 如何贡献 | How to Contribute

1. **Fork** 这个仓库 | Fork this repository
2. **创建** 你的特性分支 | Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. **提交** 你的更改 | Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. **推送** 到分支 | Push to the branch (`git push origin feature/AmazingFeature`)
5. **打开** 一个Pull Request | Open a Pull Request

### 代码规范 | Code Standards

- 遵循Go官方代码规范 | Follow Go official code standards
- 使用`gofmt`格式化代码 | Use `gofmt` to format code
- 编写单元测试 | Write unit tests
- 添加适当的注释 | Add appropriate comments

## 许可证 | License

本项目采用Apache 2.0许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.

## 目标受众 | Target Audience

- **基础设施工程师** 转向Go | Infrastructure engineers transitioning to Go
- **云原生应用开发者** | Cloud-native application developers  
- **SRE团队** 构建自定义工具 | SRE teams building custom tooling
- **计算机科学学生** 学习系统编程 | Computer science students studying systems programming
- **技术负责人** 架构Go平台 | Tech leads architecting Go-based platforms

## 独特价值 | Unique Value Proposition

与现有的Go书籍不同，它们要么专注于语言基础，要么专注于特定领域，这本百科全书将深度技术内容与整个云原生栈的生产验证实现相结合，为专业Go开发者提供单一权威参考。

Unlike existing Go books that focus on either language basics or specific domains, this encyclopedia will combine deep technical content with production-proven implementations across the entire cloud-native stack, serving as a single authoritative reference for professional Go developers.

## 项目目标 | Project Goals

- **创建** 随Go生态系统发展的活文档 | Create living documentation that evolves with the Go ecosystem
- **包含** 带有验证输出的可执行代码示例 | Include executable code samples with verified outputs
- **涵盖** 理论CS概念和实际实现 | Cover both theoretical CS concepts and practical implementations
- **维护** 与主要云平台的兼容性 | Maintain compatibility with major cloud platforms
- **提供** 从其他语言的迁移指南 | Provide migration guides from other languages

## 联系方式 | Contact

- **GitHub**: [@LCY2013](https://github.com/LCY2013)
- **Issues**: [GitHub Issues](https://github.com/LCY2013/thinking-in-go/issues)
- **Discussions**: [GitHub Discussions](https://github.com/LCY2013/thinking-in-go/discussions)

---

⭐ 如果这个项目对你有帮助，请给它一个星标！| If this project helps you, please give it a star!

📚 **让我们一起构建最全面的Go语言学习资源！** | **Let's build the most comprehensive Go language learning resource together!**