# 🌐 NetPlay - DevOps Networking Demo

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![MySQL](https://img.shields.io/badge/MySQL-8.0+-4479A1?style=flat&logo=mysql)](https://mysql.com)
[![Nginx](https://img.shields.io/badge/Nginx-1.24+-009639?style=flat&logo=nginx)](https://nginx.com)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A comprehensive demonstration of **networking concepts in DevOps** including Reverse Proxy, Load Balancing, CORS, Network Isolation, Client IP Preservation, Rate Limiting, and Request Headers.

## 🎯 Networking Concepts Demonstrated

| Concept | Implementation | Why It Matters |
|---------|---------------|----------------|
| **Reverse Proxy** | Nginx routes API requests to backend | Hides backend infrastructure, provides single entry point |
| **Client IP Preservation** | `X-Forwarded-For`, `X-Real-IP` headers | Maintains original client identity through proxy chain |
| **CORS (Cross-Origin Resource Sharing)** | Configured at proxy level with proper headers | Enables secure cross-origin requests from browsers |
| **Load Balancing** | Upstream backend configuration with least_conn | Distributes traffic, improves reliability |
| **Rate Limiting** | 100 requests per 60 seconds per IP | Prevents API abuse and DDoS attacks |
| **Request/Response Headers** | Custom security and cache headers | Enhances security and performance |
| **Network Isolation** | Docker bridge network (when containerized) | Separates services for security |
| **Health Checks** | `/health` endpoint monitoring | Enables automated recovery and monitoring |
| **Access Logging** | Detailed request/response logging | Debugging, monitoring, and analytics |

## 🏗️ Architecture & Network Flow
# NetPlay - Networking Concepts with Go, MySQL & Nginx

A full-stack application designed to demonstrate real-world networking concepts including **reverse proxying, load balancing, CORS, rate limiting, client IP preservation, security headers, and network monitoring** using Go, MySQL, and Nginx.

---

## 📖 Overview

This project showcases how modern web applications are deployed behind an Nginx reverse proxy while maintaining security, scalability, and observability.

### Architecture

```text
┌─────────────┐
│  Browser    │
│   Client    │
└──────┬──────┘
       │ HTTP Request
       ▼
┌─────────────────────────────────────────────────────────────┐
│                NGINX Reverse Proxy (Port 80)               │
│                                                             │
│  Features:                                                  │
│  • Load Balancing (least_conn)                              │
│  • CORS Headers                                             │
│  • Rate Limiting (100 req/min)                              │
│  • Client IP Forwarding (X-Forwarded-For)                   │
│  • Static File Serving                                      │
│                                                             │
│  ┌──────────────┐ ┌──────────────┐ ┌──────────────┐         │
│  │   /api/*     │ │      /       │ │   /health    │         │
│  │ Proxy Pass   │ │    Static    │ │ Health Check │         │
│  └──────┬───────┘ └──────┬───────┘ └──────┬───────┘         │
└─────────┼────────────────────┼────────────────────┼─────────┘
          │                    │                    │
          ▼                    ▼                    ▼
   ┌─────────────┐      ┌─────────────┐      ┌─────────────┐
   │  Backend    │      │  Frontend   │      │   Health    │
   │   :8080     │      │ Static Files│      │   Monitor   │
   └──────┬──────┘      └─────────────┘      └─────────────┘
          │
          ▼
   ┌─────────────┐
   │    MySQL    │
   │    :3306    │
   └─────────────┘
```

---

# 📡 Network Headers Flow

## Browser Request

```http
GET /api/tasks
Host: localhost
Origin: http://localhost
```

## Nginx Adds Headers

```http
X-Forwarded-For: client-ip
X-Real-IP: client-ip
X-Request-ID: unique-id
```

## Backend Receives

```text
✓ Original client IP
✓ Request chain information
✓ Accurate logging data
```

## Response Headers

```http
Access-Control-Allow-Origin: *
X-Frame-Options: SAMEORIGIN
Cache-Control: public, immutable
```

---

# ✨ Features

* Reverse Proxy with Nginx
* Load Balancing using `least_conn`
* Rate Limiting (100 requests/minute)
* Client IP Preservation
* Security Headers
* CORS Configuration
* Static File Serving
* Health Monitoring Endpoint
* Request Logging
* Network Debugging Endpoints

---

# 🚀 Quick Start

## Prerequisites

* Go 1.21+
* MySQL 8.0+
* Nginx 1.24+

## Installation

### 1. Clone Repository

```bash
git clone https://github.com/yourusername/netplay.git
cd netplay
```

### 2. Setup Database

```bash
mysql -u root -p < sql/init.sql
```

### 3. Start Backend

```bash
cd backend

go mod download
go run cmd/main.go
```

Backend will start on:

```text
http://localhost:8080
```

### 4. Configure Nginx

```bash
sudo cp nginx/nginx.conf /etc/nginx/sites-available/netplay

sudo ln -s \
/etc/nginx/sites-available/netplay \
/etc/nginx/sites-enabled/

sudo nginx -t

sudo systemctl reload nginx
```

### 5. Open Application

```text
http://localhost
```

---

# 🔬 Testing Network Concepts

## 1. Test Reverse Proxy

### Direct Backend

```bash
curl http://localhost:8080/api/health
```

### Through Nginx

```bash
curl http://localhost/api/health
```

### Compare Headers

```bash
curl -I http://localhost/api/health

curl -I http://localhost:8080/api/health
```

Observe the forwarded headers added by Nginx.

---

## 2. Test Client IP Preservation

```bash
curl http://localhost/api/network-info
```

Expected response:

```json
{
  "client_ip": "x.x.x.x",
  "x_forwarded_for": "x.x.x.x",
  "x_real_ip": "x.x.x.x",
  "remote_addr": "127.0.0.1"
}
```

---

## 3. Test CORS Headers

```bash
curl -X OPTIONS http://localhost/api/tasks \
  -H "Origin: http://example.com" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

Verify:

```http
Access-Control-Allow-Origin: *
```

---

## 4. Test Rate Limiting

```bash
for i in {1..150}; do
  curl -s http://localhost/api/tasks > /dev/null
  echo "Request $i"
done
```

After 100 requests within the configured window:

```text
HTTP 429 Too Many Requests
```

---

## 5. Test Load Balancing

Start a second backend instance:

```bash
cd backend

PORT=8081 go run cmd/main.go
```

Update Nginx upstream configuration:

```nginx
upstream backend {
    least_conn;

    server 127.0.0.1:8080 weight=3;
    server 127.0.0.1:8081 weight=2;
}
```

Reload Nginx:

```bash
sudo nginx -t
sudo systemctl reload nginx
```

Test request distribution:

```bash
for i in {1..20}; do
  curl http://localhost/api/network-info | jq '.hostname'
done
```

---

## 6. Monitor Network Traffic

### Nginx Logs

```bash
tail -f /var/log/nginx/access.log
```

### Backend Logs

```bash
tail -f backend/logs/app.log
```

### Rate Limiting Events

```bash
tail -f /var/log/nginx/error.log | grep limiting
```

---

# 📊 API Endpoints

| Method | Endpoint            | Description       | Network Feature |
| ------ | ------------------- | ----------------- | --------------- |
| GET    | `/api/tasks`        | Get all tasks     | Load Balanced   |
| POST   | `/api/tasks`        | Create task       | Rate Limited    |
| GET    | `/api/tasks/{id}`   | Get single task   | CORS Enabled    |
| PUT    | `/api/tasks/{id}`   | Update task       | IP Preserved    |
| DELETE | `/api/tasks/{id}`   | Delete task       | Header Injected |
| GET    | `/api/network-info` | Network debugging | Shows headers   |
| GET    | `/api/health`       | Health check      | Monitoring      |

---

# 🧪 Network Debugging Tools

## Browser DevTools

Press **F12**

### Network Tab

Inspect:

* Request headers
* Response headers
* Response times
* Status codes

### Timing Tab

Measure:

* DNS lookup
* TCP handshake
* Proxy latency
* Server response time

### Console Tab

Check:

* CORS errors
* Network failures
* JavaScript request issues

---

## Command Line Tools

### DNS Resolution

```bash
nslookup localhost
```

### Route Tracing

```bash
tracert localhost
```

### Open Ports

```bash
netstat -an
```

### Linux Port Monitoring

```bash
lsof -i :80
lsof -i :8080
```

---

# 📈 Performance Impact

| Feature        | Overhead | Benefit                |
| -------------- | -------- | ---------------------- |
| Reverse Proxy  | ~5ms     | Security & Scalability |
| Rate Limiting  | <1ms     | Prevents Abuse         |
| Load Balancing | Minimal  | High Availability      |
| CORS Headers   | <0.5ms   | Browser Security       |
| Logging        | ~2ms     | Monitoring & Debugging |

---

# 🔒 Security Headers

```nginx
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;

add_header Access-Control-Allow-Origin * always;
add_header Access-Control-Allow-Methods "GET, POST, PUT, DELETE, OPTIONS";
```

---

# 🐛 Troubleshooting

## Common Issues

| Issue              | Cause                 | Solution                  |
| ------------------ | --------------------- | ------------------------- |
| 404 Not Found      | Incorrect proxy path  | Verify `proxy_pass`       |
| CORS Errors        | Missing headers       | Check CORS config         |
| HTTP 429           | Rate limit exceeded   | Wait for reset            |
| Wrong Client IP    | Headers not forwarded | Verify `proxy_set_header` |
| Connection Refused | Backend offline       | Start backend             |

## Useful Commands

### Check Listening Ports

```bash
netstat -tulpn
```

### Test Connectivity

```bash
telnet localhost 80

telnet localhost 8080
```

### View Logs

```bash
tail -f /var/log/nginx/access.log

tail -f /var/log/nginx/error.log
```

### Test Different User Agents

```bash
curl -A "Mobile" http://localhost/api/health
```

---

# 📚 Learning Outcomes

After completing this project you will understand:

* Reverse Proxy Architecture
* Load Balancing Strategies
* CORS Configuration
* Client IP Preservation
* Rate Limiting Techniques
* HTTP Headers
* Network Monitoring
* Security Hardening
* Nginx Request Processing
* Backend Traffic Routing

---

# 🛠️ Built With

* Go 1.21
* MySQL 8.0
* Nginx 1.24+
* HTML5
* CSS3
* JavaScript

---

# 🤝 Contributing

Contributions are welcome.

If you discover additional networking concepts or improvements, feel free to submit a Pull Request.

---

# 📄 License

This project is licensed under the MIT License.

See the `LICENSE` file for details.

---

# 🙏 Acknowledgments

* Nginx Documentation
* Go Community
* DevOps Community
* Open Source Contributors

Built with ❤️ to demonstrate networking concepts in a practical, real-world environment.
