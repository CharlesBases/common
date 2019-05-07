package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"common/log"
	"common/web/jwt"
	"common/web/weberror"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	_, err := jwt.GetUser(c)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, weberror.NewTokenInvalidError(err.Error()))
		return
	}

	file, _ := c.FormFile("upload")

	path := strings.Builder{}
	path.WriteString("/Users/sun/Downloads")
	path.WriteString("/picture")
	path.WriteString(fmt.Sprintf(`/%s`, time.Now().Format("20060102")))
	path.WriteString(fmt.Sprintf(`/%d%s`, time.Now().UnixNano(), file.Filename))

	err = c.SaveUploadedFile(file, path.String())
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusOK, weberror.NewFileUploadBaseError(err.Error()))
		return
	}
	c.JSON(http.StatusOK, weberror.Success(path))
}
