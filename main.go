package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	alsa "github.com/Narsil/alsa-go"
	"gopkg.in/gin-gonic/gin.v1"
)

func aplay(filename string) error {
	rate := 48000
	channels := 2

	handle := alsa.New()
	err := handle.Open("default", alsa.StreamTypePlayback, alsa.ModeBlock)
	if err != nil {
		return err
	}
	defer handle.Close()
	handle.SampleFormat = alsa.SampleFormatS16LE
	handle.SampleRate = rate
	handle.Channels = channels
	err = handle.ApplyHwParams()
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	_, err = handle.Write(buf)
	return err
}

func main() {
	token := os.Getenv("TOKEN")
	fmt.Println("token:", token)

	r := gin.Default()

	r.POST("/play", func(c *gin.Context) {
		if q, ok := c.GetPostForm("token"); !ok || q != token {
			c.JSON(http.StatusOK, gin.H{
				"attachments": []map[string]string{
					map[string]string{
						"title": "error",
						"text":  "token is invalid",
						"color": "#bf271b",
					},
				},
			})
			return
		}
		err := aplay("/usr/local/share/bell.wav")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"attachments": []map[string]string{
					map[string]string{
						"title": "error",
						"text":  err.Error(),
						"color": "#bf271b",
					},
				},
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"text": "呼び出し中です..."})
	})

	r.Run()
}
