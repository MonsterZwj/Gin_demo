package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"time"
)

type UserInfo struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func loginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	//post请求，form表单获取提交的参数
	//name := c.PostForm("username")
	//pwd := c.PostForm("password")
	//name := c.DefaultPostForm("username", "xiaoming")
	//pwd := c.DefaultPostForm("password", "888")
	name, ok := c.GetPostForm("username")
	if !ok {
		name = "xiaozhang"
	}
	pwd, ok := c.GetPostForm("password")
	if !ok {
		pwd = "xiaozhang"
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Name": name,
		"Password": pwd,
	})
}

func getHttp(c *gin.Context) {
	//get请求获取参数
	//name := c.Query("query")
	//name := c.DefaultQuery("query", "somebody")
	name, ok := c.GetQuery("query")
	if !ok {
		name = "somebody"
	}
	c.JSON(http.StatusOK, gin.H{
		"name": name,
	})
}

func getPath(c *gin.Context) {
	//路径参数获取
	//URL匹配不要冲突
	year := c.Param("year")
	month := c.Param("month")
	c.JSON(http.StatusOK, gin.H{
		"year": year,
		"month": month,
	})
}

func bindPost(c *gin.Context) {
	//绑定参数获取，省略获取多个参数的写法，支持get，post，json传参，form传参 牛逼的
	var u UserInfo
	err := c.ShouldBind(&u)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		fmt.Printf("%#v\n", u)
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"username": u.Username,
			"password": u.Password,
		})
	}
}

func uploadGet(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}

func uploadPost(c *gin.Context) {
	//文件上传
	 f, err := c.FormFile("f1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	} else {
		dst := path.Join("./data", f.Filename)
		err := c.SaveUploadedFile(f, dst)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		}

	}
}

func relative(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
}

func any(c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		c.JSON(http.StatusOK, gin.H{"status": "get"})
	case http.MethodPost:
		c.JSON(http.StatusOK, gin.H{"status": "post"})
	}
}

func m1(c *gin.Context) {
	fmt.Println("m1 in...")
	start := time.Now() // 计时
	c.Next() // 调用后续处理函数
	//c.Abort() // 阻止调用后续的处理函数
	cost := time.Since(start)
	fmt.Printf("cost: %v\n", cost)
	fmt.Println("m1 out...")
}

func m2(c *gin.Context) {
	fmt.Println("m2 in...")
	//go funsaasd(c.Copy()) 如果框架中使用goroutine做并发操作，必须拷贝c，不可直接调用
	c.Set("name", "Tom")
	c.Next()
	//c.Abort()
	//return // 立即退出m2
	fmt.Println("m2 out...")
}

func authMiddleware(doCheck bool)gin.HandlerFunc{
	return func(c *gin.Context) {
		if doCheck{
			// 判断用户是否登录
			// if 登录
			// c.Next()
			// else
			// c.Abort()
		}else {
			c.Next()
		}
	}
}

func main() {
	//r := gin.New()
	r := gin.Default() // Default()默认使用Logger(), Recovery()中间件，如果不想使用就用New()
	//r.LoadHTMLFiles("./login.html", "./index.html")
	r.LoadHTMLGlob("./templates/**/*.html")
 	//get请求
	r.GET("/get", getHttp)
	//post请求
	r.GET("/login", loginGet)
	r.POST("/login", loginPost)
	//绑定参数写法，支持所有请求类型
	r.POST("/bind", bindPost)
	//文件上传
	r.GET("/upload", uploadGet)
	r.POST("/upload", uploadPost)
	//路径参数
	r.GET("blog/:year/:month", getPath)
	//重定向
	r.GET("/relative", relative)
	r.GET("/a", func(c *gin.Context) {
		c.Request.URL.Path = "/b" // 把请求的URI修改
		r.HandleContext(c) // 继续后续处理
	})
	r.GET("/b", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "b",
		})
	})
	//Any: 请求方法大集合，支持所有请求方法
	r.Any("/Any", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "haha",
		})
	})
	r.Any("/any", any)
	//404请求
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
	//路由组
	videoGroup := r.Group("/video")
	{
		videoGroup.GET("/name", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"result": "name"})
		})
		videoGroup.POST("/time", func(c *gin.Context) {
			from, ok := c.GetPostForm("from")
			if !ok {
				from = "2021/01/01"
			}
			to, ok := c.GetPostForm("to")
			if !ok {
				to = "2022/02/02"
			}
			c.JSON(http.StatusOK, gin.H{
				"to": to,
				"from": from,
			})
		})
		//路由组还可以继续嵌套路由组
		gameGroup := videoGroup.Group("/game")
		{
			gameGroup.GET("/dnf", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"game": "dnf"})
			})
		}
	}
	//中间件
	r.GET("/middleware", m1, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"result": "ok"})
	})
	//路由组中间件
	middlewareGroup := r.Group("/middleware2", authMiddleware(true))
	{
		//middlewareGroup.Use(authMiddleware(true))
		middlewareGroup.GET("/name", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"result": "ok"})
		})
	}
	r.Use(m1, m2) // 全局中间件
	r.GET("/middleware1", func(c *gin.Context) {
		// 获取m2传递的参数
		//name, ok := c.Get("name1")
		//if !ok {
		//	name = "Dane"
		//}
		name := c.MustGet("name")
		c.JSON(http.StatusOK, gin.H{"result": name})
	})

	r.Run(":8000")
}
