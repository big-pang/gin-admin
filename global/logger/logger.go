package logger

import (
	"fmt"
	"gin-admin/global"
	"gin-admin/utils"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var timer *time.Timer
var sigs chan os.Signal

func Sugar() *zap.SugaredLogger {
	return Logger.Sugar()
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// getLumberjackLogger returns a Lumberjack logger with the current date in the filename
func getLumberjackLogger(serviceName string) *lumberjack.Logger {

	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", serviceName), // 日志文件路径
		MaxSize:    100,                                // 单个日志文件最大尺寸（以 MB 为单位）
		MaxBackups: 7,                                  // 保留的旧日志文件的最大数量
		MaxAge:     7,                                  // 日志文件的最大保存天数
		Compress:   true,                               // 是否压缩旧的日志文件
	}
}

func Sync() {
	Logger.Sync()
}

type LoggerInterface interface {
	Write(p []byte) (n int, err error)
}

func HandleSignals(serviceName string) {
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1)

	go func() {
		for {
			sig := <-sigs
			switch sig {
			case syscall.SIGUSR1:
				Logger.Sync() // flushes buffer, if any
				fmt.Println("Received SIGUSR1, flushing logs")
				LogInit(serviceName) // reinitialize logger
			}
		}
	}()
}

func LogInit(serviceName string) {
	fmt.Println(utils.IsRunningUnderSystemd())
	// 判断是否为systemctl启动
	if utils.IsRunningUnderSystemd() {
		pidFile := fmt.Sprintf("/var/run/%s.pid", serviceName)
		err := utils.CreatePIDFile(pidFile)
		if err != nil {
			fmt.Println("创建PID文件失败", err)
			os.Exit(1)
			return
		}
		HandleSignals(serviceName)
		systemdLoggerInit(serviceName)
	} else {
		logInit(serviceName)
	}

}

func systemdLoggerInit(serviceName string) {
	lumberjackLogger := getLumberjackLogger("/var/log/dingtalk/" + serviceName)

	zapInit(lumberjackLogger)
}

func logInit(serviceName string) {

	err := utils.CheckAndCreateDir("logs")
	if err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}
	currentDate := time.Now().Format("2006-01-02")
	// 创建一个 Lumberjack logger
	lumberjackLogger := getLumberjackLogger("./logs/" + serviceName + "-" + currentDate)
	// defer LOG.Sync() // 确保日志在程序退出前被写入
	// 定时重启日志
	timer = time.AfterFunc(utils.TimeUntilTomorrowMidnight()+5*time.Minute, func() {
		logInit(serviceName)
	})
	zapInit(lumberjackLogger)

}
func zapInit(lumberjackLogger LoggerInterface) {
	// 配置 Zap 使用 Lumberjack 作为日志输出
	writeSyncer := zapcore.AddSync(lumberjackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	Logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 使用 JSON 编码器
		writeSyncer,                           // 设置日志写入目标
		zapcore.InfoLevel,                     // 设置日志级别
	), zap.AddCaller())
	global.LOG = Logger
}
