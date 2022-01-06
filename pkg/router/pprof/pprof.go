/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-31 16:58:21
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-31 16:58:21
 */
package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/pkg/errors"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp" // 暴露runtime信息
)

var (
	_perfOnce sync.Once
)

type Prof struct {
	Port int `yaml:"port"`
}

func (p *Prof) StartPerf() {
	_perfOnce.Do(func() {
		engine := gin.New()
		prefixRouter := engine.Group("/debug/pprof")
		{
			prefixRouter.GET("/", p.pprofHandler(pprof.Index))
			prefixRouter.GET("/cmdline", p.pprofHandler(pprof.Cmdline))
			prefixRouter.GET("/profile", p.pprofHandler(pprof.Profile))
			prefixRouter.POST("/symbol", p.pprofHandler(pprof.Symbol))
			prefixRouter.GET("/symbol", p.pprofHandler(pprof.Symbol))
			prefixRouter.GET("/trace", p.pprofHandler(pprof.Trace))
			prefixRouter.GET("/allocs", p.pprofHandler(pprof.Handler("allocs").ServeHTTP))
			prefixRouter.GET("/block", p.pprofHandler(pprof.Handler("block").ServeHTTP))
			prefixRouter.GET("/goroutine", p.pprofHandler(pprof.Handler("goroutine").ServeHTTP))
			prefixRouter.GET("/heap", p.pprofHandler(pprof.Handler("heap").ServeHTTP))
			prefixRouter.GET("/mutex", p.pprofHandler(pprof.Handler("mutex").ServeHTTP))
			prefixRouter.GET("/threadcreate", p.pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
		}
		gin.SetMode(gin.ReleaseMode)
		// prometheus
		engine.GET("/metrics", p.monitor())
		addr := fmt.Sprintf(":%d", p.Port)
		s := &http.Server{
			Addr:           addr,
			Handler:        engine,
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		go func() {
			if err := s.ListenAndServe(); err != nil {
				panic(errors.Errorf("perf: listen %s: error(%v)", addr, err))
			}
		}()
	})
}

func (p *Prof) pprofHandler(h func(http.ResponseWriter, *http.Request)) gin.HandlerFunc {
	return func(c *gin.Context) {
		h(c.Writer, c.Request)
	}
}

func (p *Prof) monitor() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := promhttp.Handler()
		h.ServeHTTP(c.Writer, c.Request)
	}
}
