package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	alsa "github.com/Narsil/alsa-go"
	"github.com/youpy/go-wav"
	"gopkg.in/gin-gonic/gin.v1"
)

func aplay(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	r := wav.NewReader(file)
	format, err := r.Format()
	if err != nil {
		return err
	}
	if format.AudioFormat != 1 {
		return fmt.Errorf("audio format (%x) is not supported", format.AudioFormat)
	}
	var sampleFormat alsa.SampleFormat
	switch format.BitsPerSample {
	case 8:
		sampleFormat = alsa.SampleFormatU8
	case 16:
		sampleFormat = alsa.SampleFormatS16LE
	default:
		return fmt.Errorf("sample format (%x) should be 8 or 16", format.BitsPerSample)
	}

	handle := alsa.New()
	err = handle.Open("default", alsa.StreamTypePlayback, alsa.ModeBlock)
	if err != nil {
		return err
	}
	defer handle.Close()
	handle.SampleFormat = sampleFormat
	handle.SampleRate = int(format.SampleRate)
	handle.Channels = int(format.NumChannels)
	err = handle.ApplyHwParams()
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(file)
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
