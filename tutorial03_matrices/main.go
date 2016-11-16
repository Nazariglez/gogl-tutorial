/**
 * Created by nazarigonzalez on 13/11/16.
 */

package main

import (
  "runtime"
  "log"

  "github.com/go-gl/glfw/v3.1/glfw"
  "github.com/go-gl/gl/v3.3-core/gl"
  "github.com/go-gl/mathgl/mgl32"

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
  window, err := glfw.CreateWindow(1024, 768, "Tutorial 03", nil, nil)
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

  //size on bytes (*4 = float32)
  sizeOfData := len(gVertexBufferData)*4

  //give our vertices to OpenGL
  gl.BufferData(gl.ARRAY_BUFFER, sizeOfData, gl.Ptr(gVertexBufferData), gl.STATIC_DRAW)
  gl.ClearColor(0, 0, 0.4, 0)

  // Ensure we can capture the escape key being pressed below
  window.SetInputMode(glfw.StickyKeysMode, glfw.True)

  //Create and compile our GLSL program from the shaders
  programID := common.LoadShader("vertexshader.vert", "fragmentshader.frag")

  //Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
  projection := mgl32.Perspective(mgl32.DegToRad(45.0), 4.0/3.0, 0.1, 100)

  //Camera matrix
  view := mgl32.LookAt(
    4, 3, 3, //Camera is at (4,3,3), in world space
    0, 0, 0, //and looks at the origin
    0, 1, 0, //head is up (set to 0, -1, 0 to look upside-down)
  )

  /*
  view := mgl32.LookAtV(
    mgl.Vec3{4,3,3},
    mgl.Vec3{0,0,0},
    mgl.Vec3{0,1,0},
  )
   */

  //model matrix : and identity matrix (model will be at te origin)
  model := mgl32.Ident4()
  //our ModelViewProjection : multiplication of our 3 matrices
  mvp := projection.Mul4(view.Mul4(model))

  //get a handle for our "MVP" uniform
  //only during the initialisation
  mvpPointer, free := gl.Strs("MVP")
  defer free()
  matrixID := gl.GetUniformLocation(programID, *mvpPointer)


  for {
    // Check if the ESC key was pressed or the window was closed
    if !(window.GetKey(glfw.KeyEscape) != glfw.Press && !window.ShouldClose()) {
      break
    }

    //exercise1 "try changing the perspective" Press A
    if window.GetKey(glfw.KeyA) == glfw.Press {
      projection = mgl32.Perspective(mgl32.DegToRad(30.0), 4.0/4.0, 0.1, 100)
      mvp = projection.Mul4(view.Mul4(model))
    }

    //exercise2 "use orthographic projection instead perspective" Press B
    if window.GetKey(glfw.KeyB) == glfw.Press {
      projection = mgl32.Ortho(-10, 10, -10, 10, 0, 100)
      mvp = projection.Mul4(view.Mul4(model))
    }

    //exercise3 "modify ModelMatrix to translate, rotate, then scale the triangle" Translate Press T
    if window.GetKey(glfw.KeyT) == glfw.Press {
      model, mvp = translate(projection, view, model)
    }

    //exercise3 "modify ModelMatrix to translate, rotate, then scale the triangle" Rotate Press R
    if window.GetKey(glfw.KeyR) == glfw.Press {
      model, mvp = rotate(projection, view, model)
    }

    //exercise3 "modify ModelMatrix to translate, rotate, then scale the triangle" Scale Press S
    if window.GetKey(glfw.KeyS) == glfw.Press {
      model, mvp = scale(projection, view, model)
    }

    //exercise 4 scale, rotate, translate (correct order)
    if window.GetKey(glfw.KeyC) == glfw.Press {
      model, mvp = scale(projection, view, model)
      model, mvp = rotate(projection, view, model)
      model, mvp = translate(projection, view, model)
    }


    gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
    gl.UseProgram(programID)

    //send our transformation to the currently bound shader, in the "MVP" uniform
    //this is donde in the main loop since each model will have a different MVP matrix (At least for the M part)
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

    //draw the triangle!
    gl.DrawArrays(gl.TRIANGLES, 0, 3) //starting from vertex 0; 3 vertices total -> 1 triangle
    gl.DisableVertexAttribArray(0)

    // Swap buffers
    window.SwapBuffers()
    glfw.PollEvents()
  }
}

func translate(projection, view, model mgl32.Mat4) (mgl32.Mat4, mgl32.Mat4) {
  translateMatrix := mgl32.Translate3D(0.5, -0.2, 0)
  model = translateMatrix.Mul4(model)
  return model, projection.Mul4(view.Mul4(model))
}

func rotate(projection, view, model mgl32.Mat4) (mgl32.Mat4, mgl32.Mat4) {
  rotationMatrix := mgl32.HomogRotate3D(mgl32.DegToRad(10), mgl32.Vec3{0,1,0})
  model = rotationMatrix.Mul4(model)
  return model, projection.Mul4(view.Mul4(model))
}

func scale(projection, view, model mgl32.Mat4) (mgl32.Mat4, mgl32.Mat4) {
  scaleMatrix := mgl32.Scale3D(1.1, 1, 0.9)
  model = scaleMatrix.Mul4(model)
  return model, projection.Mul4(view.Mul4(model))
}