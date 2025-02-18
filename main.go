package main

import (
	"image/color"
	"log"
	"strings"
	"voskhod/protocol"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	GRAY  color.RGBA = color.RGBA{R: 18, G: 18, B: 18, A: 255}
	WHITE            = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	RED              = color.RGBA{R: 255, G: 18, B: 18, A: 255}
	BLUE             = color.RGBA{R: 18, G: 18, B: 255, A: 255}
)

func main() {
	url := "gemini://geminiprotocol.net/"
	db, err := protocol.InitCertsDB()
	if err != nil {
		log.Fatalln(err.Error())
	}
	response, err := protocol.FetchUrl(url, db)
	if err != nil {
		log.Fatalln(err.Error())
	}

	const px = 32
	const sep = 10

	// println(response)
	rl.InitWindow(1600, 900, "voskhod")

	for !rl.WindowShouldClose() {
		lines := strings.Split(response, "\n")[1:]
		rl.BeginDrawing()
		rl.ClearBackground(GRAY)

		y := int32(sep)
		for _, line := range lines {
			col := WHITE
			if strings.HasPrefix(line, "#") {
				col = RED
			} else if strings.HasPrefix(line, "=>") {
				col = BLUE
				mouse_y := rl.GetMouseY()
				if y <= mouse_y && mouse_y <= y+px+sep && rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
					new_url := strings.TrimLeft(line, "=>\t ")
					idx := strings.IndexAny(new_url, " \t")
					if idx != -1 {
						new_url = new_url[:idx]
					}
					url = url + new_url
					println("going to", url)
					response, err = protocol.FetchUrl(url, db)
					if err != nil {
						log.Fatalln(err.Error())
					}
				}
			}
			rl.DrawText(line, sep, y, px, col)
			y += px + sep
		}

		rl.EndDrawing()
	}
}
