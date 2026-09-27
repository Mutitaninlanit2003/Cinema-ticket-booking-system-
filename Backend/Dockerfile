# ใช้ Go เวอร์ชัน 1.24
FROM golang:1.24

# กำหนด Working Directory
WORKDIR /app

# Copy go.mod และ go.sum ก่อน
COPY go.mod go.sum ./

# ดาวน์โหลด dependencies
RUN go mod download

# Copy source code ทั้งหมด
COPY . .

# Build Application
RUN go build -o main ./cmd

# เปิด Port
EXPOSE 8080

# รันโปรแกรม
CMD ["./main"]