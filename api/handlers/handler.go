package handlers

import (
	"database/sql"
	"net/http"
	"skillsapi/api/models"
	"skillsapi/db"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func GetSkill(c *gin.Context) {
	key := c.Param("key")

	var skill models.Skill
	err := db.DB.QueryRow("SELECT key, name, description, logo, tags FROM skills WHERE key = $1", key).Scan(
		&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Skill not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func GetAllSkills(c *gin.Context) {
	rows, err := db.DB.Query("SELECT key, name, description, logo, tags FROM skills")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}
	defer rows.Close()

	var skills []models.Skill
	for rows.Next() {
		var skill models.Skill
		if err := rows.Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
			return
		}
		skills = append(skills, skill)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skills,
	})
}

func CreateSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var existingKey string
	err := db.DB.QueryRow("SELECT key FROM skills WHERE key = $1", skill.Key).Scan(&existingKey)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if existingKey != "" {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Skill already exists"})
		return
	}

	_, err = db.DB.Exec("INSERT INTO skills (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5)",
		skill.Key, skill.Name, skill.Description, skill.Logo, pq.Array(skill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func UpdateSkill(c *gin.Context) {
	key := c.Param("key")
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := db.DB.Exec("UPDATE skills SET name = $1, description = $2, logo = $3, tags = $4 WHERE key = $5",
		skill.Name, skill.Description, skill.Logo, pq.Array(skill.Tags), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to update skill"})
		return
	}

	skill.Key = key

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func UpdateSkillName(c *gin.Context) {
	key := c.Param("key")
	var payload struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := db.DB.Exec("UPDATE skills SET name = $1 WHERE key = $2", payload.Name, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to update skill name"})
		return
	}

	var skill models.Skill
	err = db.DB.QueryRow("SELECT key, name, description, logo, tags FROM skills WHERE key = $1", key).Scan(
		&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func UpdateSkillDescription(c *gin.Context) {
	key := c.Param("key")
	var payload struct {
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := db.DB.Exec("UPDATE skills SET description = $1 WHERE key = $2", payload.Description, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to update skill description"})
		return
	}

	var skill models.Skill
	err = db.DB.QueryRow("SELECT key, name, description, logo, tags FROM skills WHERE key = $1", key).Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func UpdateSkillLogo(c *gin.Context) {
	key := c.Param("key")
	var payload struct {
		Logo string `json:"logo"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := db.DB.Exec("UPDATE skills SET logo = $1 WHERE key = $2", payload.Logo, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to update skill logo"})
		return
	}

	var skill models.Skill
	err = db.DB.QueryRow("SELECT key, name, description, logo, tags FROM skills WHERE key = $1", key).Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func UpdateSkillTags(c *gin.Context) {
	key := c.Param("key")
	var payload struct {
		Tags []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	result, err := db.DB.Exec("UPDATE skills SET tags = $1 WHERE key = $2", pq.Array(payload.Tags), key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to update skill tags"})
		return
	}

	var skill models.Skill
	err = db.DB.QueryRow("SELECT key, name, description, logo, tags FROM skills WHERE key = $1", key).Scan(&skill.Key, &skill.Name, &skill.Description, &skill.Logo, pq.Array(&skill.Tags))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   skill,
	})
}

func DeleteSkill(c *gin.Context) {
	key := c.Param("key")

	result, err := db.DB.Exec("DELETE FROM skills WHERE key = $1", key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to delete skill"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Unable to delete skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Skill deleted",
	})
}
