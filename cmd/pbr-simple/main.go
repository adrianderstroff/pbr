package main

import (
	"fmt"
	"math"
	"runtime"
	"strconv"

	"github.com/adrianderstroff/pbr/pkg/cgm"
	"github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/adrianderstroff/pbr/pkg/core/interaction"
	"github.com/adrianderstroff/pbr/pkg/core/window"
	"github.com/adrianderstroff/pbr/pkg/gui"
	"github.com/adrianderstroff/pbr/pkg/scene/camera/trackball"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	SHADER_PATH  = "./assets/shaders/"
	TEX_PATH     = "./assets/images/textures/material4/"
	CUBEMAP_PATH = "./assets/images/cubemap/hdr/"
	OBJ_PATH     = "./assets/objects/"
	OUT_PATH     = "./"

	WIDTH  int = 1200
	HEIGHT int = 800
)

func init() {
	// has to be called when using opengl context
	runtime.LockOSThread()
}

func updatePos(t float32, pos *mgl32.Vec3) {
	angle := t * 2 * math.Pi
	r := pos.Len()
	x := cgm.Cos32(angle) * r
	z := cgm.Sin32(angle) * r
	pos[0] = x
	pos[2] = z
}

func main() {
	// setup window
	title := "PBR"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	interaction.AddInteractable(window)
	defer window.Close()

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 5)
	interaction.AddInteractable(&camera)

	// make passes
	pbrpass := MakePbrPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH, OBJ_PATH)
	sunpass := MakeSunPass(WIDTH, HEIGHT, SHADER_PATH)
	//normalpass := MakeNormalPass(WIDTH, HEIGHT, SHADER_PATH)

	// setup gui
	gui := gui.New(window.Window)
	interaction.AddInteractable(gui)

	// init state
	state := State{
		albedo:         mgl32.Vec4{1, 1, 1, 1},
		roughness:      1.0,
		metalness:      0.0,
		lightpos:       mgl32.Vec3{10, 0, 10},
		lightintensity: mgl32.Vec3{100, 100, 100},
		speed:          0.0003,
		wireframe:      false,
	}

	// render loop
	var t float32 = 0
	renderloop := func() {
		// update title
		samplecount := strconv.Itoa(int(pbrpass.samples))
		params := fmt.Sprintf("%v %v - %v %v - %v", samplecount,
			pbrpass.globalroughness, pbrpass.roughness, pbrpass.metallic,
			pbrpass.lightintensity)

		window.SetTitle(title + " " + window.GetFPSFormatted() + " " + params)

		// update camera
		camera.Update()

		// update light pos
		updatePos(t, &state.lightpos)

		// execute pbr pass
		if state.wireframe {
			gl.Wireframe()
		}
		pbrpass.SetState(state)
		pbrpass.Render(&camera)

		// execute sun path
		sunpass.SetState(state)
		sunpass.Render(&camera)
		if state.wireframe {
			gl.Fill()
		}

		// execute normal pass
		//normalpass.Render(&camera)

		// render GUI
		gui.Begin()
		if open := gui.BeginWindow("Options", 0, 0, 250, 455); open {
			if open := gui.BeginGroup("Material", 135); open {
				gui.ColorPicker("albedo", &state.albedo)
				gui.SliderFloat32("roughness", &state.roughness, 0, 1, 0.01)
				gui.SliderFloat32("metalness", &state.metalness, 0, 1, 0.01)
				gui.EndGroup()
			}

			if open := gui.BeginGroup("Light", 200); open {
				//gui.Slider3("pos", &state.lightpos, -20, 20, 1)
				gui.Slider3("color", &state.lightintensity, 0, 100, 1)
				gui.SliderFloat32("speed", &state.speed, 0, 0.001, 0.00001)
				gui.EndGroup()
			}

			if open := gui.BeginGroup("Debug", 80); open {
				//gui.Slider3("pos", &state.lightpos, -20, 20, 1)
				gui.Checkbox("Wireframe", &state.wireframe)
				gui.EndGroup()
			}
		}
		gui.EndWindow()
		gui.End()

		t += state.speed
	}
	window.RunMainLoop(renderloop)
}

// State describes the gui state used for the pbr pass.
type State struct {
	// material
	albedo    mgl32.Vec4
	roughness float32
	metalness float32

	// light
	lightpos       mgl32.Vec3
	lightintensity mgl32.Vec3
	speed          float32

	// debug
	wireframe bool
}
