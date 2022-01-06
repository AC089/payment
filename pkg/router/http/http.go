/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-24 14:38:55
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-24 14:38:57
 */
package http

import (
	"fmt"
	"net/http"
	"payment/pkg/facade"
	"payment/pkg/facade/impl"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NewHttpServer() *HttpServer {
	return &HttpServer{
		tradeDetailFacade: impl.NewTradeDetailFacadeImpl(),
	}
}

type HttpServer struct {
	Port uint32 `yaml:"port"`
	PemPath string `yaml:"pem_path"`
	KeyPath string `yaml:"key_path"`
	ReadTimeout time.Duration `yaml:"readTimeout"`
	WriteTimeout time.Duration `yaml:"writeTimeout"`
	RunMode string `yaml:"runMode"`
	srv *http.Server

	tradeDetailFacade facade.TradeDetailFacade
}

func (h HttpServer) CreateServer() error {
	h.srv = &http.Server{
		Addr:           fmt.Sprintf(":%d", h.Port),
		Handler:        h.handle(),
		ReadTimeout:    h.ReadTimeout * time.Second,
		WriteTimeout:   h.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := h.srv.ListenAndServeTLS(h.PemPath, h.KeyPath); err != nil && err != http.ErrServerClosed {
		log.Printf("webServer shutdown: %s\n", err)
		return err
	} else if err == http.ErrServerClosed{
		log.Printf("webServer shutdown.")
		return err
	}
	return nil
}

func (h HttpServer) handle() *gin.Engine {
	r := gin.Default()
	// 跨域处理
	//r.Use(cors.Default()) // 跨域，允许所有源
	gin.SetMode(h.RunMode)
	corsMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			//c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost")
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

			if c.Request.Method == "OPTIONS" {
				fmt.Println("OPTIONS")
				c.AbortWithStatus(200)
			} else {
				c.Next()
			}
		}
	}
	r.Use(corsMiddleware())

	apiNotify := r.Group("/notify")
	apiNotify.POST("/alipay", h.NotifyAlipay)    
	apiNotify.POST("/wechat/:confId", h.NotifyWechat)
	log.Infof("Create http server success!port is %d", h.Port)
	return r
}
