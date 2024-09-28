package hellofresh

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "import_hellofresh_index.tmpl", gin.H{})
}
