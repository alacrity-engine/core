#version 460 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec2 aTexCoord;
layout (location = 2) in vec4 aColor;

layout(binding = 1) uniform samplerBuffer models;
layout(binding = 3) uniform usamplerBuffer projectionsIdx;
layout(binding = 4) uniform usamplerBuffer viewsIdx;

uniform int numSprites;
uniform mat4 projections[{{ .maxNumCanvases }}];
uniform mat4 views[{{ .maxNumCanvases }}];

out vec2 texCoord;
out vec4 color;

mat4 assembleModel(int spriteIdx, samplerBuffer models) {
    mat4 model;

    model[0] = texelFetch(models, spriteIdx).r;
    model[1] = texelFetch(models, spriteIdx + 1).r;
    model[2] = texelFetch(models, spriteIdx + 2).r;
    model[3] = texelFetch(models, spriteIdx + 3).r;

    model[4] = texelFetch(models, spriteIdx + 4).r;
    model[5] = texelFetch(models, spriteIdx + 5).r;
    model[6] = texelFetch(models, spriteIdx + 6).r;
    model[7] = texelFetch(models, spriteIdx + 7).r;

    model[8] = texelFetch(models, spriteIdx + 8).r;
    model[9] = texelFetch(models, spriteIdx + 9).r;
    model[10] = texelFetch(models, spriteIdx + 10).r;
    model[11] = texelFetch(models, spriteIdx + 11).r;

    model[12] = texelFetch(models, spriteIdx + 12).r;
    model[13] = texelFetch(models, spriteIdx + 13).r;
    model[14] = texelFetch(models, spriteIdx + 14).r;
    model[15] = texelFetch(models, spriteIdx + 15).r;
}

void main() {
    int spriteIdx = mod(gl_VertexID, numSprites);
    mat4 projection = projections[texelFetch(projectionsIdx, spriteIdx).r];
    mat4 view = views[texelFetch(viewsIdx, spriteIdx).r];
    mat4 model = assembleModel(spriteIdx, models);

    gl_Position = projection * view * model * vec4(aPos.xyz, 1.0);
    texCoord = aTexCoord;
    color = aColor;
}
