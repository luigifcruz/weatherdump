import React, { Component } from 'react';

import 'styles/About';
import 'styles/tabview';

class Feedback extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="tab-view-body">
				Your feedback is very important to us.
			</div>
        );
	}
}

export default Feedback