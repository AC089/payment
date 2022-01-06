/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-28 11:02:42
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-28 11:02:47
 */
package infra

import (
	"fmt"
	"os"
	"path/filepath"
	"payment/pkg/common/util"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

type Logging struct {
	ServerName string `yaml:"serverName"`
	Level string `yaml:"level"`
	MaxRemain time.Duration `yaml:"maxRemain"`
	RotationTime time.Duration `yaml:"rotationTime"`
	Path string `yaml:"path"`
}

func GetLogLevel(s string) log.Level {
	switch s {
	case "Debug":
		return log.DebugLevel
	case "Info":
		return log.InfoLevel
	case "Warn":
		return log.WarnLevel
	case "Error":
		return log.ErrorLevel
	case "Fatal":
		return log.FatalLevel
	case "Panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func (p *Logging) newLfsHook(logName string, errLogName string) (log.Hook, log.Hook) {
	//取容器id前12
	// out, _ := util.Exec_shell("cat /proc/self/cgroup | grep -o -e \"docker/.*\"| head -n 1 |sed \"s/docker\\/\\(.*\\)/\\1/\"|cut -c1-10")
	writer, err := rotatelogs.New(
		// logName+"."+out+".%Y%m%d%H",
		logName+".%Y%m%d%H",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*p.RotationTime),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(time.Hour*p.MaxRemain),
		//rotatelogs.WithRotationCount(maxRemain),
	)
	writer2, err := rotatelogs.New(
		// errLogName+"."+out+".%Y%m%d%H",
		errLogName+".%Y%m%d%H",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(errLogName),
		rotatelogs.WithRotationTime(time.Hour*p.RotationTime),
		rotatelogs.WithMaxAge(time.Hour*p.MaxRemain),
	)
	if err != nil {
		log.Errorf("config local file system for logger error: %v", err)
	}

	log.SetLevel(GetLogLevel(p.Level))

	lfsHook1 := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer,
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
		//}, &log.TextFormatter{DisableColors: true})
	}, &log.JSONFormatter{})

	lfsHook2 := lfshook.NewHook(lfshook.WriterMap{
		log.ErrorLevel: writer2,
		log.FatalLevel: writer2,
		log.PanicLevel: writer2,
		//}, &log.TextFormatter{DisableColors: true})
	}, &log.JSONFormatter{})

	return lfsHook1, lfsHook2
}

func (p *Logging) Setup(buildUser, version string) {
	fn := fmt.Sprintf("%s_%s_%s", p.ServerName, buildUser, version)
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	//	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, DisableColors: false})
	log.SetReportCaller(true)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.InfoLevel)

	if !util.Exists(p.Path) {
		if !util.Mkdir(p.Path) {
			log.WithFields(log.Fields{
				"err": "can not mkdir",
			}).Fatal()
		}
	}

	filename := filepath.Join(p.Path, fn)
	// 添加错误日志目录
	errPath := filepath.Join(p.Path, "error")
	if !util.Exists(errPath) {
		if !util.Mkdir(errPath) {
			log.WithFields(log.Fields{
				"err": "can not mkdir",
			}).Fatal()
		}
	}
	errfilename := filepath.Join(errPath, fn)
	lfsHook1, lfsHook2 := p.newLfsHook(filename, errfilename)
	log.AddHook(lfsHook1)
	log.AddHook(lfsHook2)
}