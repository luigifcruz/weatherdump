import 'styles/meta';

import React, { Component } from 'react';
import { engineVersion, version } from '../../../package.json';

class About extends Component {
	render() {
		return (
			<div className="about">
				<div className="about-left">
					<figure>
						<img className="about-left-icon" src="/icon_by_eucalyp.png" />
						<figcaption>Icon made by <a target="_blank" href="https://www.flaticon.com/authors/eucalyp">Eucalyp</a> from <a target="_blank" href="https://www.flaticon.com">Flaticon</a>.</figcaption>
					</figure>
				</div>
				<div className="about-right">
					<div className="about-right-title">WeatherDump</div>
					<div className="about-right-subtitle">by <a target="_blank" href="https://github.com/opensatelliteproject">Open Satellite Project</a></div>
					<div className="about-right-body">
						<div>Interface Version: {version}</div>
						<div>Engine Version: {engineVersion}</div>
					</div>
					<div className="about-right-body">
						<div>This program comes with absolutely no warranty.</div>
					</div>
				</div>
			</div>
        );
	}
}

export default About