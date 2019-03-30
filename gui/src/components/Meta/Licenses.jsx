import React, { Component } from 'react';

import '../../styles/About.scss';
import '../../styles/TabView.scss';

class Licenses extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="TabViewBody">
				Open-source licenses will be put in this area.
			</div>
        );
	}
}

export default Licenses