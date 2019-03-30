import React, { Component } from 'react';

import '../../styles/About.scss';
import '../../styles/TabView.scss';

class Feedback extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="TabViewBody">
				Your feedback is very important to us.
			</div>
        );
	}
}

export default Feedback