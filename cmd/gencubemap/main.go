// gencubemap is a utility program to turn an equirectangular texture into a set
// of six cube map textures.
package main

import (
	"runtime"

	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH = "./assets/shaders/"
	IN_PATH     = "./assets/images/textures/hdr/the_sky_is_on_fire_16k.hdr"
	OUT_PATH    = "./assets/images/cubemap/sky/"

	WIDTH       int = 800
	HEIGHT      int = 600
	TEXTURE_RES int = 2048 // has to be power of 2
)

func main() {
	// has to be called when using opengl context
	runtime.LockOSThread()

	// setup window
	title := "Generate Cubemap"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	interaction.AddInteractable(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 2)
	interaction.AddInteractable(&camera)

	// make passes
	genpass := MakeGenPass(SHADER_PATH, IN_PATH, TEXTURE_RES)
	genpass.Render()
	cubemap := genpass.GetCubeMap()
	cubemappass := MakeCubemapPass(SHADER_PATH, cubemap)

	// save all cube map sides to file
	cubemapimages, err := cubemap.DownloadCubeMapImages(gl.RGB, gl.FLOAT)
	if err != nil {
		panic(err)
	}

	// save all images
	filenames := []string{
		"right",
		"left",
		"top",
		"bottom",
		"front",
		"back",
	}
	for i, img := range cubemapimages {
		img.SaveToPath(OUT_PATH + filenames[i] + ".hdr")
	}

	// render loop
	renderloop := func() {
		// update title
		window.SetTitle(title + " " + window.GetFPSFormatted())

		// update camera
		camera.Update()

		// execute render passes
		cubemappass.Render(&camera)
	}
	window.RunMainLoop(renderloop)
}
