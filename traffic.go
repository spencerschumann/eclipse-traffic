package main

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strconv"
	"time"

	"context"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
)

// TrafficScreenshot returns a screenshot of a map of the current U.S. traffic conditions.
func TrafficScreenshot() ([]byte, error) {
	var err error
	var webSession selenium.WebDriver

	cmdCtx, cmdCancel := context.WithCancel(context.Background())
	defer cmdCancel()

	port := "4444"
	cmd := exec.CommandContext(cmdCtx, "phantomjs", "--webdriver=localhost:"+port)
	err = cmd.Start()
	if err != nil {
		return nil, errors.Wrap(err, "could not launch phantomjs")
	}

	time.Sleep(time.Second)

	// Connect to the WebDriver instance running locally.
	caps := selenium.Capabilities{}
	for webSession == nil {
		time.Sleep(time.Second)
		webSession, err = selenium.NewRemote(caps, "http://localhost:"+port)
		if err != nil {
			fmt.Println("error launching webdriver:", err)
		}
	}

	wh, _ := webSession.CurrentWindowHandle()
	webSession.ResizeWindow(wh, 2000, 1500)

	url := "https://www.google.com/maps/@43.2927733,-96.7809004,5z/data=!5m1!1e1"

	if err := webSession.Get(url); err != nil {
		return nil, err
	}

	// Give the page time to finish loading
	time.Sleep(10 * time.Second)

	img, err := webSession.Screenshot()
	if err != nil {
		return nil, err
	}

	return img, nil
}

// PostProcess post-processes the screenshot image to prepare it to be used as a video frame.
func PostProcess(img image.Image, t time.Time) image.Image {
	// Create a 720p output image
	out := image.NewRGBA(image.Rect(0, 0, 1280, 720))

	gc := draw2dimg.NewGraphicContext(out)
	// Scale the image to fit the frame. Although the screenshot could have been generated
	// directly in the target size, doing so would leave unwanted UI artifacts in it.
	gc.Translate(-270, -450)
	gc.Scale(0.9, 0.9)
	gc.DrawImage(img)

	gc.SetMatrixTransform(draw2d.NewIdentityMatrix())

	gc.SetFontData(draw2d.FontData{Name: "luxi", Family: draw2d.FontFamilyMono, Style: draw2d.FontStyleBold})
	gc.SetFillColor(image.Black)
	gc.SetFontSize(14)
	gc.FillStringAt(t.String(), 900, 50)

	return out
}

func main() {
	os.Mkdir("frames", 0755)
	frameNumber := 0

	for {
		imgData, err := TrafficScreenshot()
		if err != nil {
			fmt.Println("Error creating screenshot: ", err)
			continue
		}

		img, _, err := image.Decode(bytes.NewBuffer(imgData))
		if err != nil {
			fmt.Println("Error decoding screenshot: ", err)
			continue
		}

		frame := PostProcess(img, time.Now().Round(time.Second))

		frameNumber++
		err = draw2dimg.SaveToPngFile("frames/"+strconv.Itoa(frameNumber)+".png", frame)
		if err != nil {
			fmt.Println("Error saving frame: ", err)
			continue
		}

		// Pause before the next snapshot
		time.Sleep(time.Minute)
	}
}
