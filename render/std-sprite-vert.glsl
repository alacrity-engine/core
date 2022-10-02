#version 460 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec4 aColor;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec2 texCoord;
out vec4 color;

void main()
{
    gl_Position = projection * view * model * vec4(aPos.xyz, 1.0);
    texCoord = aTexCoord;
    color = aColor;
}
