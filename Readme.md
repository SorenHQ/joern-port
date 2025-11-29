
# Joern Lab 🧠

A lightweight Go tool for interacting with a running [Joern](https://joern.io) server and fetching source code from GitHub repositories.

## Features

- 🔌 **Proxy** — Connect seamlessly to a running Joern server  
- 📄 **StdOutParser** — Parse Joern query output easily  
- ☁️ **GitHub Service** — Download repositories directly from GitHub  

---

## 🧰 Setup

### 1. Download Joern
Get the latest Joern release from:  
👉 [https://github.com/joernio/joern/releases](https://github.com/joernio/joern/releases)

### 2. Run Joern in Server Mode
```bash
./joern --server --server-host localhost --server-port 8081
````

### 3. Run Your Go App

```bash
go mod tidy
go run app.go
```


---

## ⚙️ Requirements

* Go 1.25+
* Running Joern server
* GitHub access token (if needed for private repos)




