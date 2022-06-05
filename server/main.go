package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"strconv"
	"findPassword/pkg/password"
)

func main() {
    router := gin.Default()
	router.LoadHTMLGlob("server/view/*")
	store := cookie.NewStore([]byte("1234"))
	router.Use(sessions.Sessions("passwordFinder", store))

	upload :=  router.Group("/upload")
	{	
		upload.GET("", func(ctx *gin.Context) {
			session := sessions.Default(ctx)
			var status string
			log.Println(session.Get("status"))
			if session.Get("status") == nil {
				status = ""
			}else{
				status = "Running"
			}
			log.Println(status)
			ctx.HTML(http.StatusOK, "demo.html", gin.H{
				"status": status,
			})
		})
		
		upload.POST("", func(ctx *gin.Context) {
			file, _ := ctx.FormFile("file")
			ctx.SaveUploadedFile(file, "server/tmp/" + file.Filename)
			startLen, err := strconv.Atoi(ctx.PostForm("startLen"))
			if err != nil {
				log.Panic(err)
			}
			endLen, err := strconv.Atoi(ctx.PostForm("endLen"))
			if err != nil {
				log.Panic(err)
			}
			isAllNumber := ctx.PostForm("AllNumber") == "on"
			isAllUppercase := ctx.PostForm("AllUppercase") == "on"
			isAllLowercase := ctx.PostForm("AllLowercase") == "on"
			isAllSpecial := ctx.PostForm("AllSpecial") == "on"

			defineOthers := ctx.PostForm("defineOthers")
			regexp := ctx.PostForm("regexp")

			passwordFinder :=  new(passwordGenerator.PasswordFinder)
			passwordFinder.Init()
			passwordFinder.SetTableClear()

			if isAllLowercase {
				passwordFinder.AddTableAllLowercase()
			}
			if isAllUppercase {
				passwordFinder.AddTableAllUppercase()
			}
			if isAllNumber {
				passwordFinder.AddTableAllNumber()
			}
			if isAllSpecial {
				passwordFinder.AddTableAllSpecial()
			}
			passwordFinder.AddTable(defineOthers)
			passwordFinder.SetReg(regexp)
			
			session := sessions.Default(ctx)
			session.Set("status", "Running")
			session.Save()
			go func () {
				ans := passwordFinder.Find7ZPassword(uint8(startLen), uint8(endLen), "server/tmp/" + file.Filename, uint8(0))
				log.Println(ans)
				
			}()

			ctx.Redirect(http.StatusMovedPermanently, "/upload")
		})
	}
	
    router.Run(":80")
}