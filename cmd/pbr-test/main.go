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
	TEX_PATH     = "./assets/images/textures/material1/"
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

func updatePos(a float32, pos *mgl32.Vec3) {
	angle := a * math.Pi / 180.0
	r := cgm.Sqrt32(pos[0]*pos[0] + pos[2]*pos[2])
	x := cgm.Cos32(angle) * r
	z := cgm.Sin32(angle) * r
	pos[0] = x
	pos[1] = 10
	pos[2] = z
}

func main() {
	// setup window
	title := "PBR test"
	window, _ := window.New(title, int(WIDTH), int(HEIGHT))
	window.LockFPS(60)
	interaction := interaction.New(window)
	interaction.AddInteractable(window)
	defer window.Close()

	gl.ClearColor(0.6, 0.6, 0.6, 1.0)

	// make camera
	camera := trackball.MakeDefault(WIDTH, HEIGHT, 5)

	// make passes
	pbrpass := MakePbrPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH, OBJ_PATH)
	envpass := MakeCubemapPass(SHADER_PATH, CUBEMAP_PATH)
	iblpass := MakeIblPass(WIDTH, HEIGHT, SHADER_PATH, TEX_PATH, &envpass.cubemap)
	sunpass := MakeSunPass(WIDTH, HEIGHT, SHADER_PATH)

	// setup gui
	gui := gui.New(window.Window)
	interaction.AddInteractable(gui)
	interaction.AddInteractable(&camera)

	// init state
	state := State{
		imageidx:        0,
		albedo:          mgl32.Vec4{1, 1, 1, 1},
		roughness:       1.0,
		metalness:       0.0,
		lightpos:        mgl32.Vec3{10, 10, 10},
		lightintensity:  mgl32.Vec3{100, 100, 100},
		angle:           62.0,
		samples:         10,
		globalroughness: 0.0,
		wireframe:       false,
		normal:          false,
		ibl:             false,
		bgcolor:         mgl32.Vec4{0.6, 0.6, 0.6, 1.0},
	}

	// set up options for the different displayable textures
	imgs := []string{"color", "diffuse", "specular", "ndf", "geometry", "fresnel"}
	images := imgs[:]
	mods := []string{"direct", "ibl"}
	modes := mods[:]
	var modeidx int32 = 0

	// render loop
	renderloop := func() {
		// update title
		samplecount := strconv.Itoa(int(pbrpass.samples))
		params := fmt.Sprintf("%v %v - %v %v - %v - %v", samplecount,
			pbrpass.globalroughness, pbrpass.roughness, pbrpass.metallic,
			pbrpass.lightintensity, state.angle)

		window.SetTitle(title + " " + window.GetFPSFormatted() + " " + params)

		// update camera
		camera.Update()

		// update light pos
		updatePos(state.angle, &state.lightpos)

		// execute pbr pass
		if state.wireframe {
			gl.Wireframe()
		}

		// execute environment map pass if enabled
		if state.ibl {
			// execute pbr pass
			iblpass.SetState(state)
			iblpass.Render(&camera)

			// execute environmap pass
			envpass.Render(&camera)
		} else {
			// execute pbr pass
			pbrpass.SetState(state)
			pbrpass.Render(&camera)

			// execute sun path
			sunpass.SetState(state)
			sunpass.Render(&camera)
			if state.wireframe {
				gl.Fill()
			}
		}

		// render GUI
		gui.Begin()
		if open := gui.BeginWindow("Options", 0, 0, 250, float32(HEIGHT)); open {
			if open := gui.BeginGroup("Render", 105); open {
				gui.Selector("mode", modes, &modeidx)
				gui.Selector("display", images, &state.imageidx)
				gui.EndGroup()
			}
			state.ibl = modeidx == 1

			if !state.ibl {

				if open := gui.BeginGroup("Other", 65); open {
					gui.ColorPicker("background", &state.bgcolor)
					gui.EndGroup()
				}

				if open := gui.BeginGroup("Material", 135); open {
					gui.ColorPicker("albedo", &state.albedo)
					gui.SliderFloat32("roughness", &state.roughness, 0, 1, 0.1)
					gui.SliderFloat32("metalness", &state.metalness, 0, 1, 0.1)
					gui.EndGroup()
				}

				if open := gui.BeginGroup("Light", 200); open {
					gui.Slider3("color", &state.lightintensity, 0, 100, 1)
					gui.SliderFloat32("angle", &state.angle, 0, 360, 1.0)
					gui.EndGroup()
				}
			} else {
				if open := gui.BeginGroup("Mat", 135); open {
					gui.SliderFloat32("glob roughness", &state.globalroughness, 0, 1, 0.1)
					gui.SliderInt32("samples", &state.samples, 1, 50, 1)
					gui.EndGroup()
				}
			}

			if open := gui.BeginGroup("Debug", 140); open {
				gui.Checkbox("Wireframe", &state.wireframe)
				gui.Checkbox("Render Normals", &state.normal)
				gui.EndGroup()
			}
		}
		gui.EndWindow()
		gui.End()

		// set background color
		gl.ClearColor(state.bgcolor.X(), state.bgcolor.Y(), state.bgcolor.Z(), state.bgcolor.W())
	}
	window.RunMainLoop(renderloop)
}

// State describes the gui state used for the pbr pass.
type State struct {
	// render
	imageidx int32

	// material
	albedo    mgl32.Vec4
	roughness float32
	metalness float32

	// light
	lightpos       mgl32.Vec3
	lightintensity mgl32.Vec3
	angle          float32

	// ibl
	samples         int32
	globalroughness float32

	// debug
	wireframe bool
	normal    bool
	ibl       bool
	bgcolor   mgl32.Vec4
}
