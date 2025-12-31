package gfx

import (
	"log"
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var window *glfw.Window

func init() {
	// OpenGL braucht OS-Thread
	runtime.LockOSThread()
}

func Init(width, height int, title string) {
	if err := glfw.Init(); err != nil {
		log.Fatalf("failed to init glfw: %v", err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	w, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		log.Fatalf("failed to create window: %v", err)
	}
	window = w
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		log.Fatalf("failed to init gl: %v", err)
	}

	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func ShouldClose() bool {
	if window == nil {
		return true
	}
	return window.ShouldClose()
}

func PollEvents() {
	glfw.PollEvents()
}

func Shutdown() {
	if window != nil {
		window.Destroy()
	}
	glfw.Terminate()
}

func BeginFrame() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func EndFrame() {
	window.SwapBuffers()
}

// Clear in RGBA (0–1)
func Clear(r, g, b, a float32) {
	gl.ClearColor(r, g, b, a)
}

// TODO: hier kannst du später VAO/VBO + Shader für echte Rechtecke bauen.
// Vorläufig lassen wir DrawRect als Platzhalter.
func DrawRect(x, y, w, h float32, r, g, b, a float32) {
	// Placeholder – später: Position in NDC umrechnen + Shader zeichnen
}
