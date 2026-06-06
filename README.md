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
