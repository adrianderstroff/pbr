package main

import (
	"runtime"

	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material4/"
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
	interaction.AddInteractable(window)
	defer window.Close()

	// make camera
	camera := trackball.Make(WIDTH, HEIGHT, 10, mgl32.Vec3{0, 0, 0}, 45.0, 0.1, 100.0)
	interaction.AddInteractable(&camera)

	// make passes
	cubemappass := MakeCubemapPass(SHADER_PATH, CUBEMAP_PATH)
	renderpass := MakeRenderPass(SHADER_PATH, 10)

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// execute render passes
		cubemappass.Render(&camera)
		renderpass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}
