import React, { Component } from 'react';
import { version, engineVersion } from '../../../package.json';

import '../../styles/About.scss';
import '../../styles/TabView.scss';

class About extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="TabViewBody">
				<div className="LeftContainer">
					<figure>
						<img className="MainIcon" src="/icon_by_eucalyp.png" />
						<figcaption>Icon made by <a target="_blank" href="https://www.flaticon.com/authors/eucalyp">Eucalyp</a> from <a target="_blank" href="https://www.flaticon.com">Flaticon</a>.</figcaption>
					</figure>
				</div>
				<div className="RightContainer">
					<div className="AppName">WeatherDump</div>
					<div className="AppSubtitle">by <a target="_blank" href="https://github.com/opensatelliteproject">Open Satellite Project</a></div>
					<div className="AppDescription">
						<div>Interface Version: {version}</div>
						<div>Engine Version: {engineVersion}</div>
						<div>Build Date: {BUILD_DATE}</div>
					</div>
					<div className="AppDescription">
						<div>This program comes with absolutely no warranty.</div>
					</div>
				</div>
			</div>
        );
	}
}

export default About