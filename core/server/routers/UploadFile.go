package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/redesblock/uploader/core/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

// @Summary list upload files
// @Schemes
// @Description pagination list upload files
// @Tags Upload File
// @Accept json
// @Produce json
// @Param page_num query int false "page number"
// @Param page_size query int false "page size"
// @Success 200 {object} model.UploadFile
// @Router /upload_files [get]
func UploadFilesHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageNum, pageSize := page(c)
		offset := (pageNum - 1) * pageSize

		var total int64
		var items []*model.UploadFile
		if err := db.Model(&model.UploadFile{}).Order("id desc").Count(&total).Offset(int(offset)).Limit(int(pageSize)).Find(&items).Error; err != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, &List{
			Total: total,
			PageTotal: func() int64 {
				pageTotal := total / pageSize
				if total%pageSize != 0 {
					pageTotal++
				}
				return pageTotal
			}(),
			Items: items,
		})
	}
}

// @Summary upload file reference
// @Schemes
// @Description upload file reference
// @Tags Upload File
// @Accept json
// @Produce json
// @Param path query string true "file path"
// @Param usable query bool false "usable"
// @Success 200 string 74e95e2785817fbad5c2b29f62a812009d5b43850856757ce35627283df7a817
// @Router /reference [get]
func FileReferenceHandler(db *gorm.DB, gateway string) func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Query("path")
		items := strings.Split(path, "/")
		cnt := len(items)
		for i := cnt; i > 0; i-- {
			tPath := strings.Join(items[:i], "/")
			var item model.UploadFile
			tx := db.Where("rel_path = ?", tPath)
			if strings.ToLower(c.Query("usable")) == "true" {
				tx = tx.Where("usable = true")
			}
			if res := tx.Find(&item); res.Error != nil {
				log.Errorf("api %s error %v", c.Request.URL.Path, res.Error)
				c.JSON(http.StatusInternalServerError, res.Error)
				return
			} else if res.RowsAffected > 0 {
				refs := []string{item.Hash}
				if i != cnt {
					refs = append(refs, items[i:]...)
					c.JSON(http.StatusOK, gateway+"/mop/"+strings.Join(refs, "/"))
				} else {
					c.JSON(http.StatusOK, gateway+"/mop/"+strings.Join(refs, "/")+"/")
				}
				return
			}
		}
		c.JSON(http.StatusOK, "")
	}
}
