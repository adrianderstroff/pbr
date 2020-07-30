// hdr is a test program to load hdr images
package main

import (
	"runtime"

	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH = "./assets/shaders/"
	TEX_PATH    = "./assets/images/textures/hdr/"

	WIDTH  int = 800
	HEIGHT int = 600
)

func main() {
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "HDR"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	interaction.AddInteractable(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 2)
	interaction.AddInteractable(&camera)

	// make passes
	texmappass := MakeHDRTexPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH)

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
