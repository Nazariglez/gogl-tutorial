/**
 * Created by nazarigonzalez on 16/11/16.
 */

/*

  Draw the cube AND the triangle, at different locations. You will need to generate 2 MVP matrices, to make 2 draw calls in the main loop, but only 1 shader is required.

  Generate the color values yourself

 */

package main

import (
  "runtime"
  "log"
  "math/rand"
  "time"

  "github.com/go-gl/glfw/v3.1/glfw"
  "github.com/go-gl/gl/v3.3-core/gl"
  "github.com/go-gl/mathgl/mgl32"

  "../common"
)

func main() {
  //https://github.com/go-gl/gl#usage
  runtime.LockOSThread()

  rand.Seed(time.Now().UnixNano())
  initialize()
}

func initialize() {
  //initialize GLFW
  if err := glfw.Init(); err != nil {
    log.Fatal(err)
  }

  defer glfw.Terminate()

  glfw.WindowHint(glfw.Samples, 4) //4x antialasing
  glfw.WindowHint(glfw.ContextVersionMajor, 3) //we want OpenGL 3.3
  glfw.WindowHint(glfw.ContextVersionMinor, 3)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) //To make MacOS happy; should not be needed
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //we don't want the old OpenGL

  //Open a window and create its OpenGL context
  window, err := glfw.CreateWindow(1024, 768, "Tutorial 04", nil, nil)
  if err != nil {
    log.Fatal(err)
  }

  //Initalize GL
  window.MakeContextCurrent()
  if err := gl.Init(); err != nil {
    log.Fatal(err)
  }

  onWindowOpen(window)

}

// Our vertices. Three consecutive floats give a 3D vertex; Three consecutive vertices give a triangle.
// A cube has 6 faces with 2 triangles each, so this makes 6*2=12 triangles, and 12*3 vertices
var gVertexBufferData = []float32{
  -1.0,-1.0,-1.0,
  -1.0,-1.0, 1.0,
  -1.0, 1.0, 1.0,

  1.0, 1.0,-1.0,
  -1.0,-1.0,-1.0,
  -1.0, 1.0,-1.0,

  1.0,-1.0, 1.0,
  -1.0,-1.0,-1.0,
  1.0,-1.0,-1.0,

  1.0, 1.0,-1.0,
  1.0,-1.0,-1.0,
  -1.0,-1.0,-1.0,

  -1.0,-1.0,-1.0,
  -1.0, 1.0, 1.0,
  -1.0, 1.0,-1.0,

  1.0,-1.0, 1.0,
  -1.0,-1.0, 1.0,
  -1.0,-1.0,-1.0,

  -1.0, 1.0, 1.0,
  -1.0,-1.0, 1.0,
  1.0,-1.0, 1.0,

  1.0, 1.0, 1.0,
  1.0,-1.0,-1.0,
  1.0, 1.0,-1.0,

  1.0,-1.0,-1.0,
  1.0, 1.0, 1.0,
  1.0,-1.0, 1.0,

  1.0, 1.0, 1.0,
  1.0, 1.0,-1.0,
  -1.0, 1.0,-1.0,

  1.0, 1.0, 1.0,
  -1.0, 1.0,-1.0,
  -1.0, 1.0, 1.0,

  1.0, 1.0, 1.0,
  -1.0, 1.0, 1.0,
  1.0,-1.0, 1.0,
}

