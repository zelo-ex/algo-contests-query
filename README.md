# algo-contests-query

## 简介

Go 练手小项目，抓取平时常访问的 OJ 平台的比赛信息并汇总，目前正在增加新功能。

## 项目结构

```
algo-contests-query/
├── README.md          # 项目说明
├── go.mod             # Go 模块定义
├── go.sum             # 依赖校验和
├── main.go            # 程序入口，并发抓取并汇总输出
├── contests.go        # Contest、ContestsList 类型及 Filter/Sort 等方法
├── atcoder.go         # AtCoder 抓取与解析
├── codeforces.go      # Codeforces API 请求与解析
├── leetcode.go        # LeetCode GraphQL 请求与解析
├── luogu.go           # Luogu JSON 提取与解析
└── nowcoder.go        # Nowcoder JSON 抓取与解析
```

## 运行方式

```bash
# Clone the repo and go to the repo directory
git clone --depth=1 https://github.com/zelo-ex/algo-contests-query
cd algo-contests-query
# Download dependencies and run source files (Go 1.22+)
go run .
```
