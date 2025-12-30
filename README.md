# Exercise API Server

一個使用 Go 和 Gin 框架實作的簡單 Task 管理 REST API 服務。

## 專案架構

```
exercise-api-server/
├── cmd/exercise-api-server/    # 應用程式入口點
├── internal/
│   ├── api/v1/                # API handlers
│   ├── cache/                 # 快取實作
│   ├── domain/                # 領域模型
│   ├── repository/memory/     # 記憶體儲存實作
│   └── task/                  # 業務邏輯服務
└── docker/                    # Docker 配置
```

## API 端點

| 方法 | 路徑 | 描述 |
|------|------|------|
| GET | `/api/v1/tasks` | 取得所有任務 |
| GET | `/api/v1/tasks/:id` | 取得特定任務 |
| POST | `/api/v1/tasks` | 創建新任務 |
| PUT | `/api/v1/tasks/:id` | 更新任務 |
| DELETE | `/api/v1/tasks/:id` | 刪除任務 |
| GET | `/health` | 健康檢查 |

## Task 資料結構

```json
{
  "id": "1",
  "name": "任務名稱",
  "status": 0,
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

- `status`: 0 = 未完成, 1 = 已完成

## 快速開始

### 本地執行

1. 確保已安裝 Go 1.19+
2. 複製專案：
   ```bash
   git clone <repository-url>
   cd exercise-api-server
   ```
3. 安裝依賴：
   ```bash
   go mod download
   ```
4. 執行服務：
   ```bash
   go run cmd/exercise-api-server/main.go
   ```

服務將在 `http://localhost:8080` 啟動。

### 使用 Docker

```bash
docker build -f docker/Dockerfile -t exercise-api-server .
docker run -p 8080:8080 exercise-api-server
```

## API 使用範例

### 創建任務
```bash
curl -X POST http://localhost:8080/api/v1/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "完成專案", "status": 0}'
```

### 取得所有任務
```bash
curl http://localhost:8080/api/v1/tasks
```

### 更新任務
```bash
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "完成專案", "status": 1}'
```

## 測試

執行單元測試：
```bash
go test ./...
```

## API 文件

Swagger API 文件位於 `docs/swagger.yaml`。你可以使用以下方式查看：

1. **線上 Swagger Editor**: 將 `docs/swagger.yaml` 內容貼到 [Swagger Editor](https://editor.swagger.io/)
2. **本地 Swagger UI**: 使用 Docker 執行
   ```bash
   docker run -p 8081:8080 -v $(pwd)/docs:/usr/share/nginx/html/docs swaggerapi/swagger-ui
   ```
   然後在瀏覽器開啟 `http://localhost:8081/?url=/docs/swagger.yaml`

## 技術特色

- **Clean Architecture**: 採用分層架構設計
- **記憶體快取**: 使用自實作的記憶體快取系統
- **RESTful API**: 遵循 REST 設計原則
- **錯誤處理**: 完整的錯誤回應機制
- **單元測試**: 包含完整的測試覆蓋