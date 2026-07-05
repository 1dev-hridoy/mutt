package handler

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/internal/middleware"
	"github.com/dishan1223/mutt/internal/service"
	"github.com/dishan1223/mutt/models"
	"github.com/gofiber/fiber/v3"
)

func newTestAppWithBackup() *fiber.App {
	app := fiber.New()

	app.Post("/api/v1/auth/signup", SignUpHandler)
	app.Post("/api/v1/auth/login", LoginHandler)

	projects := app.Group("/api/v1/projects", middleware.AuthRequired)
	projects.Post("/", CreateProjectHandler)

	app.Get("/api/v1/backup", middleware.AuthRequired, ExportBackupHandler)
	app.Post("/api/v1/backup/import", middleware.AuthRequired, ImportBackupHandler)
	return app
}

func loginBackupUser(t *testing.T, app *fiber.App, email string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"email": email, "password": "password123"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "access_token" {
			return cookie.Value
		}
	}
	t.Fatal("access_token cookie not found")
	return ""
}

func createBackupProject(t *testing.T, app *fiber.App, token, name string) uint {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"name": name})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/projects/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ := app.Test(req)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return uint(result["id"].(float64))
}

func seedBackupData(t *testing.T, projectID uint) {
	t.Helper()

	group1 := models.ErrorGroup{
		ProjectID:   projectID,
		Fingerprint: service.ComputeFingerprint("stack1", "NilPointerError"),
		Title:       "NilPointerError",
		Status:      "critical",
		Count:       2,
		LastSeenAt:  time.Now(),
	}
	config.DB.Create(&group1)

	group2 := models.ErrorGroup{
		ProjectID:   projectID,
		Fingerprint: service.ComputeFingerprint("stack2", "TimeoutError"),
		Title:       "TimeoutError",
		Status:      "critical",
		Count:       1,
		LastSeenAt:  time.Now(),
	}
	config.DB.Create(&group2)

	config.DB.Create(&models.Error{
		ErrorGroupID: group1.ID,
		ProjectID:    projectID,
		Log:          "goroutine 1 [running]: main.foo()",
		StackTrace:   "stack1",
		Severity:     "error",
		OccurredAt:   time.Now(),
	})
	config.DB.Create(&models.Error{
		ErrorGroupID: group1.ID,
		ProjectID:    projectID,
		Log:          "goroutine 1 [running]: main.foo()",
		StackTrace:   "stack1",
		Severity:     "error",
		OccurredAt:   time.Now(),
	})
	config.DB.Create(&models.Error{
		ErrorGroupID: group2.ID,
		ProjectID:    projectID,
		Log:          "context deadline exceeded",
		StackTrace:   "stack2",
		Severity:     "error",
		OccurredAt:   time.Now(),
	})
}

