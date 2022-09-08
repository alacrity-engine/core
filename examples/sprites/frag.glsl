#version 330 core

uniform sampler2D spriteTexture;

in vec2 texCoord;

out vec4 FragColor;

void main()
{
    FragColor = texture(spriteTexture, texCoord);
}