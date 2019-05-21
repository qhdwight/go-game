package main

import (
	"biomquest/entities"
	"biomquest/graphics"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl64"
	"math"
	"runtime"
)

const (
	width    = 640
	height   = 480
	lookSens = 0.005
	speed    = 5
	fov      = 70
	vertPath = "lit.vert"
	fragPath = "lit.frag"
)

var (
	cubeNorms = []float32{
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,
		0.0, 0.0, -1.0,

		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,
		0.0, 0.0, 1.0,

		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,
		-1.0, 0.0, 0.0,

		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,
		1.0, 0.0, 0.0,

		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,
		0.0, -1.0, 0.0,

		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
		0.0, 1.0, 0.0,
	}
	cubeVerts = []float32{
		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,

		-0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, -0.5, 0.5,

		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, -0.5,
		-0.5, -0.5, 0.5,
		-0.5, 0.5, 0.5,

		0.5, 0.5, 0.5,
		0.5, 0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, 0.5, 0.5,

		-0.5, -0.5, -0.5,
		0.5, -0.5, -0.5,
		0.5, -0.5, 0.5,
		0.5, -0.5, 0.5,
		-0.5, -0.5, 0.5,
		-0.5, -0.5, -0.5,

		-0.5, 0.5, -0.5,
		0.5, 0.5, -0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		-0.5, 0.5, 0.5,
		-0.5, 0.5, -0.5,
	}
	ents     []*entities.Entity
	shaders  []*graphics.Shader
	programs []*graphics.Program
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := initGlfw(); err != nil {
		panic(err)
	}
	window, err := initWindow()
	if err != nil {
		panic(err)
	}
	if err := gl.Init(); err != nil {
		panic(err)
	}
	vertShader, err := newShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragShader, err := newShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	prog, err := newProgram(vertShader, fragShader)
	if err != nil {
		panic(err)
	}
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.0, 0.2, 1.0)
	modelUniform, err := prog.GetUniformLocation("model")
	viewUniform, err := prog.GetUniformLocation("view")
	projUniform, err := prog.GetUniformLocation("projection")
	lightPosUniform, err := prog.GetUniformLocation("lightPos")
	viewPosUniform, err := prog.GetUniformLocation("viewPos")
	lightColorUniform, err := prog.GetUniformLocation("lightColor")
	objectColorUniform, err := prog.GetUniformLocation("objectColor")
	cubeModel := graphics.NewModel(cubeVerts, cubeNorms)
	player := entities.Player{
		VisualEntity: entities.VisualEntity{
			Entity: entities.Entity{
				Transform: entities.Transform{
					Pos: mgl64.Vec3{0, 0, 0},
				},
			},
		},
	}
	makeCube := func(x, y, z float64) entities.VisualEntity {
		return entities.VisualEntity{
			Transform: entities.Transform{
				Pos: mgl64.Vec3{x, y, z},
			},
			Model: cubeModel,
		}
	}
	const size = 2
	visuals := [(2*size + 1) * (2*size + 1) * (2*size + 1)]entities.VisualEntity{}
	i := 0
	for x := -size; x <= size; x++ {
		for y := -size; y <= size; y++ {
			for z := -size; z <= size; z++ {
				visuals[i] = makeCube(float64(x*2), float64(y*2), float64(z*2))
				i++
			}
		}
	}
	ents = append(ents, &player.VisualEntity.Entity)
	cam := graphics.MakeCamera(fov, width, height)
	lastRenderTime := 0.0
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			window.SetShouldClose(true)
		}
		time := glfw.GetTime()
		delta := time - lastRenderTime
		lastRenderTime = time
		mouseX, mouseY := window.GetCursorPos()
		window.SetCursorPos(width/2.0, height/2.0)
		player.VisualEntity.Entity.AddInput((mouseX-width/2.0)*lookSens, (height/2.0-mouseY)*lookSens)
		playerTransform := &player.VisualEntity.Entity.Transform
		fwd, right, up := entities.CalcRelVecs(playerTransform.Pitch, playerTransform.Yaw)
		moveOpt := func(key glfw.Key, vec mgl64.Vec3) {
			if window.GetKey(key) == glfw.Press {
				playerTransform.Pos = playerTransform.Pos.Add(vec.Mul(delta * speed))
			}
		}
		moveOpt(glfw.KeyW, fwd)
		moveOpt(glfw.KeyS, fwd.Mul(-1))
		moveOpt(glfw.KeyD, right)
		moveOpt(glfw.KeyA, right.Mul(-1))
		worldUp := mgl64.Vec3{0, 1, 0}
		moveOpt(glfw.KeyLeftShift, worldUp)
		moveOpt(glfw.KeyLeftControl, worldUp.Mul(-1))
		viewMat := mgl64.LookAtV(playerTransform.Pos, playerTransform.Pos.Add(fwd), up)
		//modelMat := mgl64.Ident4()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(prog.Handle)
		graphics.SetUniformMat4(projUniform, cam.ProjMat)
		graphics.SetUniformMat4(viewUniform, viewMat)
		cameraPos := mgl64.Vec3{math.Cos(time) * 10, math.Sin(time) * math.Cos(time) * 10, math.Sin(time) * 5}
		graphics.SetUniformVec3(lightPosUniform, cameraPos)
		graphics.SetUniformVec3(viewPosUniform, playerTransform.Pos)
		graphics.SetUniformVec3(lightColorUniform, mgl64.Vec3{1, 1, 1})
		graphics.SetUniformVec3(objectColorUniform, mgl64.Vec3{1, 0.2, 0.2})
		visuals[0].Transform.Pos = cameraPos
		for _, visual := range visuals {
			translateMat := mgl64.Translate3D(visual.Transform.Pos.Elem())
			rotMat := mgl64.HomogRotate3DX(visual.Transform.Pitch).Mul4(mgl64.HomogRotate3DY(visual.Transform.Yaw))
			modelMat := translateMat.Mul4(rotMat)
			graphics.SetUniformMat4(modelUniform, modelMat)
			visual.Model.BindVertexArray()
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cubeVerts)/3))
		}
		window.SwapBuffers()
		glfw.PollEvents()
	}
	cleanup()
}

func cleanup() {
	for _, shader := range shaders {
		gl.DeleteShader(shader.Handle)
	}
	for _, program := range programs {
		gl.DeleteProgram(program.Handle)
	}
	glfw.Terminate()
}

func initGlfw() error {
	if err := glfw.Init(); err != nil {
		return err
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	return nil
}

func newShader(path string, shaderType uint32) (*graphics.Shader, error) {
	shader, err := graphics.CompileShaderFromPath(path, shaderType)
	if err != nil {
		return nil, err
	}
	shaders = append(shaders, shader)
	return shader, nil
}

func newProgram(shaders ...*graphics.Shader) (*graphics.Program, error) {
	program, err := graphics.NewProgram(shaders)
	if err != nil {
		return nil, err
	}
	programs = append(programs, program)
	return program, nil
}

func initWindow() (*glfw.Window, error) {
	window, err := glfw.CreateWindow(width, height, "test", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()
	window.SetInputMode(glfw.StickyKeysMode, gl.TRUE)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	glfw.SwapInterval(0)
	return window, nil
}
