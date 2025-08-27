package log

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
	"path"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	hertzzap "github.com/hertz-contrib/logger/zap"
	"go.uber.org/zap/zapcore"
)

// logWriter 实现 io.Writer 接口，将标准 log 输出重定向到 hlog
type logWriter struct{}

func (w *logWriter) Write(p []byte) (n int, err error) {
	hlog.Info(string(p))
	return len(p), nil
}

// getCurrentLogFile 获取当前日志文件名
func getCurrentLogFile(logPath, prefix string) string {
	return path.Join(logPath, prefix+".log")
}

// getLogFileByDate 根据日期获取日志文件名
func getLogFileByDate(logPath, prefix string, date time.Time) string {
	return path.Join(logPath, prefix+"_"+date.Format("2006-01-02")+".log")
}

// rotateLogFile 如果需要，将当前日志文件按创建时间重命名
func rotateLogFile(logPath, prefix string) error {
	currentFile := getCurrentLogFile(logPath, prefix)
	if _, err := os.Stat(currentFile); os.IsNotExist(err) {
		return nil // 如果当前日志文件不存在，不需要轮转
	}

	// 获取当前日志文件的创建时间
	fileInfo, err := os.Stat(currentFile)
	if err != nil {
		return err
	}

	// 使用文件的创建时间来命名
	datedFile := getLogFileByDate(logPath, prefix, fileInfo.ModTime())

	// 如果目标文件已存在，添加一个后缀
	if _, err := os.Stat(datedFile); err == nil {
		for i := 1; ; i++ {
			newName := datedFile[:len(datedFile)-4] + fmt.Sprintf("_%d.log", i)
			if _, err := os.Stat(newName); os.IsNotExist(err) {
				datedFile = newName
				break
			}
		}
	}

	// 重命名旧文件
	if err := os.Rename(currentFile, datedFile); err != nil {
		return fmt.Errorf("failed to rename log file: %v", err)
	}

	return nil
}

func Config(cof *LogConf, env constant.EnvMode) {
	if env == constant.Development {
		//Set log output location
		hlog.SetOutput(os.Stdout)
		hlog.SetLevel(hlog.LevelDebug)
		// 开发模式下同时输出到控制台
		log.SetOutput(io.MultiWriter(os.Stdout, &logWriter{}))
	} else {
		// Customizable output directory。
		logFilePath := cof.OutPath
		if logFilePath == "" {
			logFilePath = "./logs/"
		}
		if err := os.MkdirAll(logFilePath, 0o777); err != nil {
			log.Println(err.Error())
			return
		}

		// 轮转旧的日志文件
		if err := rotateLogFile(logFilePath, cof.FilePrefix); err != nil {
			log.Printf("Failed to rotate log file: %v", err)
		}

		// 使用固定的前缀作为当前日志文件
		fileName := getCurrentLogFile(logFilePath, cof.FilePrefix)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if _, err := os.Create(fileName); err != nil {
				log.Println(err.Error())
				return
			}
		}

		// 创建 hertz logger
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		logger := hertzzap.NewLogger(
			hertzzap.WithCoreEnc(zapcore.NewJSONEncoder(encoderConfig)),
		)

		maxSize := 20
		if cof.MaxSize != 0 {
			maxSize = int(cof.MaxSize)
		}
		maxBackups := 5
		if cof.MaxBackups != 0 {
			maxBackups = int(cof.MaxBackups)
		}
		maxAge := 10
		if cof.MaxAge != 0 {
			maxAge = int(cof.MaxAge)
		}

		// 配置日志轮转
		lumberjackLogger := &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   cof.Compress,
		}

		logger.SetOutput(lumberjackLogger)
		logger.SetLevel(hlog.Level(cof.Level))
		hlog.SetLogger(logger)

		// 将标准 log 输出重定向到 hlog
		log.SetOutput(&logWriter{})
	}

	// 设置标准 log 的输出格式
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)
}
