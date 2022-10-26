package server

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/redesblock/uploader/core/model"
	"github.com/redesblock/uploader/core/server/routers"
	"github.com/redesblock/uploader/core/syncer"
	"github.com/redesblock/uploader/core/util"
	"github.com/redesblock/uploader/docs"
	log "github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func Start(port string, db *gorm.DB, interval string, gateway string) error {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)
	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	})
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	docs.SwaggerInfo.BasePath = "/api"
	v1 := router.Group("/api")
	v1.GET("/vouchers", routers.VouchersHandler(db))
	v1.GET("/add_voucher", routers.AddVoucherHandler(db))
	v1.GET("/remove_voucher", routers.RemoveVoucherHandler(db))
	v1.GET("/watch_files", routers.WatchFilesHandler(db))
	v1.GET("/add_watch_file", routers.AddWatchFileHandler(db))
	v1.GET("/remove_watch_file", routers.RemoveWatchFileHandler(db))
	v1.GET("/upload_files", routers.UploadFilesHandler(db))
	v1.GET("/reference", routers.FileReferenceHandler(db, gateway))

	ignoreHidden := true
	syncer := syncer.New(ignoreHidden)
	_ = syncer

	scheduler := gocron.NewScheduler(time.UTC)
	if _, err := scheduler.Every(interval).Do(func() {
		duration, _ := time.ParseDuration(interval)
		logger := log.WithField("upload scheduler", interval)
		var vouchers []*model.Voucher
		if err := db.Model(&model.Voucher{}).Order("id desc").Where("usable = true").Find(&vouchers).Error; err != nil {
			logger.WithField("error", err).Errorf("load vouchers")
			return
		}

		voucherCnt := len(vouchers)
		if voucherCnt == 0 {
			logger.WithField("error", fmt.Errorf("no usable vouchers")).Errorf("load vouchers")
			return
		}

		var items []*model.WatchFile
		if err := db.Model(&model.WatchFile{}).Order("id desc").Find(&items).Error; err != nil {
			logger.WithField("error", err).Errorf("load watch files")
			return
		}

		voucherIndex := 0
		for _, item := range items {
			if err := filepath.Walk(item.Path, func(walkPath string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}
				isHidden, err := util.IsHiddenFile(walkPath)
				if err != nil {
					return fmt.Errorf("isHidden error %v", err)
				}
				if ignoreHidden && isHidden {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if info.IsDir() {
					return nil
				}
				// index name
				if strings.HasSuffix(info.Name(), item.IndexExt) {
					if time.Now().Sub(info.ModTime()) < duration {
						// ignore upload before duration
						return nil
					}
					// upload directory contains index name
					path := filepath.Dir(walkPath)
					relPath, err := filepath.Rel(item.Path, path)
					if err != nil {
						return err
					}

					var f model.UploadFile
					if res := db.Model(&model.UploadFile{}).Where("rel_path = ?", relPath).Find(&f); res.Error != nil {
						return fmt.Errorf("find upload file error %v", res.Error)
					} else if res.RowsAffected > 0 {
						if f.IndexName != info.Name() || info.ModTime().Sub(f.ModifyAt) < duration {
							// modify ignore before duration
							return nil
						}
						if f.Path != path {
							return fmt.Errorf("relPath %s already found in watch file %s", relPath, f.Path)
						}
					}

					hash := ""
					for i := 0; i < voucherCnt; i++ {
						voucherIndex += i
						voucher := vouchers[voucherIndex%voucherCnt]
						reference, err := syncer.Upload(voucher.Node, voucher.Voucher, path, item.IndexExt)
						if err != nil {
							logger.WithField("relPath", f.RelPath).WithField("index", f.IndexName).WithField("error", err).Warningf("synced to mop failed, try %d", i)
							continue
						} else {
							hash = reference
							break
						}
					}
					f.Path = path
					f.RelPath = relPath
					f.IndexName = info.Name()
					f.ModifyAt = info.ModTime()
					f.Usable = false
					if len(hash) != 0 {
						f.Hash = hash
						if err := db.Save(&f).Error; err != nil {
							return fmt.Errorf("save upload file error %v", err)
						}
						logger.WithField("relPath", f.RelPath).WithField("index", f.IndexName).WithField("reference", f.Hash).Info("synced to mop")
					} else {
						logger.WithField("relPath", f.RelPath).WithField("index", f.IndexName).Error("synced to mop failed")
					}
				}
				return nil
			}); err != nil {
				logger.WithField("path", item.Path).WithField("error", err).Errorf("walk watch file")
			}
		}
	}); err != nil {
		return err
	}

	if _, err := scheduler.Every("1m").Do(func() {
		logger := log.WithField("usable scheduler", interval)
		var vouchers []*model.Voucher
		if err := db.Model(&model.Voucher{}).Find(&vouchers).Error; err != nil {
			logger.WithField("error", err).Errorf("load vouchers")
			return
		}
		for _, voucher := range vouchers {
			usable, err := util.VoucherUsabe(voucher.Node, voucher.Voucher)
			if err != nil {
				logger.WithField("error", err).Errorf("find voucher usable")
			}
			if voucher.Usable != usable {
				voucher.Usable = usable
				if err := db.Save(&voucher).Error; err != nil {
					logger.WithField("error", err).Errorf("save voucher")
				}
			}
		}

		var items []*model.UploadFile
		if err := db.Model(&model.UploadFile{}).Where("usable = false").Find(&items).Error; err != nil {
			logger.WithField("error", err).Errorf("load upload files")
			return
		}

		for _, item := range items {
			usable, err := util.ReferenceUsabe(gateway, item.Hash)
			if err != nil {
				logger.WithField("error", err).Errorf("find reference usable")
			}
			if item.Usable != usable {
				item.Usable = usable
				if err := db.Save(&item).Error; err != nil {
					logger.WithField("error", err).Errorf("save upload file")
				}
			}
		}
	}); err != nil {
		return err
	}
	scheduler.StartAsync()

	log.WithField("port", port).Info("starting server")
	return router.Run(port)
}