func onWindowOpen(window *glfw.Window) {
  gColorBufferData := []float32{}
  for i := 0; i < len(gVertexBufferData); i++ {
    gColorBufferData = append(gColorBufferData, rand.Float32())
  }

  var vertexArrayId uint32
  gl.GenVertexArrays(1, &vertexArrayId)
  gl.BindVertexArray(vertexArrayId)

  //this will identify our vertex buffer
  var vertexBuffer uint32
  //generate 1 buffer, put the resulting identifier in vertexbuffer
  gl.GenBuffers(1, &vertexBuffer)
  //the following comands will talk about our 'vertexbuffer' buffer
  gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)

  //size on bytes (*4 = float32)
  sizeOfData := len(gVertexBufferData)*4

  //give our vertices to OpenGL
  gl.BufferData(gl.ARRAY_BUFFER, sizeOfData, gl.Ptr(gVertexBufferData), gl.STATIC_DRAW)


  //triangle
  var triangleVertexArrayID uint32
  gl.GenVertexArrays(1, &triangleVertexArrayID)
  gl.BindVertexArray(vertexArrayId)

  gTriangleBufferData := []float32{
    -1.0, -1.0, 0.0,
    1.0, -1.0, 0.0,
    0.0, 1.0, 0.0,
  }

  var triangleBuffer uint32
  gl.GenBuffers(1, &triangleBuffer)
  gl.BindBuffer(gl.ARRAY_BUFFER, triangleBuffer)
  gl.BufferData(gl.ARRAY_BUFFER, len(gTriangleBufferData)*4, gl.Ptr(gTriangleBufferData), gl.STATIC_DRAW)


  //color
  var colorBuffer uint32
  gl.GenBuffers(1, &colorBuffer)
  gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)
  gl.BufferData(gl.ARRAY_BUFFER, len(gColorBufferData)*4, gl.Ptr(gColorBufferData), gl.STATIC_DRAW)

  gl.ClearColor(0, 0, 0.4, 0)

  // Ensure we can capture the escape key being pressed below
  window.SetInputMode(glfw.StickyKeysMode, glfw.True)

  //Create and compile our GLSL program from the shaders
  programID := common.LoadShader("vertexshader.vert", "fragmentshader.frag")

  projection := mgl32.Perspective(mgl32.DegToRad(45.0), 4.0/3.0, 0.1, 100)
  view := mgl32.LookAtV(
    mgl32.Vec3{4,3,3}, //4,3,-3
    mgl32.Vec3{0,0,-1},
    mgl32.Vec3{0,1,0},
  )

  model := mgl32.Ident4()
  mvp := projection.Mul4(view.Mul4(model))

  projection2 := mgl32.Perspective(mgl32.DegToRad(30.0), 4.0/3.0, 0.1, 100)
  view2 := mgl32.LookAtV(
    mgl32.Vec3{4,3,-3}, //4,3,-3
    mgl32.Vec3{0,0,1},
    mgl32.Vec3{0,1,0},
  )

  model2 := mgl32.Ident4()
  mvp2 := projection2.Mul4(view2.Mul4(model2))

  mvpPointer, free := gl.Strs("MVP")
  defer free()
  matrixID := gl.GetUniformLocation(programID, *mvpPointer)

  //enable depth test
  gl.Enable(gl.DEPTH_TEST)
  //accept fragment if it close to the camera than the former one
  gl.DepthFunc(gl.LESS)


  for {
    // Check if the ESC key was pressed or the window was closed
    if !(window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose()) {
      break
    }

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(programID)

    gl.UniformMatrix4fv(matrixID, 1, false, &mvp[0])

    //1rst attribute buffer : vertices
    gl.EnableVertexAttribArray(0)
    gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
    gl.VertexAttribPointer(
      0, //attribute 0. No particular reason fo 0, but must match the layout in the shader
      3, //size
      gl.FLOAT, //type
      false, //normalized?
      0, //stride
      nil, //array buffer offset
    )

    // 2nd attribute buffer : colors
    gl.EnableVertexAttribArray(1)
    gl.BindBuffer(gl.ARRAY_BUFFER, colorBuffer)

    changeColor(&gColorBufferData)

    gl.BufferData(gl.ARRAY_BUFFER, len(gColorBufferData)*4, gl.Ptr(gColorBufferData), gl.STATIC_DRAW)
    gl.VertexAttribPointer(
      1, // attribute. No particular reason for 1, but must match the layout in the shader.
      3, // size
      gl.FLOAT, // type
      false, // normalized?
      0, // stride
      nil, // array buffer offset
    )

    //triangle
    gl.EnableVertexAttribArray(2)
    gl.BindBuffer(gl.ARRAY_BUFFER, triangleBuffer)
    gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 0, nil)

    //draw the cube!
    gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // 12*3 indices starting at 0 -> 12 triangles -> 6 squares

    gl.UniformMatrix4fv(matrixID, 1, false, &mvp2[0])
    gl.DrawArrays(gl.TRIANGLES, 0, 3)

    gl.DisableVertexAttribArray(0)
    gl.DisableVertexAttribArray(1)
    gl.DisableVertexAttribArray(2)

    // Swap buffers
    window.SwapBuffers()
    glfw.PollEvents()
  }
}

func changeColor(colors *[]float32) {
  for i := 0; i < len(*colors); i++ {
    (*colors)[i] = rand.Float32()
  }
}