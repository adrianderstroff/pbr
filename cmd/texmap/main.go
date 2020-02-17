package main

import (
	"runtime"

	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material3/"
	CUBEMAP_PATH = "./assets/images/cubemap/water/"
	OUT_PATH     = "./"

	WIDTH  int = 800
	HEIGHT int = 600
)

func main() {
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "PBR"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 2)
	interaction.AddInteractable(&camera)

	// make passes
	texmappass := MakeTexMapPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH)
	interaction.AddInteractable(&texmappass)

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// execute render passes
		texmappass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}
