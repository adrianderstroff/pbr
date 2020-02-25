// Package engine provides an abstraction layer on top of OpenGL.
// It contains entities relevant for rendering.
package window

import (
	"strconv"
	"time"

	gl "github.com/adrianderstroff/pbr/pkg/core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Window takes care of window creation and interaction.
type Window struct {
	Window *glfw.Window
	Width  int
	Height int

	fpsLock float64
	lastFps float64

	loopCursor bool
}

// New returns a pointer to a Window with the specified window title and window width and height.
func New(title string, width, height int) (*Window, error) {
	// init glfw
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	// set glfw window hints
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	//glfw.WindowHint(glfw.Samples, 4)

	// create window
	glfwWindow, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}
	// actually creating the OpenGL context
	glfwWindow.MakeContextCurrent()

	// init OpenGL
	gl.Init()

	// set default values
	window := Window{
		Window:  glfwWindow,
		Width:   width,
		Height:  height,
		fpsLock: -1.0,
	}

	return &window, nil
}

// Close cleans up the Window.
func (window *Window) Close() {
	glfw.Terminate()
}

// RunMainLoop calls the specified render function each frame until the window is being closed.
func (window *Window) RunMainLoop(render func()) {
	for !window.Window.ShouldClose() {
		// set frame start
		frameStart := time.Now()
		// reset gl states
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		// render user defined function
		render()
		// swap front with back buffer
		window.Window.SwapBuffers()
		// get inputs
		glfw.PollEvents()
		// get the time after the rendering
		frameEnd := time.Now()

		// frame lock if specified
		deltaTime := frameEnd.Sub(frameStart).Seconds() * 1000.0
		timeToWait := (1000.0 / window.fpsLock) - deltaTime
		if timeToWait > 0.0 && window.fpsLock > 0.0 {
			time.Sleep(time.Duration(timeToWait/1000) * time.Second)
			deltaTime = deltaTime + timeToWait
		}
		window.lastFps = 1000.0 / deltaTime
	}
}

// LockFPS provides an upper bound for the FPS.
// The fps has to be greater than zero.
func (window *Window) LockFPS(fps float64) {
	window.fpsLock = fps
}

// GetFPS returns the fps of the previous frame.
func (window *Window) GetFPS() float64 {
	return window.lastFps
}

// GetFPSFormatted returns the fps as formatted string.
func (window *Window) GetFPSFormatted() string {
	return strconv.FormatFloat(window.GetFPS(), 'f', 0, 64) + "FPS"
}

// SetTitle updates the window title.
func (window *Window) SetTitle(title string) {
	window.Window.SetTitle(title)
}

// SetClearColor updates the color used for a new frame and when clearing a FBO.
func (window *Window) SetClearColor(r, g, b float32) {
	gl.ClearColor(r, g, b, 1.0)
}

// OnCursorPosMove is a callback handler that is called every time the cursor moves.
func (window *Window) OnCursorPosMove(x, y, dx, dy float64) bool {
	return false
}

// OnMouseButtonPress is a callback handler that is called every time a mouse button is pressed or released.
func (window *Window) OnMouseButtonPress(leftPressed, rightPressed bool) bool {
	return false
}

// OnMouseScroll is a callback handler that is called every time the mouse wheel moves.
func (window *Window) OnMouseScroll(x, y float64) bool {
	return false
}

// OnKeyPress is a callback handler that is called every time a keyboard key is pressed.
func (window *Window) OnKeyPress(key, action, mods int) bool {
	return false
}

// OnResize is a callback handler that is called every time the window is resized.
func (window *Window) OnResize(width, height int) bool {
	window.Width, window.Height = width, height
	gl.Viewport(0, 0, int32(width), int32(height))
	return false
}
