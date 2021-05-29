/**
 * @Description exchange-list-monitor
 * @Date 2021.03.19
 **/
package main

import (
	"flag"
	"fmt"
	"github.com/qcozof/exchange-list-monitor/utils"
	"log"
	"time"

	"github.com/qcozof/exchange-list-monitor/app"
	"github.com/qcozof/exchange-list-monitor/global"
	"github.com/qcozof/exchange-list-monitor/initialize"
)

var sysConfig = global.SERVER_CONFIG.SystemConfig
var mxConfig = global.SERVER_CONFIG.MxcConfig
var gsConfig = global.SERVER_CONFIG.GrayScaleConfig

var commonUtils  utils.CommonUtils

func main()  {
/*	initialize.RedisInitialize()
	app.GrayScaleListMonitor()
	return*/

	arg1 := flag.String("arg1","","参数:list、ticker")
	flag.Parse()

	switch *arg1 {
	case "ticker":
		initialize.RedisInitialize()

		mt := time.NewTicker(time.Second * time.Duration(mxConfig.IntervalExecSeconds))
		myFunc := func() {
			beginTime :=  parseTime(mxConfig.BeginTime)
			endTime :=  parseTime(mxConfig.EndTime)
			timeNow := time.Now().UTC().Add(time.Hour*8)

			if timeNow.After(beginTime) && timeNow.Before(endTime){
				app.MxcTickerMonitor()
				return
			}
			fmt.Printf("------------time not reach, now: %s   time region:%s ~ %s-----------\n", timeNow,beginTime,endTime)
		}
		run(mt, myFunc )

		break

	case "list":
		fmt.Printf("------------ExchangeListMonitor %s-----------\n", commonUtils.NowStr())

		t := time.NewTicker(time.Second * time.Duration(sysConfig.IntervalExecSeconds))
		run(t, app.ExchangeListMonitor)
		break

	case "grayscale":
		initialize.RedisInitialize()
		fmt.Printf("------------GrayScaleListMonitor %s-----------\n", commonUtils.NowStr())

		t := time.NewTicker(time.Second * time.Duration(gsConfig.IntervalExecSeconds))
		run(t, app.GrayScaleListMonitor)
		break

	default:
		fmt.Println("--help 查看帮助 \n exchange-list-monitor[.exe] -arg1=ticker、list或grayscale")
		break
	}
}

func run( t  *time.Ticker, myFunc func() )  {
	for{
		select{
		case <-t.C:
			myFunc()
			break
		}
	}
}

func parseTime(timeStr string) time.Time {
	today := commonUtils.NowStr()[0:11]
	_,time,err :=  commonUtils.Str2TimeAndStamp(today+timeStr)
	if err !=nil{
		log.Fatal(err)
	}
	return time
}
