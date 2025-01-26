#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;
uniform float time;

out vec4 finalColor;

float rand(vec2 co) {
    return fract(sin(dot(co.xy ,vec2(12.9898,78.233))) * 43758.5453);
}

vec3 brightnessContrast(vec3 color, float brightness, float contrast) {
    return (color - 0.5) * contrast + 0.5 + brightness;
}

vec3 saturation(vec3 color, float adjustment) {
    float luminance = dot(color, vec3(0.299, 0.587, 0.114));
    return mix(vec3(luminance), color, adjustment);
}

vec3 warmth(vec3 color, float temperature) {
    vec3 warm = vec3(0.1, 0.0, -0.1);
    return color + warm * temperature;
}

float vignette(vec2 uv, float intensity) {
    vec2 coord = (uv - 0.5) * 2.0;
    return 1.0 - dot(coord, coord) * intensity;
}

vec3 chromaticAberration(sampler2D tex, vec2 uv, float offset) {
    vec3 color;
    color.r = texture(tex, uv + vec2(offset, 0.0)).r;
    color.g = texture(tex, uv).g;
    color.b = texture(tex, uv - vec2(offset, 0.0)).b;
    return color;
}

float tiltShift(vec2 uv) {
    float center = 0.5;
    float falloff = 0.3;
    return smoothstep(0.0, falloff, abs(uv.y - center));
}

void main()
{
    vec2 uv = fragTexCoord;

    vec3 color = chromaticAberration(texture0, uv, 0.002 *  smoothstep(0.0, 1, abs(uv.y - 0.5)));
    color *= colDiffuse.rgb * fragColor.rgb;

    float vignetteIntensity = 0.4 + sin(time * 0.10) * 0.02;
    float vig = vignette(uv, vignetteIntensity);
    color *= mix(1.0, vig, 0.3);

    color = brightnessContrast(color, 0.1, 1.2);
    color = saturation(color, 1.3);
    color = warmth(color, smoothstep(0.0, 1, abs(1 - uv.y - 0.5)));

    float grain = rand(uv + time) * 0.05;
    color += grain;

    finalColor = vec4(color, 1.0);
}