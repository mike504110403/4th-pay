FROM golang:1.23-alpine AS builder

# 安裝 git
RUN apk add --no-cache git

# 設置 Go 環境變量以允許私有模組
RUN go env -w GOPRIVATE=ec2-15-168-3-237.ap-northeast-3.compute.amazonaws.com

# modules token
ARG MODULES_TOKEN

# 設置 Git 認證輔助程序來自動輸入憑證
RUN git config --global credential.helper '!f() { echo "username=git"; echo "password=${MODULES_TOKEN}"; }; f'

# 設置 Git URL 替換
RUN git config --global url."http://${MODULES_TOKEN}@ec2-15-168-3-237.ap-northeast-3.compute.amazonaws.com".insteadOf "http://ec2-15-168-3-237.ap-northeast-3.compute.amazonaws.com"

# 設置 GOINSECURE 環境變量以允許 http 連接
ENV GOINSECURE=ec2-15-168-3-237.ap-northeast-3.compute.amazonaws.com

# 設置工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 文件並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製其餘的應用程序文件
COPY . .

# 構建應用程序
RUN go build -o main .

# 使用輕量級的 Alpine 作為運行鏡像
FROM alpine:3.18

# 安裝 tzdata 以設置時區
RUN apk add --no-cache tzdata

# 設置時區為 Asia/Shanghai（UTC+8）
ENV TZ=Asia/Shanghai

# 設置工作目錄
WORKDIR /app

# 複製編譯好的應用程序到運行鏡像
COPY --from=builder /app/main .

# 確保可執行文件具有執行權限
RUN chmod +x ./main

# 設置容器啟動命令
CMD ["./main"]
