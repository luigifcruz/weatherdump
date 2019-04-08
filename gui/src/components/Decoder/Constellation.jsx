import "styles/progressbar";

import { GLSL, Node, Shaders } from "gl-react";
import React, { Component } from 'react';

import { Surface } from "gl-react-dom";

const shaders = Shaders.create({
    constellation: {
        frag: GLSL`
        precision highp float;
        uniform float complex[200];

        float circle(in vec2 _st, in vec2 pos, in float _radius){
            vec2 dist = _st-vec2(pos);
            return 1.-smoothstep(_radius-(_radius*0.01),
                                 _radius+(_radius*0.01),
                                 dot(dist,dist)*4.0);
        }

        float map(float x, float in_min, float in_max, float out_min, float out_max) {
            return (x - in_min) * (out_max - out_min) / (in_max - in_min) + out_min;
        }

        void main() { 
            vec3 color = vec3(0.0);
            vec2 st = gl_FragCoord.xy/250.0;
            color = mix(color, vec3(0.12974,0.13725,0.18823), 1.0);
            
            for (int i=0; i < 200; i+=2) {
                float x, y;
    
                if (complex[i] > 127.0) {
                    x = map(complex[i], 127.0, 255.0, 1.0, 0.5);
                } else {
                    x = map(complex[i], 0.0, 127.0, 0.5, 0.0);
                }

                if (complex[i+1] > 127.0) {
                    y = map(complex[i+1], 127.0, 255.0, 1.0, 0.5);
                } else {
                    y = map(complex[i+1], 0.0, 127.0, 0.5, 0.0);
                }
                
                vec2 pos = vec2(x, y);
                color = mix(color, vec3(0.14901,0.40392,1.000), circle(st, pos, 0.0005));
            }

            gl_FragColor = vec4(color, 1.0);
        }
    `}
});

class Constellation extends Component {
    render() {
        const { percentage, stats, complex } = this.props;
        return (
            <div className="Constellation">
                <Surface width={250} height={250}>
                    <Node shader={shaders.constellation} uniforms={{ complex }} />
                </Surface>
                <div className="Label">PSK Constellation</div>
                <div className="progress-bar progress-bar-orange-dark">
                    <div className="bar">
                        <div style={{ width: percentage + "%" }} className="filler"></div>
                    </div>
                    <div className="text">
                        <div className="description">{stats.TaskName}</div>
                        <div className="percentage">{percentage.toFixed(2)}%</div>
                    </div>
                </div>
            </div>
        );
    }
}

export default Constellation