run:
	@echo "running server..."
	go run ./cmd/main.go

BackupTest:
	@echo "running backup test..."
	go test ./server/handler/ -run TestBackup -v

TestBackup_WriteToFile:
	@echo "running backup test..."
	go test ./server/handler/ -run TestBackup_WriteToFile -v