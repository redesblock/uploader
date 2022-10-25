package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redesblock/uploader/core/model"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

// @Summary list vouchers
// @Schemes
// @Description pagination list vouchers
// @Tags Voucher
// @Accept json
// @Produce json
// @Param page_num query int false "page number"
// @Param page_size query int false "page size"
// @Success 200 {object} model.Voucher
// @Router /vouchers [get]
func VouchersHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		pageNum, pageSize := page(c)
		offset := (pageNum - 1) * pageSize

		var total int64
		var items []*model.Voucher
		if err := db.Model(&model.Voucher{}).Order("id desc").Count(&total).Offset(int(offset)).Limit(int(pageSize)).Find(&items).Error; err != nil {
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

// @Summary add voucher
// @Schemes
// @Description add voucher
// @Tags Voucher
// @Accept json
// @Produce json
// @Param voucher query string true "voucher"
// @Param host query string true "host"
// @Success 200 {object} model.Voucher
// @Router /add_voucher [get]
func AddVoucherHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		voucher := c.Query("voucher")
		if len(voucher) == 0 {
			log.Errorf("api %s error %v", c.Request.URL.Path, fmt.Errorf("invalid voucher (%s)", voucher))
			c.JSON(http.StatusBadRequest, fmt.Errorf("invalid voucher (%s)", voucher))
			return
		}
		host := c.Query("host")
		if len(host) == 0 {
			log.Errorf("api %s error %v", c.Request.URL.Path, fmt.Errorf("invalid host (%s)", host))
			c.JSON(http.StatusBadRequest, fmt.Errorf("invalid host (%s)", host))
			return
		}

		var item model.Voucher
		if res := db.Where("voucher = ?", voucher).Find(&item); res.Error != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, res.Error)
			c.JSON(http.StatusInternalServerError, res.Error)
			return
		} else if res.RowsAffected == 0 {
			item.Voucher = voucher
			item.Usable = true
		}
		item.Host = host

		if err := db.Save(&item).Error; err != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

// @Summary remove voucher
// @Schemes
// @Description remove voucher
// @Tags Voucher
// @Accept json
// @Produce json
// @Param voucher query string true "voucher"
// @Success 200 int 0 "affect rows"
// @Router /remove_voucher [get]
func RemoveVoucherHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		voucher := c.Query("voucher")
		if len(voucher) == 0 {
			log.Errorf("api %s error %v", c.Request.URL.Path, fmt.Errorf("invalid voucher (%s)", voucher))
			c.JSON(http.StatusBadRequest, fmt.Errorf("invalid voucher (%s)", voucher))
			return
		}

		res := db.Unscoped().Delete(&model.Voucher{}, "voucher = ?", voucher)
		if res.Error != nil {
			log.Errorf("api %s error %v", c.Request.URL.Path, res.Error)
			c.JSON(http.StatusInternalServerError, res.Error)
			return
		}

		c.JSON(http.StatusOK, res.RowsAffected)
	}
}
