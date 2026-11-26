.PHONY: run build test lint migrate seed clean help

APP_NAME = siarsip-api
BUILD_DIR = ./bin

help: ## Menampilkan bantuan
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Menjalankan server (development, dengan live reload)
	go run ./cmd/server/main.go

build: ## Build binary
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server/main.go
	@echo "✅ Build selesai: $(BUILD_DIR)/$(APP_NAME)"

test: ## Menjalankan semua test
	go test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Test selesai"

test-cover: ## Melihat coverage report
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

lint: ## Menjalankan linter
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run ./... || echo "⚠️  golangci-lint tidak terinstall, jalankan: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
	go vet ./...

fmt: ## Format kode dengan gofmt
	gofmt -l -w .
	@echo "✅ Format kode selesai"

migrate: ## Menjalankan migrasi database
	go run ./cmd/server/main.go --migrate

seed: ## Menjalankan seeder (buat akun default)
	go run ./cmd/server/main.go --seed

clean: ## Membersihkan build artifacts
	rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "✅ Bersih"

deps: ## Install dependencies
	go mod tidy
	go mod download
	@echo "✅ Dependencies terinstall"

audit: ## Audit keamanan dependencies
	go mod verify
	@echo "✅ Module verification selesai"
