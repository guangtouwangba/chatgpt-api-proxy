# ChatGPT API Proxy
ChatGPT API Proxy是一个基于Golang和Gin框架的API代理服务，旨在提供对ChatGPT服务的转发和增强功能。ChatGPT是一个自然语言处理模型，可以根据用户输入生成自然的回复，而ChatGPT
API Proxy可以帮助您更方便地使用ChatGPT服务。

# 功能
ChatGPT API Proxy支持以下功能：

- [ ] 转发ChatGPT请求：您可以将来自客户端的ChatGPT请求发送到ChatGPT API Proxy，然后将请求转发到ChatGPT服务。支持基于API和基于USE认证形式的请求。
- [ ] 流式API返回：ChatGPT API Proxy支持流式API返回，可以更快地返回结果，提高用户体验。
- [ ] Token使用计量：ChatGPT API Proxy支持对Token的使用计量，以及价格计算，方便您管理和计费。
- [ ] API日志监控：ChatGPT API Proxy提供API日志监控功能，可以帮助您追踪和分析API使用情况。
- [ ] Prompt查询：ChatGPT API Proxy支持Prompt查询功能，可以帮助您更方便地查询和管理Prompt。

# 使用方法

请参考项目文档，了解如何安装和配置ChatGPT API Proxy，以及如何使用API和命令行工具。

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

# 贡献

欢迎您参与贡献ChatGPT API Proxy项目，包括提交Bug报告、提出改进建议、编写文档、编写代码等。请参考项目贡献指南，了解如何参与贡献。

# 许可证

ChatGPT API Proxy采用MIT许可证开源，您可以自由地使用、复制、修改、合并、出版发行、散布、再授权和/或销售本软件的副本。详见LICENSE文件。
