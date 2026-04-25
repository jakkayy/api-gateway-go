# API Gateway (Go)

API Gateway เขียนด้วย Go สำหรับ routing และ load balancing request ไปยัง upstream services โดยไม่พึ่ง framework ภายนอก ใช้แค่ Go standard library เป็นหลัก

## Features

- **Reverse Proxy** — ส่ง request ต่อไปยัง upstream service และ copy response กลับมาให้ client
- **Round-Robin Load Balancing** — กระจาย request ไปยัง target หลายตัวแบบผลัดเปลี่ยน
- **Rate Limiting** — จำกัด 5 request/นาที ต่อ IP
- **Request Logging** — บันทึก method, path, และเวลาที่ใช้ทุก request
- **Panic Recovery** — ดักจับ panic และตอบ 500 แทนการ crash

## โครงสร้างโปรเจค

```
api-gateway-go/
├── cmd/gateway/
│   └── main.go                  # จุดเริ่มต้นโปรแกรม
├── configs/
│   └── routes.yaml              # กำหนด route และ target services
├── internal/
│   ├── config/
│   │   └── config.go            # อ่านและ parse config file
│   ├── router/
│   │   └── router.go            # ประกอบ route + middleware + proxy
│   ├── balancer/
│   │   └── round.robin.go       # Round-robin load balancer
│   ├── proxy/
│   │   └── reverse_proxy.go     # Reverse proxy ไปยัง upstream
│   ├── middleware/
│   │   ├── chain.go             # ประกอบ middleware เป็น stack
│   │   ├── recovery.go          # ดักจับ panic
│   │   ├── logging.go           # บันทึก request log
│   │   └── ratelimit.go         # จำกัด request ต่อ IP
│   └── service/
│       └── upstream.go          # Upstream struct (reserved)
└── pkg/
    └── logger/
        └── logger.go            # Logger wrapper
```

## เส้นทางการเดินทางของ Request

```
Client
  │
  ▼
HTTP Server :8080
  │
  ▼
Router  (จับคู่ path กับ handler)
  │
  ▼
Middleware Chain
  ├─ 1. Recovery    (ครอบสุดนอก — ดักจับ panic)
  ├─ 2. Logging     (จับเวลาและ log request)
  └─ 3. RateLimit   (เช็ค IP — ถ้าเกินตอบ 429)
  │
  ▼
Reverse Proxy
  │
  ▼
Round-Robin Load Balancer  (เลือก target)
  │
  ▼
Upstream Service  (user-service / order-service)
  │
  ▼ (response ย้อนกลับมา)
Client
```

## การตั้งค่า Route

แก้ไขได้ที่ [configs/routes.yaml](configs/routes.yaml)

```yaml
routes:
  - path: /users
    targets:
      - http://localhost:8001
      - http://localhost:8003   # request จะสลับกันระหว่างสองตัวนี้

  - path: /orders
    targets:
      - http://localhost:8002
```

- `path` — URL path ที่ gateway จะรับ
- `targets` — list ของ upstream service URL (load balance แบบ round-robin)

## วิธีรัน

### Prerequisites
- Go 1.25+

### รัน Gateway

```bash
cd api-gateway-go
go run cmd/gateway/main.go
```

Gateway จะเปิดที่ `http://localhost:8080`

### รัน Upstream Services (ตัวอย่าง)

```bash
# Terminal 1 — user-service (port 8001)
go run ../user-service/main.go

# Terminal 2 — user-service-2 (port 8003)
go run ../user-service-2/main.go

# Terminal 3 — order-service (port 8002)
go run ../order-service/main.go
```

### ทดสอบ

```bash
# ทดสอบ route /users
curl http://localhost:8080/users

# ทดสอบ load balancing (request สลับระหว่าง :8001 และ :8003)
curl http://localhost:8080/users
curl http://localhost:8080/users

# ทดสอบ route /orders
curl http://localhost:8080/orders
```

## Rate Limit

ค่า default: **5 request ต่อนาที ต่อ IP**

เมื่อเกิน limit จะได้รับ response:
```
HTTP 429 Too Many Requests
Too many requests
```

แก้ไขค่าได้ใน [internal/middleware/ratelimit.go](internal/middleware/ratelimit.go):

```go
limit  = 5           // จำนวน request สูงสุด
window = time.Minute // ช่วงเวลา
```

## Dependencies

| Package | การใช้งาน |
|---------|-----------|
| `gopkg.in/yaml.v3` | Parse config YAML |
| Go stdlib | HTTP server, reverse proxy, logging |
