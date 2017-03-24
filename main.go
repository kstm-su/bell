package main

import (
	"io/ioutil"
	"os"

	alsa "github.com/Narsil/alsa-go"
	"gopkg.in/gin-gonic/gin.v1"
)

func aplay(filename string) error {
	rate := 44100
	channels := 1

	handle := alsa.New()
	err := handle.Open("default", alsa.StreamTypePlayback, alsa.ModeBlock)
	if err != nil {
		return err
	}
	handle.SampleFormat = alsa.SampleFormatU8
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

	_, err := handle.Write(buf)
	if err != nil {
		return err
	}
}

func main() {
	r := gin.Default()

	r.Get("/play", func(c *gin.Context) {
		err := aplay("/usr/local/share/bell.wav")
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOk, gin.H{})
	})

	r.Run()
}
