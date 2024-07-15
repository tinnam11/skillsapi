package routes

import (
	"skillsapi/api/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		skills := api.Group("/skills")
		{
			skills.GET("/", handlers.GetAllSkills)
			skills.GET("/:key", handlers.GetSkill)
			skills.POST("/", handlers.CreateSkill)
			skills.PUT("/:key", handlers.UpdateSkill)
			skills.PATCH("/:key/actions/name", handlers.UpdateSkillName)
			skills.PATCH("/:key/actions/description", handlers.UpdateSkillDescription)
			skills.PATCH("/:key/actions/logo", handlers.UpdateSkillLogo)
			skills.PATCH("/:key/actions/tags", handlers.UpdateSkillTags)
			skills.DELETE("/:key", handlers.DeleteSkill)
		}
	}
}
