/*
 * @Description:
 * @Author: Allen
 * @Date: 2021-05-22 10:55:27
 * @LastEditors: Allen
 * @LastEditTime: 2021-05-22 10:55:27
 */
package main

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"payment/pkg/common/util"
	"payment/pkg/router/grpc"
	"payment/pkg/router/http"
	prof "payment/pkg/router/pprof"
	"payment/pkg/version"

	"payment/pkg/infra"

	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
)

var printVersion bool

type Payment struct {
	HttpServer *http.HttpServer `yaml:"http_server"`
	GrpcServer *grpc.GrpcServer `yaml:"grpc_server"`
	MysqlDB *infra.MysqlDB `yaml:"mysql"`
	RedisConfig *infra.RedisConfig `yaml:"redis"`
	Logging *infra.Logging `yaml:"logging"`
	Pprof *prof.Prof `yaml:"pprof"`
	Worker *util.Worker `yaml:"worker"`
	Nsq *infra.NsqAgent `yaml:"nsq"`
}

func init() {
    flag.BoolVar(&printVersion, "version", false, "print program build version")
    flag.Parse()
}

func main() {
	//输出编译信息
	version.PrintCLIVersion()
	if printVersion {
		return
	}
	payment, err := initialization()
	if err != nil {
		println("failed to Genesis, err:%s", err.Error())
	}
	payment.Logging.Setup(version.BuildUser, version.Version)
	//pprof 性能分析器
	payment.Pprof.StartPerf()
	infra.SetupDB(payment.MysqlDB)
	payment.RedisConfig.SetupRedis()
	//start nsq
	payment.Nsq.InitProducer()

	ctx, cancel := context.WithCancel(context.Background())
	group, errCtx := errgroup.WithContext(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("logAgent init producer failed")
	}
	group.Go(func() error {
		if err := payment.GrpcServer.CreateServer(); err != nil {
			defer cancel()
			return err
		}
		return nil
	})
	group.Go(func() error {
		if err := payment.HttpServer.CreateServer(); err != nil {
			defer cancel()
			return err
		}
		return nil
	})
	group.Go(func() error {
		defer func() {
			payment.GrpcServer.StopServer()
			cancel()
		}()
		for {
			select {
			case <-errCtx.Done():
				log.Println("recv ctx:", errCtx.Err().Error())
				return errCtx.Err()
			}
		}
	})
	if err := group.Wait(); err != nil {
		log.Info("----------- All goroutine done. err:", err.Error())
	} else {
		log.Info("----------- All goroutine done.")
	}
}


func initialization() (*Payment, error) {
	//读取配置文件
	yamlFile, err := ioutil.ReadFile("config/env/" + os.Getenv("ENV") + ".yaml")
	if err != nil {
		println("failed to readFile, err:%s", err.Error())
		return nil, err
	}
	payment := &Payment{
		HttpServer: http.NewHttpServer(),
		GrpcServer: grpc.NewGrpcServer(),
		MysqlDB: new(infra.MysqlDB),
		RedisConfig: new(infra.RedisConfig),
		Logging: new(infra.Logging),
		Pprof: new(prof.Prof),
		Worker: util.WorkerInstance(),
		Nsq: infra.CreateNsqInstance(),
	}
	err = yaml.Unmarshal(yamlFile, payment)
	if err != nil {
		println("failed to yaml.Unmarshal, err:%s", err.Error())
		return nil, err
	}

	return payment, nil
}
