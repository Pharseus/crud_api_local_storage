# CRUD API Local Storage

## 🛠️ Tech Stack

- **Language**: Go 1.21+
- **Router**: Chi v5
- **Documentation**: Swaggo/Swag
- **Storage**: Local File System


## 🔧 Installation

### 1. Clone Repository
```bash
git clone <https://github.com/Pharseus/DMS_REST_API_FILES.git>
cd crud_api_local_storage
```

### 2. Install Dependencies
```bash
go mod download
go mod tidy
```

### 3. Install Swag CLI
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 4. Generate Swagger Documentation
```bash
swag init
```

### 5. Setup Storage Directory
Buat folder untuk menyimpan file (contoh: `D:/storage-local`), bisa edit path di `controllers/files_controller.go`:
```go
const storageRoot = "D:/storage-local"
```

### 6. Run Server
```bash
go run main.go
```

## 📚 API Documentation

Akses Swagger UI untuk dokumentasi interaktif:
```
http://localhost:8080/swagger/index.html
```