func TestBackup_ExportAndReimport(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "backupuser", "backup@example.com", "password123", "+1111111111")
	token := loginBackupUser(t, app, "backup@example.com")

	// create project and seed data directly
	projectID := createBackupProject(t, app, token, "Backup Test Project")
	seedBackupData(t, projectID)

	// verify data seeded
	var groups []models.ErrorGroup
	config.DB.Where("project_id = ?", projectID).Find(&groups)
	if len(groups) != 2 {
		t.Fatalf("expected 2 error groups, got %d", len(groups))
	}

	// export backup
	req := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("export request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("export expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)

	var backup models.BackupData
	if err := json.Unmarshal(respBody, &backup); err != nil {
		t.Fatalf("failed to parse backup JSON: %v", err)
	}
	if len(backup.Projects) != 1 {
		t.Fatalf("expected 1 project in backup, got %d", len(backup.Projects))
	}
	if backup.Projects[0].Name != "Backup Test Project" {
		t.Fatalf("expected project name 'Backup Test Project', got %v", backup.Projects[0].Name)
	}
	if len(backup.Projects[0].ErrorGroups) != 2 {
		t.Fatalf("expected 2 error groups in backup, got %d", len(backup.Projects[0].ErrorGroups))
	}

	// re-import as a different user
	seedUser(t, "importuser", "import@example.com", "password123", "+2222222222")
	importToken := loginBackupUser(t, app, "import@example.com")

	// build multipart form
	importBody, _ := json.Marshal(backup)
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "backup.json")
	part.Write(importBody)
	writer.Close()

	importReq := httptest.NewRequest(http.MethodPost, "/api/v1/backup/import", &buf)
	importReq.Header.Set("Authorization", "Bearer "+importToken)
	importReq.Header.Set("Content-Type", writer.FormDataContentType())

	importResp, err := app.Test(importReq)
	if err != nil {
		t.Fatalf("import request failed: %v", err)
	}
	if importResp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(importResp.Body)
		t.Fatalf("import expected 200, got %d: %s", importResp.StatusCode, string(body))
	}

	var importResult map[string]interface{}
	json.NewDecoder(importResp.Body).Decode(&importResult)
	imported := importResult["imported"].(map[string]interface{})
	if imported["projects"].(float64) != 1 {
		t.Fatalf("expected 1 imported project, got %v", imported["projects"])
	}
	if imported["error_groups"].(float64) != 2 {
		t.Fatalf("expected 2 imported error groups, got %v", imported["error_groups"])
	}
	if imported["errors"].(float64) != 3 {
		t.Fatalf("expected 3 imported errors, got %v", imported["errors"])
	}

	// verify import user has their own project
	var importUser models.User
	config.DB.Where("email = ?", "import@example.com").First(&importUser)
	var userProjects []models.Project
	config.DB.Where("user_id = ?", importUser.ID).Find(&userProjects)
	if len(userProjects) != 1 {
		t.Fatalf("expected 1 project for import user, got %d", len(userProjects))
	}
}

func TestBackup_ExportNoAuth(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestBackup_ExportEmpty(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "emptybackup", "empty@example.com", "password123", "+3333333333")
	token := loginBackupUser(t, app, "empty@example.com")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var backup models.BackupData
	json.NewDecoder(resp.Body).Decode(&backup)
	if len(backup.Projects) != 0 {
		t.Fatalf("expected 0 projects, got %d", len(backup.Projects))
	}
}

func TestBackup_ImportNoFile(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "importno", "importno@example.com", "password123", "+4444444444")
	token := loginBackupUser(t, app, "importno@example.com")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/backup/import", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestBackup_ImportInvalidJSON(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "importbad", "importbad@example.com", "password123", "+5555555555")
	token := loginBackupUser(t, app, "importbad@example.com")

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, _ := writer.CreateFormFile("file", "backup.json")
	part.Write([]byte("not valid json"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/backup/import", &buf)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestBackup_ExportGzipCompression(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "gzipuser", "gzip@example.com", "password123", "+6666666666")
	token := loginBackupUser(t, app, "gzip@example.com")

	projectID := createBackupProject(t, app, token, "Gzip Project")

	// seed enough data to exceed gzip threshold (10KB)
	for i := 0; i < 50; i++ {
		group := models.ErrorGroup{
			ProjectID:   projectID,
			Fingerprint: service.ComputeFingerprint("stack"+string(rune(i)), "BigError"),
			Title:       "BigError",
			Status:      "critical",
			Count:       1,
			LastSeenAt:  time.Now(),
		}
		config.DB.Create(&group)
		config.DB.Create(&models.Error{
			ErrorGroupID: group.ID,
			ProjectID:    projectID,
			Log:          "goroutine 1 [running]: main.bigFunction()\nruntime.goexit()\ngoroutine 2 [running]: main.otherFunction()\nruntime.goexit()",
			StackTrace:   "stack" + string(rune(i)),
			Severity:     "error",
			OccurredAt:   time.Now(),
		})
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)

	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(bytes.NewReader(respBody))
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer gz.Close()
		decompressed, _ := io.ReadAll(gz)
		var backup models.BackupData
		if err := json.Unmarshal(decompressed, &backup); err != nil {
			t.Fatalf("failed to parse decompressed backup: %v", err)
		}
		if len(backup.Projects) != 1 {
			t.Fatalf("expected 1 project, got %d", len(backup.Projects))
		}
	} else {
		var backup models.BackupData
		if err := json.Unmarshal(respBody, &backup); err != nil {
			t.Fatalf("failed to parse backup: %v", err)
		}
		if len(backup.Projects) != 1 {
			t.Fatalf("expected 1 project, got %d", len(backup.Projects))
		}
	}
}

