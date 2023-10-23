package main

import (
	"github.com/duolabmeng6/goefun/egin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
)

// 定义全局变量 G 里面有 x 和 y
var G = struct {
	X int
	Y int
}{}

func main() {
	// 创建一个Gin引擎
	router := gin.Default()

	// 设置静态文件目录
	router.Static("/static", "./static")

	// 定义路由
	router.GET("/", func(c *gin.Context) {
		//跳转到 /static/index.html
		c.Redirect(302, "/static/index.html")
	})

	// 定义处理鼠标移动的路由
	router.POST("/api/mouserelativeposition", handleMouseMove)
	router.POST("/api/mouseclick", handleMouseClick)
	router.POST("/api/message", handleMessage)

	router.POST("/api/backspace", handleBackspace)

	// 启动HTTP服务器
	router.Run(":8080")
}

// 处理鼠标移动的函数
func handleMouseMove(c *gin.Context) {
	// 获取鼠标移动的距离
	xStr := egin.I(c, "deltaX", "")
	yStr := egin.I(c, "deltaY", "")
	actionType := egin.I(c, "actionType", "")
	println(xStr, xStr)
	xDelta, errX := strconv.Atoi(xStr)
	yDelta, errY := strconv.Atoi(yStr)

	if errX != nil || errY != nil {
		println(errX.Error(), errY.Error())
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid coordinates"})
		return
	}

	if actionType == "TouchStart" {
		G.X, G.Y = robotgo.GetMousePos()
		println("TouchStart", G.X, G.Y)
		c.JSON(http.StatusOK, gin.H{"status": "success", "newX": 0, "newY": 0})

	}
	if actionType == "Mouse" {
		if G.X == 0 && G.Y == 0 {
			G.X, G.Y = robotgo.GetMousePos()
			println("Mouse", G.X, G.Y)
		}
		newX := G.X + xDelta
		newY := G.Y + yDelta
		println(newX, newY)

		//robotgo.Move(newX, newY)
		robotgo.MoveSmooth(newX, newY, 1.0, 1.0)
		c.JSON(http.StatusOK, gin.H{"status": "success", "newX": newX, "newY": newY})
	}
	if actionType == "ScrollBar" {
		//计算出是上滑还是下滑
		if yDelta > 0 {
			robotgo.ScrollMouse(yDelta, "up")
		} else {
			robotgo.ScrollMouse(-yDelta, "down")
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "newX": 0, "newY": 0})
	}

}

// 处理点击事件的函数
func handleMouseClick(c *gin.Context) {
	clickType := egin.I(c, "type", "")
	println("clickType", clickType)
	x, y := robotgo.GetMousePos()

	if clickType == "click" {
		robotgo.Click("left", false)
	}
	if clickType == "dblclick" {
		robotgo.Click("left", true)
	}
	if clickType == "rightclick" {
		robotgo.Click("right", false)

	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "type": "click", "x": x, "y": y})

}

func handleBackspace(context *gin.Context) {
	robotgo.KeyTap("backspace")
	context.JSON(http.StatusOK, gin.H{"status": "success", "type": "backspace"})
}

func handleMessage(context *gin.Context) {
	message := egin.I(context, "message", "")

	println(message)

	robotgo.TypeStr(message)
	context.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}
