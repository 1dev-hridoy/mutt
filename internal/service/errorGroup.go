package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/models"
	"gorm.io/gorm/clause"
)

func ComputeFingerprint(stackTrace string, title string) string {
	h := sha256.New()
	h.Write([]byte(title + "\n" + stackTrace))
	return hex.EncodeToString(h.Sum(nil))
}

func FindOrCreateErrorGroup(projectID uint, fingerprint, title string) (*models.ErrorGroup, error) {
	var group models.ErrorGroup
	now := time.Now()

	err := config.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "project_id"}, {Name: "fingerprint"}},
		DoUpdates: clause.AssignmentColumns([]string{"count", "last_seen_at"}),
	}).Create(&models.ErrorGroup{
		ProjectID:   projectID,
		Fingerprint: fingerprint,
		Title:       title,
		Status:      "critical",
		Count:       1,
		LastSeenAt:  now,
	}).Error

	if err != nil {
		return nil, err
	}

	config.DB.Where("project_id = ? AND fingerprint = ?", projectID, fingerprint).First(&group)
	return &group, nil
}
