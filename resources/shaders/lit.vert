#version 410

layout(location = 0) in vec3 vert;
layout(location = 1) in vec3 norm;

out vec3 fragPos, transformedNorm;

uniform mat4 model, view, projection;

void main() {
    fragPos = vec3(model * vec4(vert, 1.0));
    transformedNorm = mat3(transpose(inverse(model))) * norm;
    gl_Position = projection * view * model * vec4(vert, 1.0);
}