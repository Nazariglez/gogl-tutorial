/**
 * Created by nazarigonzalez on 13/11/16.
 */

package main

import (
  "runtime"
  "log"

  "github.com/go-gl/glfw/v3.1/glfw"
  "github.com/go-gl/gl/v3.3-core/gl"
)

func main() {
  //https://github.com/go-gl/gl#usage
  runtime.LockOSThread()

  initialize()
}

func initialize() {
  //initialize GLFW
  if err := glfw.Init(); err != nil {
    log.Panic(err)
  }

  defer glfw.Terminate()

  glfw.WindowHint(glfw.Samples, 4) //4x antialasing
  glfw.WindowHint(glfw.ContextVersionMajor, 3) //we want OpenGL 3.3
  glfw.WindowHint(glfw.ContextVersionMinor, 3)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True) //To make MacOS happy; should not be needed
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) //we don't want the old OpenGL

  //Open a window and create its OpenGL context
  window, err := glfw.CreateWindow(1024, 768, "Tutorial 01", nil, nil)
  if err != nil {
    log.Panic(err)
  }

  //Initalize GL
  window.MakeContextCurrent()
  if err := gl.Init(); err != nil {
    log.Panic(err)
  }

  onWindowOpen(window)

}

func onWindowOpen(window *glfw.Window) {
  var vertexArrayId uint32
  gl.GenVertexArrays(1, &vertexArrayId)
  gl.BindVertexArray(vertexArrayId)

  //An array of 3 vectors which represents 3 vertices
  gVertexBufferData := []float32{
    -1.0, -1.0, 0.0,
    1.0, -1.0, 0.0,
    0.0, 1.0, 0.0,
  }

  //this will identify our vertex buffer
  var vertexBuffer uint32
  //generate 1 buffer, put the resulting identifier in vertexbuffer
  gl.GenBuffers(1, &vertexBuffer)
  //the following comands will talk about our 'vertexbuffer' buffer
  gl.BindBuffer(gl.ARRAY_BUFFER, vertexBuffer)
  //give our vertices to OpenGL
  gl.BufferData(gl.ARRAY_BUFFER, len(gVertexBufferData), gl.Ptr(gVertexBufferData), gl.STATIC_DRAW)



  // Ensure we can capture the escape key being pressed below
  window.SetInputMode(glfw.StickyKeysMode, glfw.True)

  for {
    // Check if the ESC key was pressed or the window was closed
    if !(window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose()) {
      break
    }

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

    //draw the triangle!
    gl.DrawArrays(gl.TRIANGLES, 0, 3) //starting from vertex 0; 3 vertices total -> 1 triangle
    gl.DisableVertexAttribArray(0)


    // Swap buffers
    window.SwapBuffers()
    glfw.PollEvents()
  }
}