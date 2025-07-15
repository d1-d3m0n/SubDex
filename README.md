# 🕵️ Subdomain Enumerator

A fast and efficient multithreaded subdomain enumeration tool written in Go. It uses custom wordlists to enumerate subdomains for a given target domain using DNS queries.

![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)

---

## 🚀 Features

- ✅ Uses custom DNS resolvers (default: `8.8.8.8`)
- ✅ Multithreaded (concurrent workers)
- ✅ Supports CNAME resolution
- ✅ Outputs results to `.csv` file
- ✅ Pretty console output

---


## Tech Stack

**Server:** Golang

## 📸 Demo


## Installation

Install my-project with npm

```bash
  git clone https://github.com/d1_d3m0n/subdex.git
  cd subdex
  go build -o subenum
  ./subenum -domain example.com -wordlist subdomains.txt -c 200
```

## Screenshots

![App Screenshot](screenshot.jpg)



## Example Output
```bash
  Output saved as: example.com.csv
  www.example.com    93.184.216.34
  blog.example.com   192.0.2.10
  cdn.example.com    203.0.113.5
```

## Author
```bash
  Made with ❤️ by d1_d3m0n
```

## 🔒 Disclaimer
This tool is intended for educational and authorized security testing purposes only. Unauthorized scanning or probing of networks you don't own is illegal and unethical.

## 📜 License

MIT License. © 2025 d1_d3m0n
