# ChatGPT API Proxy

ChatGPT API Proxy是一个基于Golang和Gin框架的API代理服务，旨在提供对ChatGPT服务的转发和增强功能。ChatGPT是一个自然语言处理模型，可以根据用户输入生成自然的回复，而ChatGPT
API Proxy可以帮助您更方便地使用ChatGPT服务。

# 功能

ChatGPT API Proxy支持以下功能：

### 目前正在开发中

- [x] 转发ChatGPT请求：您可以将来自客户端的ChatGPT请求发送到ChatGPT API Proxy，然后将请求转发到ChatGPT服务。支持基于API和基于USE认证形式的请求。
- [x] 流式API返回：ChatGPT API Proxy支持流式API返回，可以更快地返回结果，提高用户体验。
- [ ] Token使用计量：ChatGPT API Proxy支持对Token的使用计量，以及价格计算，方便您管理和计费。
- [ ] API日志监控：ChatGPT API Proxy提供API日志监控功能，可以帮助您追踪和分析API使用情况。
- [ ] Prompt查询：ChatGPT API Proxy支持Prompt查询功能，可以帮助您更方便地查询和管理Prompt。

# 使用方法

## API请求

### 请求地址

目前支持的API请求地址如下：

GPT3:  `${host}/api/openai/completions`

GPT3.5, 4: `${host}/api/openai/chat`

### 请求参数

所有参数和openai官方文档一致，具体请参考[openai官方文档](https://beta.openai.com/docs/api-reference/completions/create)。

请求Header中可以包含`Authorization`字段，其值为`Bearer ${token}`，`${token}`为您的OpenAI API Token。

如果不传，默认使用部署时使用的环境变量`OPENAI_API_KEY`。

# 部署

## Railway部署

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/new/template/ZweBXA)

在railway中需要配置的环境变量有：

```markdown
OPENAI_API_KEY
PORT
```

## Vercel部署

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fvercel%2Fnext.js%2Ftree%2Fcanary%2Fexamples%2Fhello-world)

在vercel中需要配置的环境变量有：

```markdown
OPENAI_API_KEY
```

# 本地开发

## 环境要求

- Golang 1.18+
- docker
- docker-compose

## golang技术栈

- Gin Web Framework
- Gorm ORM
- golang-migrate 数据库迁移
- viper 配置文件管理
- logrus 日志管理
- golangci-lint 代码检查

# 贡献

欢迎您参与贡献ChatGPT API Proxy项目，包括提交Bug报告、提出改进建议、编写文档、编写代码等。请参考项目贡献指南，了解如何参与贡献。

# 许可证

ChatGPT API Proxy采用MIT许可证开源，您可以自由地使用、复制、修改、合并、出版发行、散布、再授权和/或销售本软件的副本。详见LICENSE文件。
