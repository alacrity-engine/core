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

out int vertexID;
out vec2 texCoord;
out vec4 color;

mat4 assembleModel(int spriteIdx, samplerBuffer models) {
    mat4 model;

    model[0][0] = texelFetch(models, spriteIdx).r;
    model[0][1] = texelFetch(models, spriteIdx + 1).r;
    model[0][2] = texelFetch(models, spriteIdx + 2).r;
    model[0][3] = texelFetch(models, spriteIdx + 3).r;

    model[1][0] = texelFetch(models, spriteIdx + 4).r;
    model[1][1] = texelFetch(models, spriteIdx + 5).r;
    model[1][2] = texelFetch(models, spriteIdx + 6).r;
    model[1][3] = texelFetch(models, spriteIdx + 7).r;

    model[2][0] = texelFetch(models, spriteIdx + 8).r;
    model[2][1] = texelFetch(models, spriteIdx + 9).r;
    model[2][2] = texelFetch(models, spriteIdx + 10).r;
    model[2][3] = texelFetch(models, spriteIdx + 11).r;

    model[3][0] = texelFetch(models, spriteIdx + 12).r;
    model[3][1] = texelFetch(models, spriteIdx + 13).r;
    model[3][2] = texelFetch(models, spriteIdx + 14).r;
    model[3][3] = texelFetch(models, spriteIdx + 15).r;

    return model;
}

void main() {
    int spriteIdx = int(mod(float(gl_VertexID), float(numSprites)));
    mat4 projection = projections[texelFetch(projectionsIdx, spriteIdx).r];
    mat4 view = views[texelFetch(viewsIdx, spriteIdx).r];
    mat4 model = assembleModel(spriteIdx, models);

    gl_Position = projection * view * model * vec4(aPos.xyz, 1.0);
    vertexID = gl_VertexID;
    texCoord = aTexCoord;
    color = aColor;
}