func TestBackup_UserIsolation(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()

	seedUser(t, "isolatedA", "isoA@example.com", "password123", "+1111111111")
	tokenA := loginBackupUser(t, app, "isoA@example.com")
	projectA := createBackupProject(t, app, tokenA, "User A Project")

	groupA := models.ErrorGroup{
		ProjectID:   projectA,
		Fingerprint: service.ComputeFingerprint("stackA", "ErrorA"),
		Title:       "ErrorA",
		Status:      "critical",
		Count:       1,
		LastSeenAt:  time.Now(),
	}
	config.DB.Create(&groupA)
	config.DB.Create(&models.Error{
		ErrorGroupID: groupA.ID,
		ProjectID:    projectA,
		Log:          "user A error",
		StackTrace:   "stackA",
		Severity:     "error",
		OccurredAt:   time.Now(),
	})

	seedUser(t, "isolatedB", "isoB@example.com", "password123", "+2222222222")
	tokenB := loginBackupUser(t, app, "isoB@example.com")
	projectB := createBackupProject(t, app, tokenB, "User B Project")

	groupB := models.ErrorGroup{
		ProjectID:   projectB,
		Fingerprint: service.ComputeFingerprint("stackB", "ErrorB"),
		Title:       "ErrorB",
		Status:      "critical",
		Count:       1,
		LastSeenAt:  time.Now(),
	}
	config.DB.Create(&groupB)
	config.DB.Create(&models.Error{
		ErrorGroupID: groupB.ID,
		ProjectID:    projectB,
		Log:          "user B error",
		StackTrace:   "stackB",
		Severity:     "error",
		OccurredAt:   time.Now(),
	})

	// user A should only see their project
	reqA := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	reqA.Header.Set("Authorization", "Bearer "+tokenA)
	respA, _ := app.Test(reqA)
	var backupA models.BackupData
	json.NewDecoder(respA.Body).Decode(&backupA)
	if len(backupA.Projects) != 1 {
		t.Fatalf("user A expected 1 project, got %d", len(backupA.Projects))
	}
	if backupA.Projects[0].Name != "User A Project" {
		t.Fatalf("expected 'User A Project', got %v", backupA.Projects[0].Name)
	}

	// user B should only see their project
	reqB := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	reqB.Header.Set("Authorization", "Bearer "+tokenB)
	respB, _ := app.Test(reqB)
	var backupB models.BackupData
	json.NewDecoder(respB.Body).Decode(&backupB)
	if len(backupB.Projects) != 1 {
		t.Fatalf("user B expected 1 project, got %d", len(backupB.Projects))
	}
	if backupB.Projects[0].Name != "User B Project" {
		t.Fatalf("expected 'User B Project', got %v", backupB.Projects[0].Name)
	}
}

func TestBackup_WriteToFile(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	app := newTestAppWithBackup()
	seedUser(t, "fileuser", "file@example.com", "password123", "+1111111111")
	token := loginBackupUser(t, app, "file@example.com")

	projectID := createBackupProject(t, app, token, "File Test Project")
	seedBackupData(t, projectID)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/backup", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("export request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("export expected 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)

	// decompress if gzipped
	data := respBody
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(bytes.NewReader(respBody))
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer gz.Close()
		data, _ = io.ReadAll(gz)
	}

	// pretty-print for readability
	var backup models.BackupData
	json.Unmarshal(data, &backup)
	pretty, _ := json.MarshalIndent(backup, "", "  ")

	if err := os.WriteFile("backup_dump.json", pretty, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
	t.Logf("backup written to server/handler/backup_dump.json (%d bytes)", len(pretty))
}
