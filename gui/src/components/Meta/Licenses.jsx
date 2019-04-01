import React, { Component } from 'react';

import 'styles/About';
import 'styles/tabview';

class Licenses extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="tab-view-body">
				Open-source licenses will be put in this area.
			</div>
        );
	}
}

export default Licenses