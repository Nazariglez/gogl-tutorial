#version 330 core

out vec3 color;
//interpolated value from the vertex shaders
in vec3 fragmentColor;

void main(){
  //Output color = color specified in the vertex shader,
  //interpolated between all 3 surronding vertices
  color = fragmentColor;
}