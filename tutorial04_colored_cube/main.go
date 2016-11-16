/**
 * Created by nazarigonzalez on 16/11/16.
 */

package main

import (
  "runtime"
  "log"

  "github.com/go-gl/glfw/v3.1/glfw"
  "github.com/go-gl/gl/v3.3-core/gl"

  "../common"
)

func main() {
  //https://github.com/go-gl/gl#usage
  runtime.LockOSThread()

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

var gColorBufferData = []float32{
  0.583,  0.771,  0.014,
  0.609,  0.115,  0.436,
  0.327,  0.483,  0.844,

  0.822,  0.569,  0.201,
  0.435,  0.602,  0.223,
  0.310,  0.747,  0.185,

  0.597,  0.770,  0.761,
  0.559,  0.436,  0.730,
  0.359,  0.583,  0.152,

  0.483,  0.596,  0.789,
  0.559,  0.861,  0.639,
  0.195,  0.548,  0.859,

  0.014,  0.184,  0.576,
  0.771,  0.328,  0.970,
  0.406,  0.615,  0.116,

  0.676,  0.977,  0.133,
  0.971,  0.572,  0.833,
  0.140,  0.616,  0.489,

  0.997,  0.513,  0.064,
  0.945,  0.719,  0.592,
  0.543,  0.021,  0.978,

  0.279,  0.317,  0.505,
  0.167,  0.620,  0.077,
  0.347,  0.857,  0.137,

  0.055,  0.953,  0.042,
  0.714,  0.505,  0.345,
  0.783,  0.290,  0.734,

  0.722,  0.645,  0.174,
  0.302,  0.455,  0.848,
  0.225,  0.587,  0.040,

  0.517,  0.713,  0.338,
  0.053,  0.959,  0.120,
  0.393,  0.621,  0.362,

  0.673,  0.211,  0.457,
  0.820,  0.883,  0.371,
  0.982,  0.099,  0.879,
}

func onWindowOpen(window *glfw.Window) {
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

  for {
    // Check if the ESC key was pressed or the window was closed
    if !(window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose()) {
      break
    }

    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(programID)

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
    gl.VertexAttribPointer(
      1, // attribute. No particular reason for 1, but must match the layout in the shader.
      3, // size
      gl.FLOAT, // type
      false, // normalized?
      0, // stride
      nil, // array buffer offset
    )

    //draw the cube!
    gl.DrawArrays(gl.TRIANGLES, 0, 12*3) // 12*3 indices starting at 0 -> 12 triangles -> 6 squares
    gl.DisableVertexAttribArray(0)

    // Swap buffers
    window.SwapBuffers()
    glfw.PollEvents()
  }
}