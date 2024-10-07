# Tahap pertama: Pembangunan
FROM golang:1.22.1 AS builder

# Set lingkungan kerja
WORKDIR /app

# Mengcopy go.mod dan go.sum agar dependensi dapat di-cache
COPY go.mod go.sum ./

# Mengunduh dependensi
RUN go mod download

# Mengcopy kode sumber ke dalam kontainer
COPY . .

# Membuat binary aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/myapp ./cmd/app

# Debugging: Cek apakah binary berhasil dibuat
RUN ls -l /app/myapp

# Tahap kedua: Produksi
FROM alpine:latest

# Instal tzdata untuk dukungan zona waktu
RUN apk --no-cache add tzdata

WORKDIR /app

# Mengcopy binary dari tahap pembangunan
COPY --from=builder /app/myapp .
COPY --from=builder /app/.env .

# Konfigurasi zona waktu
ENV TZ=Asia/Jakarta
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

# Membuat direktori /assets/images dan mengatur izin
RUN mkdir -p /app/assets/images && chmod 777 /app/assets/images

# Menambahkan logging izin direktori
RUN ls -l /app/assets && ls -l /app/assets/images

# Menjalankan aplikasi
EXPOSE 8080

CMD ["./myapp"]
