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

  // Ensure we can capture the escape key being pressed below
  window.SetInputMode(glfw.StickyKeysMode, glfw.True)


  // Check if the ESC key was pressed or the window was closed
  for window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose() {
    draw(window)
  }

}

func draw(window *glfw.Window) {
  // Draw nothing, see you in tutorial 2!

  // Swap buffers
  window.SwapBuffers()
  glfw.PollEvents()
}