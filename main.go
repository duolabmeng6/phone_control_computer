package main

import (
	"github.com/duolabmeng6/goefun/egin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vgo/robotgo"
)

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
	router.POST("/api/mousedblclick", handleMouseDoubleClick)
	router.POST("/api/mouserightclick", handleMouseRightClick)
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

	// 将字符串转换为整数
	xDelta, errX := strconv.Atoi(xStr)
	yDelta, errY := strconv.Atoi(yStr)

	if errX != nil || errY != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid coordinates"})
		return
	}
	//println(xDelta, yDelta)
	// 获取当前鼠标位置
	x, y := robotgo.GetMousePos()
	//println(x, y)
	// 计算新的鼠标位置
	newX := x + xDelta
	newY := y + yDelta
	//println("newX", newX, newY)
	//// 调整本地机器上的鼠标位置
	//robotgo.Move(newX, newY)
	robotgo.MoveMouseSmooth(newX, newY, 1.0, 1.0)
	// 返回成功的响应
	c.JSON(http.StatusOK, gin.H{"status": "success", "newX": newX, "newY": newY})
}

// 处理点击事件的函数
func handleMouseClick(c *gin.Context) {
	// 处理鼠标点击事件
	x, y := robotgo.GetMousePos()
	// 这里可以添加额外的逻辑
	robotgo.Click("left", false)

	c.JSON(http.StatusOK, gin.H{"status": "success", "type": "click", "x": x, "y": y})
}

// 处理双击事件的函数
func handleMouseDoubleClick(c *gin.Context) {
	// 处理鼠标双击事件
	x, y := robotgo.GetMousePos()
	// 这里可以添加额外的逻辑
	robotgo.Click("left", true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "type": "dblclick", "x": x, "y": y})
}

func handleMouseRightClick(context *gin.Context) {
	// 处理鼠标右键点击事件
	x, y := robotgo.GetMousePos()
	// 这里可以添加额外的逻辑
	robotgo.Click("right", false)

	context.JSON(http.StatusOK, gin.H{"status": "success", "type": "rightclick", "x": x, "y": y})

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
