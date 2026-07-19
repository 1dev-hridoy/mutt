package handler

import (
	"time"

	"github.com/dishan1223/mutt/internal/config"
	"github.com/dishan1223/mutt/models"
	"github.com/gofiber/fiber/v3"
)

func StatsHandler(c fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var projectIDs []uint
	config.DB.Model(&models.Project{}).Where("user_id = ?", userID).Pluck("id", &projectIDs)

	stats := models.StatsResponse{
		ByStatus: map[string]int{"critical": 0, "resolved": 0, "recovered": 0},
	}

	stats.TotalProjects = len(projectIDs)

	if len(projectIDs) == 0 {
		return c.JSON(stats)
	}

	config.DB.Model(&models.ErrorGroup{}).Where("project_id IN ?", projectIDs).Select("COUNT(*)").Scan(&stats.TotalErrorGroups)

	config.DB.Model(&models.Error{}).Where("project_id IN ?", projectIDs).Select("COUNT(*)").Scan(&stats.TotalErrors)

	var statusCounts []struct {
		Status string
		Count  int
	}
	config.DB.Model(&models.ErrorGroup{}).Where("project_id IN ?", projectIDs).Select("status, COUNT(*) as count").Group("status").Scan(&statusCounts)
	for _, s := range statusCounts {
		stats.ByStatus[s.Status] = s.Count
	}

	since := time.Now().Add(-24 * time.Hour)
	config.DB.Model(&models.Error{}).Where("project_id IN ? AND occurred_at >= ?", projectIDs, since).Select("COUNT(*)").Scan(&stats.ErrorsLast24h)

	return c.JSON(stats)
}
