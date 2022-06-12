package main

import (
	"findPassword/pkg/password"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	passwordGenerator.Set7zPath(`C:\Program Files\7-Zip\7z.exe`)

	var fileConter int64 = 0
	fileConterMutex := new(sync.Mutex)

    router := gin.Default()
	router.LoadHTMLGlob("server/view/*")
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(""))
	router.Use(sessions.Sessions("mysession", store))
	upload :=  router.Group("/upload")
	{	
		upload.GET("", func(ctx *gin.Context) {
			session := sessions.Default(ctx)
			status := session.Get("status")
			anwser := session.Get("answer")
			
			if anwser != nil {
				session.Delete("answer")
				session.Save()
			}
			ctx.HTML(http.StatusOK, "demo.html", gin.H{
				"status": status,
				"answer": anwser,
			})
		})

		upload.POST("", func(ctx *gin.Context) {
			file, _ := ctx.FormFile("file")
			fileConterMutex.Lock()
			fileConter++
			fileneme := strconv.FormatInt(fileConter, 10) + file.Filename
			fileConterMutex.Unlock()

			ctx.SaveUploadedFile(file, "server/tmp/" + fileneme)

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

			go func (isAllLowercase bool, isAllUppercase bool, isAllNumber bool, isAllSpecial bool, defineOthers string, regexp string) {
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

				filePath, _:= filepath.Abs("server/tmp/" + fileneme)
				ans := passwordFinder.Find7ZPassword(uint8(startLen), uint8(endLen), filePath, uint8(1))
				os.Remove(filePath)
				
				session.Set("answer", ans)
				session.Set("status", "not running")
				session.Save()

			}(isAllLowercase, isAllUppercase, isAllNumber, isAllSpecial, defineOthers, regexp)

			ctx.Redirect(http.StatusMovedPermanently, "/upload")
		})
	}
	
    router.Run(":80")
}