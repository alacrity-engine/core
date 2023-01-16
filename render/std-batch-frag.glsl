#version 460 core

layout(binding = 0) uniform sampler2D batchTexture;
layout(binding = 2) uniform usamplerBuffer shouldDraw;

uniform int numSprites;

in vec2 texCoord;
in vec4 color;

out vec4 FragColor;

void main() {
    int spriteIdx = gl_VertexID % numSprites;
    int shouldDrawFlag = texelFetch(shouldDraw, spriteIdx).r;

    FragColor = texture(spriteTexture, texCoord) * color * shouldDrawFlag;
}