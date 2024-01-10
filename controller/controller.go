package controller

import (
	"BBoltViewer/db"
	"BBoltViewer/g"
	"BBoltViewer/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var resp utils.Response

func Connect() gin.HandlerFunc {
	var param struct {
		DBPath string `json:"db_path" binding:"required"`
	}
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if !utils.IsPathExist(param.DBPath) {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db path not exist"))
			return
		}
		var (
			visiter *db.Visiter
			err     error
			ok      bool
		)
		if visiter, ok = g.CacheDBConn[param.DBPath]; !ok {
			visiter, err = db.Open(param.DBPath, &db.Options{Timeout: time.Second * 1})
			if err != nil {
				c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
				return
			}
			g.CacheDBConn[param.DBPath] = visiter
		}
		if buckets, err := visiter.Buckets(); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
		} else {
			c.JSON(http.StatusOK, resp.Success(buckets))
		}
	}
}

func Close() gin.HandlerFunc {
	var param struct {
		DBPath string `json:"db_path" binding:"required"`
	}
	return func(c *gin.Context) {
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if err := g.CacheDBConn[param.DBPath].Close(); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		delete(g.CacheDBConn, param.DBPath)
		c.JSON(http.StatusOK, resp.Success(param.DBPath))
	}
}

func DataFromBucket() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			DBPath string `json:"db_path" binding:"required"`
			Bucket string `json:"bucket" binding:"required"`
		}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if !utils.IsPathExist(param.DBPath) {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db path not exist"))
			return
		}
		var (
			visiter *db.Visiter
			ok      bool
		)
		if visiter, ok = g.CacheDBConn[param.DBPath]; !ok {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db not connect"))
			return
		}
		if data, err := visiter.DataFromBucket(param.Bucket); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
		} else {
			c.JSON(http.StatusOK, resp.Success(data))
		}
	}
}

func Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			DBPath string `json:"db_path" binding:"required"`
			Bucket string `json:"bucket" binding:"required"`
			Key    string `json:"key" binding:"required"`
			Value  string `json:"value" binding:"required"`
		}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if !utils.IsPathExist(param.DBPath) {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db path not exist"))
			return
		}
		var (
			visiter *db.Visiter
			ok      bool
		)
		if visiter, ok = g.CacheDBConn[param.DBPath]; !ok {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db not connect"))
			return
		}
		if err := visiter.Set(param.Bucket, param.Key, []byte(param.Value)); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		c.JSON(http.StatusOK, resp.Success(param.Key))
	}
}

func Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			DBPath string `json:"db_path" binding:"required"`
			Bucket string `json:"bucket" binding:"required"`
			Key    string `json:"key" binding:"required"`
		}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if !utils.IsPathExist(param.DBPath) {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db path not exist"))
			return
		}
		var (
			visiter *db.Visiter
			ok      bool
		)
		if visiter, ok = g.CacheDBConn[param.DBPath]; !ok {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db not connect"))
			return
		}
		if data, err := visiter.Get(param.Bucket, param.Key); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
		} else {
			c.JSON(http.StatusOK, resp.Success(string(data)))
		}
	}
}

func Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var param struct {
			DBPath string `json:"db_path" binding:"required"`
			Bucket string `json:"bucket" binding:"required"`
			Key    string `json:"key" binding:"required"`
		}
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		if !utils.IsPathExist(param.DBPath) {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db path not exist"))
			return
		}
		var (
			visiter *db.Visiter
			ok      bool
		)
		if visiter, ok = g.CacheDBConn[param.DBPath]; !ok {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc("db not connect"))
			return
		}
		if err := visiter.Delete(param.Bucket, param.Key); err != nil {
			c.JSON(http.StatusBadRequest, resp.Failure().WithDesc(err.Error()))
			return
		}
		c.JSON(http.StatusOK, resp.Success(param.Key))
	}
}
