#version 460 core

uniform sampler2D spriteTexture;

in vec2 texCoord;
in vec4 color;

out vec4 FragColor;

void main()
{
    FragColor = texture(spriteTexture, texCoord) * color;
}
