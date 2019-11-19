package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material1/"
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
	cubemappass := MakeCubemapPass(SHADER_PATH, CUBEMAP_PATH)
	pbrpass := MakePbrPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH, &cubemappass.cubemap)
	interaction.AddInteractable(&pbrpass)

	gl.Disable(gl.CULL_FACE)

	// render loop
	renderloop := func() {
		// update title
		samplecount := strconv.Itoa(int(pbrpass.samples))
		roughness := fmt.Sprintf("%f", pbrpass.globalroughness)
		window.SetTitle(title + " " + window.GetFPSFormatted() + " " + samplecount + " " + roughness)

		// update camera
		camera.Update()

		// execute render passes
		cubemappass.Render(&camera)
		pbrpass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}
