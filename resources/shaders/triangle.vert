#version 460
layout(location = 0) in vec3 vp;
uniform mat4 mvp;
void main() {
    gl_Position = mvp * vec4(vp, 1.0);
}