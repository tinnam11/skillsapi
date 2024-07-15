package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"skillsapi/api/models"
	"skillsapi/db"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/api/v1/skills/:key", GetSkill)
	r.GET("/api/v1/skills", GetAllSkills)
	r.POST("/api/v1/skills", CreateSkill)
	r.PUT("/api/v1/skills/:key", UpdateSkill)
	r.PATCH("/api/v1/skills/:key/actions/name", UpdateSkillName)
	r.PATCH("/api/v1/skills/:key/actions/description", UpdateSkillDescription)
	r.PATCH("/api/v1/skills/:key/actions/logo", UpdateSkillLogo)
	r.PATCH("/api/v1/skills/:key/actions/tags", UpdateSkillTags)
	r.DELETE("/api/v1/skills/:key", DeleteSkill)
	return r
}

func insertTestSkill(t *testing.T, key, name, description, logo string, tags []string) {
	deleteTestSkill(t, key) // Ensure the skill does not already exist
	_, err := db.DB.Exec("INSERT INTO skills (key, name, description, logo, tags) VALUES ($1, $2, $3, $4, $5)",
		key, name, description, logo, pq.Array(tags))
	if err != nil {
		t.Fatalf("Error inserting test skill: %v", err)
	}
}

func deleteTestSkill(t *testing.T, key string) {
	_, err := db.DB.Exec("DELETE FROM skills WHERE key = $1", key)
	if err != nil && err != sql.ErrNoRows {
		t.Fatalf("Error deleting test skill: %v", err)
	}
}

func clearSkillsTable(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM skills")
	if err != nil {
		t.Fatalf("Error clearing skills table: %v", err)
	}
}

func TestMain(m *testing.M) {
	db.Init()
	defer db.Close()
	clearSkillsTable(nil)
	code := m.Run()
	os.Exit(code)
}

func TestGetSkill(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})

	req, _ := http.NewRequest("GET", "/api/v1/skills/python", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Python", responseBody["data"].(map[string]interface{})["name"])

	deleteTestSkill(t, "python")
}

func TestGetAllSkills(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	insertTestSkill(t, "golang", "Go", "Go is a statically typed, compiled programming language designed at Google.", "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg", []string{"programming language", "system"})

	req, _ := http.NewRequest("GET", "/api/v1/skills", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])

	data, ok := responseBody["data"].([]interface{})
	assert.True(t, ok)
	assert.Len(t, data, 2)

	deleteTestSkill(t, "python")
	deleteTestSkill(t, "golang")
}

func TestCreateSkill(t *testing.T) {
	router := setupRouter()
	skill := models.Skill{
		Key:         "golang",
		Name:        "Go",
		Description: "Go is a statically typed, compiled programming language designed at Google.",
		Logo:        "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
		Tags:        []string{"programming language", "system"},
	}
	body, _ := json.Marshal(skill)

	req, _ := http.NewRequest("POST", "/api/v1/skills", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Go", responseBody["data"].(map[string]interface{})["name"])

	deleteTestSkill(t, "golang")
}

func TestUpdateSkill(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	updatedSkill := models.Skill{
		Name:        "Python 3",
		Description: "Python 3 is the latest version of Python programming language.",
		Logo:        "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg",
		Tags:        []string{"data"},
	}
	body, _ := json.Marshal(updatedSkill)

	req, _ := http.NewRequest("PUT", "/api/v1/skills/python", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Python 3", responseBody["data"].(map[string]interface{})["name"])

	deleteTestSkill(t, "python")
}

func TestUpdateSkillName(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	payload := map[string]string{"name": "Python 3"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "/api/v1/skills/python/actions/name", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Python 3", responseBody["data"].(map[string]interface{})["name"])

	deleteTestSkill(t, "python")
}

func TestUpdateSkillDescription(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	payload := map[string]string{"description": "Python 3 is the latest version of Python programming language."}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "/api/v1/skills/python/actions/description", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Python 3 is the latest version of Python programming language.", responseBody["data"].(map[string]interface{})["description"])

	deleteTestSkill(t, "python")
}

func TestUpdateSkillLogo(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	payload := map[string]string{"logo": "https://newlogo.com/logo.svg"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "/api/v1/skills/python/actions/logo", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "https://newlogo.com/logo.svg", responseBody["data"].(map[string]interface{})["logo"])

	deleteTestSkill(t, "python")
}

func TestUpdateSkillTags(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})
	payload := map[string][]string{"tags": {"programming language", "data"}}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PATCH", "/api/v1/skills/python/actions/tags", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, []interface{}{"programming language", "data"}, responseBody["data"].(map[string]interface{})["tags"])

	deleteTestSkill(t, "python")
}

func TestDeleteSkill(t *testing.T) {
	router := setupRouter()
	insertTestSkill(t, "python", "Python", "Python is an interpreted, high-level, general-purpose programming language.", "https://upload.wikimedia.org/wikipedia/commons/c/c3/Python-logo-notext.svg", []string{"programming language", "scripting"})

	req, _ := http.NewRequest("DELETE", "/api/v1/skills/python", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Equal(t, "success", responseBody["status"])
	assert.Equal(t, "Skill deleted", responseBody["message"])
}
