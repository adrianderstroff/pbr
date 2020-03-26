package main

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/gui"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material-gun/"
	CUBEMAP_PATH = "./assets/images/cubemap/hdr/"
	MESH_PATH    = "./assets/objects/"
	OUT_PATH     = "./"

	WIDTH  int = 1200
	HEIGHT int = 800
)

func init() {
	// has to be called when using opengl context
	runtime.LockOSThread()
}

func main() {
	// setup window
	title := "PBR"
	window, err := window.New(title, int(WIDTH), int(HEIGHT))
	if err != nil {
		panic(err)
	}
	window.LockFPS(60)
	window.SetClearColor(0.0, 0.53, 1.0)
	interaction := interaction.New(window)
	interaction.AddInteractable(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 2)
	interaction.AddInteractable(&camera)

	// make passes
	cubemappass := MakeCubemapPass(SHADER_PATH, CUBEMAP_PATH)
	pbrpass := MakePbrPass(WIDTH, HEIGHT, MESH_PATH, SHADER_PATH, TEX_PATH, &cubemappass.cubemap)

	// setup gui
	gui := gui.New(window.Window)
	interaction.AddInteractable(gui)

	// init state
	state := State{
		roughness: 0.1,
		samples:   10,
		wireframe: false,
	}

	// render loop
	renderloop := func() {
		// update title
		samplecount := strconv.Itoa(int(state.samples))
		params := fmt.Sprintf("%v %v", samplecount, state.roughness)

		window.SetTitle(title + " " + window.GetFPSFormatted() + " " + params)

		// update camera
		camera.Update()

		// execute render passes
		cubemappass.Render(&camera)
		pbrpass.SetState(&state)
		pbrpass.Render(&camera)

		// render GUI
		gui.Begin()
		if open := gui.BeginWindow("Options", 0, 0, 250, 355); open {
			if open := gui.BeginGroup("Material", 135); open {
				gui.SliderFloat32("roughness", &state.roughness, 0, 1, 0.01)
				gui.SliderInt32("samples", &state.samples, 1, 100, 1)
				gui.EndGroup()
			}

			if open := gui.BeginGroup("Debug", 80); open {
				gui.Checkbox("Wireframe", &state.wireframe)
				gui.EndGroup()
			}
		}
		gui.EndWindow()
		gui.End()
	}
	window.RunMainLoop(renderloop)
}

// State describes the gui state used for the pbr pass.
type State struct {
	// material
	samples   int32
	roughness float32

	// debug
	wireframe bool
}
