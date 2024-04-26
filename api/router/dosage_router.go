package router

import "github.com/gin-gonic/gin"

func initializeDosageRouter(r *gin.Engine) {
	dosage := r.Group("/dosage")
	{
		dosage.GET("/medications", func(c *gin.Context) {
			c.JSON(200, gin.H{"test": "medications"})
		})
		return
	}
}
