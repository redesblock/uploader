package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redesblock/uploader/core/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// @Summary list watch files
// @Schemes
// @Description pagination list watch files
// @Tags Watch File
// @Accept json
// @Produce json
// @Param page_num query int false "page number"
// @Param page_size query int false "page size"
// @Success 200 {object} model.WatchFile
// @Router /watch_files [get]
func WatchFilesHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageNum, pageSize := page(c)
		offset := (pageNum - 1) * pageSize

		var total int64
		var items []*model.WatchFile
		if err := db.Model(&model.WatchFile{}).Order("id desc").Count(&total).Offset(int(offset)).Limit(int(pageSize)).Find(&items).Error; err != nil {
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

// @Summary add watch file
// @Schemes
// @Description add watch file
// @Tags Watch File
// @Accept json
// @Produce json
// @Param path query string true "watch file path"
// @Param index query string true "index file or index file ext"
// @Success 200 {object} model.WatchFile
// @Router /add_watch_file [get]
func AddWatchFileHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		index := c.Query("index")
		if len(index) == 0 || strings.ContainsRune(index, os.PathSeparator) {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid index (%s)", index))
			return
		}

		path := c.Query("path")
		if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
			if err == nil {
				err = fmt.Errorf("only support folder")
			}
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid path %v", err))
			return
		}

		fullPath, err := filepath.Abs(path)
		if err != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusBadRequest, fmt.Errorf("invalid path %v", err))
			return
		}

		var item model.WatchFile
		if res := db.Where("path = ?", fullPath).Find(&item); res.Error != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, res.Error)
			c.JSON(http.StatusInternalServerError, res.Error)
			return
		} else if res.RowsAffected > 0 {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("path %s already wathced", fullPath))
			return
		}
		item.Path = fullPath
		item.IndexExt = index

		if err := db.Save(&item).Error; err != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// @Summary remove watch file
// @Schemes
// @Description remove watch file
// @Tags Watch File
// @Accept json
// @Produce json
// @Param path query string true "watch file path"
// @Success 200 int 0 "affect rows"
// @Router /remove_watch_file [get]
func RemoveWatchFileHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		path := c.Query("path")
		fi, err := os.Stat(path)
		if err != nil || !fi.IsDir() {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid path (%s)", path))
			return
		}

		fullPath, err := filepath.Abs(path)
		if err != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		res := db.Unscoped().Delete(&model.WatchFile{}, "path = ?", fullPath)
		if res.Error != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, res.Error)
			c.JSON(http.StatusInternalServerError, res.Error)
			return
		}

		c.JSON(http.StatusOK, res.RowsAffected)
	}
}
