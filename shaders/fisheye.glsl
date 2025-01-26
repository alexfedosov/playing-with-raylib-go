#version 330

in vec2 fragTexCoord;

out vec4 fragColor;

uniform sampler2D texture0;
uniform vec4 colDiffuse;

uniform float strength;

const float PI = 3.1415926535;

vec4 bloom(sampler2D image, vec2 uv) {
    vec4 sum = vec4(0);
    int samples = 5;
    float intensity = 1 - strength;

    vec4 source = texture(image, uv);
    float luminance = dot(source.rgb, vec3(0.2126, 0.7152, 0.0722));
    vec4 brightPass = (luminance > 0.3) ? source : vec4(0.0);

    float quality = 2.5;
    vec2 radius = vec2(strength * 0.05);

    for(int i = -samples; i < samples; i++) {
        for(int j = -samples; j < samples; j++) {
            vec2 offset = vec2(i, j) * radius * quality;
            sum += texture(image, uv + offset);
        }
    }

    sum /= float(samples * samples * 4);
    return mix(source, source + sum * intensity, strength);
}

vec4 toGrayscale(vec4 color) {
    float gray = dot(color.rgb, vec3(0.299, 0.587, 0.114));
    return vec4(gray, gray, gray, color.a);
}

void main()
{
    float aperture = 378.0;
    float apertureHalf = 0.5 * aperture * (PI / 180.0);
    float maxFactor = sin(apertureHalf);

    vec2 uv = vec2(0);
    vec2 xy = 4.0 * fragTexCoord.xy - 2.0;
    float d = length(xy);

    if (d < (2.0 - maxFactor))
    {
        d = length(xy * maxFactor);
        float z = sqrt(1.0 - d * d);
        float r = atan(d, z) / PI;
        float phi = atan(xy.y, xy.x);

        uv.x = r * cos(phi) + 0.5;
        uv.y = r * sin(phi) + 0.5;
        uv = mix(fragTexCoord.xy, uv, strength * 0.3);
        fragColor = bloom(texture0, uv);
    }
    else
    {
        uv = fragTexCoord.xy;
        fragColor = toGrayscale(texture(texture0, uv));
    }
}