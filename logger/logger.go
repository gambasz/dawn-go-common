package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	errors "gitlab.cs.umd.edu/dawn/dawn-go-common/errors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

type RequestLog struct {
	Date       string
	PID        string
	Level      string
	RequestId  string
	Error      *errors.DawnError
	StatusCode string
	Method     string
	Path       string
}

type Request struct {
	Headers fasthttp.RequestHeader
}

type MessageLog struct {
	Date      string
	Level     string
	PID       string
	RequestId string
	Message   string
}

var LEVEL_FORMAT_STRING string = "%-5s"

func buildMessageLog(c *fiber.Ctx, message string) MessageLog {
	const layout = "2006-01-02 03:04:05"
	requestId := c.Locals("requestId")

	messageLog := MessageLog{
		Date:      time.Now().UTC().Format(layout),
		RequestId: fmt.Sprintf("%s", requestId),
		PID:       strconv.Itoa(os.Getpid()),
		Message:   message,
	}
	return messageLog
}

func cleanRequest(c *fiber.Ctx, r *fasthttp.Request) Request {
	headers := fasthttp.AcquireRequest().Header
	r.Header.CopyTo(&headers)
	fmt.Println(headers.String())
	return Request{
		Headers: headers,
	}
}

func BuildMessage(c *fiber.Ctx) RequestLog {
	const layout = "2006-01-02 03:04:05"
	requestId := c.Locals("requestId")

	message := RequestLog{
		Date:       time.Now().UTC().Format(layout),
		RequestId:  fmt.Sprintf("%s", requestId),
		Level:      "INFO",
		StatusCode: strconv.Itoa(c.Response().StatusCode()),
		Method:     c.Method(),
		Path:       c.Path(),
		PID:        strconv.Itoa(os.Getpid()),
	}
	return message
}

func LogRequest(message RequestLog) {
	logString := ""
	if message.Error != nil {
		message.Level = "ERROR"
	}
	if viper.GetString("app.logType") == "json" {
		tempLogString, _ := json.MarshalIndent(message, "", "  ")
		logString = string(tempLogString)
	} else {
		logString = fmt.Sprintf("[%s] %s %s %s %s - %s %s", fmt.Sprintf(LEVEL_FORMAT_STRING, message.Level), message.Date, message.PID, message.RequestId, message.StatusCode, message.Method, message.Path)
		if message.Error != nil {
			logString += " - " + message.Error.Error()
		}
	}
	fmt.Println(logString)
}

func FiberLogger() fiber.Handler {

	return func(c *fiber.Ctx) error {
		errHandler := c.App().Config().ErrorHandler
		chainErr := c.Next()

		message := BuildMessage(c)

		if chainErr != nil {
			dawnError := ErrorHandler(c, chainErr)
			message.Error = dawnError
		}

		LogRequest(message)

		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		return nil
	}
}

func ErrorHandler(ctx *fiber.Ctx, err error) *errors.DawnError {
	var returnError *errors.DawnError
	if e, ok := err.(*errors.DawnError); ok {
		returnError = e
	} else {
		returnError = errors.Build(err)
	}

	return returnError
}

/// LOG LEVELS

func stringToLevel(str string) int {
	switch str {
	case "TRACE":
		return 1
	case "DEBUG":
		return 2
	case "INFO":
		return 3
	}
	return 1
}

func TRACE(c *fiber.Ctx, message string) {
	if stringToLevel("TRACE") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "TRACE", message)
	}
}
func DEBUG(c *fiber.Ctx, message string) {
	if stringToLevel("DEBUG") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "DEBUG", message)
	}
}
func INFO(c *fiber.Ctx, message string) {
	if stringToLevel("INFO") >= stringToLevel(viper.GetString("app.logLevel")) {
		_log(c, "INFO", message)
	}
}

func _log(c *fiber.Ctx, level, message string) {
	lg := buildMessageLog(c, message)
	lg.Level = level
	logString := ""
	if viper.GetString("app.logType") == "json" {
		tempLogString, _ := json.MarshalIndent(lg, "", "  ")
		logString = string(tempLogString)
	} else {
		logString = fmt.Sprintf("[%s] %s %s %s %s", fmt.Sprintf(LEVEL_FORMAT_STRING, lg.Level), lg.Date, lg.PID, lg.RequestId, lg.Message)
	}
	fmt.Println(logString)
}
