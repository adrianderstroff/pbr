package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material4/"
	CUBEMAP_PATH = "./assets/images/cubemap/hdr/"
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
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 2)
	interaction.AddInteractable(&camera)

	// make passes
	cubemappass := MakeCubemapPass(SHADER_PATH, CUBEMAP_PATH)
	pbrpass := MakePbrPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH, &cubemappass.cubemap)
	interaction.AddInteractable(&pbrpass)

	// render loop
	renderloop := func() {
		// update title
		samplecount := strconv.Itoa(int(pbrpass.samples))
		params := fmt.Sprintf("%v %v - %v %v - %v", samplecount,
			pbrpass.globalroughness, pbrpass.roughness, pbrpass.metallic,
			pbrpass.lightintensity)

		window.SetTitle(title + " " + window.GetFPSFormatted() + " " + params)

		// update camera
		camera.Update()

		// execute render passes
		cubemappass.Render(&camera)
		pbrpass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}
