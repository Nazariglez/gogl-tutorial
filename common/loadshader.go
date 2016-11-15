/**
 * Created by nazarigonzalez on 15/11/16.
 */

package common

import (
  "github.com/go-gl/gl/v3.3-core/gl"
  "log"
  "io/ioutil"
  "fmt"
  "strings"
)

func LoadShader(vertexFilePath, fragmentFilePath string) uint32 {
  //create the shaders
  vertexShaderID := gl.CreateShader(gl.VERTEX_SHADER)
  fragmentShaderID := gl.CreateShader(gl.FRAGMENT_SHADER)

  //read the shader code from the file
  vertexShaderCode := readShaderFile(vertexFilePath)
  fragmentShaderCode := readShaderFile(fragmentFilePath)

  var result int32
  var infoLogLength int32

  //compile vertex shader
  vertexSourcePointer, free := gl.Strs(vertexShaderCode)
  defer free()

  gl.ShaderSource(vertexShaderID, 1, vertexSourcePointer, nil)
  gl.CompileShader(vertexShaderID)

  //check vertex shader
  gl.GetShaderiv(vertexShaderID, gl.COMPILE_STATUS, &result)
  gl.GetShaderiv(vertexShaderID, gl.INFO_LOG_LENGTH, &infoLogLength)
  if result != gl.TRUE && infoLogLength > 0 {
    errorLog := strings.Repeat("\x00", int(infoLogLength+1))
    gl.GetShaderInfoLog(vertexShaderID, infoLogLength, nil, gl.Str(errorLog))
    fmt.Printf("[%s]:\n%s\n", vertexFilePath, errorLog)
  }

  //compile fragment shader
  fragmentSourcePointer, free := gl.Strs(fragmentShaderCode)
  defer free()

  gl.ShaderSource(fragmentShaderID, 1, fragmentSourcePointer, nil)
  gl.CompileShader(fragmentShaderID)

  //check fragment shader
  gl.GetShaderiv(fragmentShaderID, gl.COMPILE_STATUS, &result)
  gl.GetShaderiv(fragmentShaderID, gl.INFO_LOG_LENGTH, &infoLogLength)
  if result != gl.TRUE && infoLogLength > 0 {
    errorLog := strings.Repeat("\x00", int(infoLogLength+1))
    gl.GetShaderInfoLog(fragmentShaderID, infoLogLength, nil, gl.Str(errorLog))
    fmt.Printf("[%s]:\n%s\n", fragmentFilePath, errorLog)
  }

  //link the program
  programID := gl.CreateProgram()
  gl.AttachShader(programID, vertexShaderID)
  gl.AttachShader(programID, fragmentShaderID)
  gl.LinkProgram(programID)

  //check the program
  gl.GetProgramiv(programID, gl.LINK_STATUS, &result)
  gl.GetProgramiv(programID, gl.INFO_LOG_LENGTH, &infoLogLength)
  if result != gl.TRUE && infoLogLength > 0 {
    errorLog := strings.Repeat("\x00", int(infoLogLength+1))
    gl.GetProgramInfoLog(programID, infoLogLength, nil, gl.Str(errorLog))
    fmt.Printf("[%s]:\n%s\n", "Program", errorLog)
  }

  gl.DetachShader(programID, vertexShaderID)
  gl.DetachShader(programID, fragmentShaderID)

  gl.DeleteShader(vertexShaderID)
  gl.DeleteShader(fragmentShaderID)

  return programID
}

func readShaderFile(filePath string) string {
  data, err := ioutil.ReadFile(filePath)
  if err != nil {
    log.Fatal(err)
  }

  return string(data) + "\x00"
}