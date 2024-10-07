# File: rename_migration.ps1

# Ambil timestamp saat ini (format: YYYYMMDDHHMMSS)
$timestamp = Get-Date -Format "yyyyMMddHHmmss"

# Nama migrasi yang diambil dari argumen pertama
$migration_name = $args[0]

# Jalankan perintah migrate create
migrate create -ext sql -dir ./db/migrations -seq $migration_name

# Dapatkan nama file yang baru dibuat
$up_file = Get-ChildItem ./db/migrations/*_$migration_name.up.sql -ErrorAction SilentlyContinue | Select-Object -First 1
$down_file = Get-ChildItem ./db/migrations/*_$migration_name.down.sql -ErrorAction SilentlyContinue | Select-Object -First 1

# Periksa apakah file ditemukan sebelum mengganti nama
if ($up_file -and $down_file) {
    # Ganti nama file dengan format timestamp
    Rename-Item $up_file.FullName "./db/migrations/${timestamp}_$migration_name.up.sql"
    Rename-Item $down_file.FullName "./db/migrations/${timestamp}_$migration_name.down.sql"
    Write-Host "Migrasi dengan nama timestamp ${timestamp}_$migration_name telah dibuat."
} else {
    Write-Host "File migrasi tidak ditemukan. Pastikan perintah migrate create berhasil dijalankan."
}
