package main

import (
	"cctv-image-receiver/conf"
	"cctv-image-receiver/service"
	"cctv-image-receiver/util"
	"fmt"
	"os"

	"github.com/robfig/cron/v3"
)

func main() {

	conf := conf.CurrentConfig

	if conf.Os.IsWindows {
		os.Setenv("ZONEINFO", conf.Os.TimeFilePath)
	}

	c := cron.New()

	fmt.Printf("執行抓取%s <%s>\n", conf.Param.ProjectName, util.GetTimeNow())

	cctvCountMap := make(map[string]int, 0)

	for _, cctv := range conf.Cctvs {

		service.SetCCTVCountMap(cctvCountMap, cctv.No)

		c.AddJob(conf.Param.CronSpec, service.ImageJob{
			CctvCountMap: cctvCountMap, ImageOutPath: conf.Param.ImageOutPath, Cctv: *cctv})
	}

	c.Start()

	defer c.Stop()
	select {}
}
