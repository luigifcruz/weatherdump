import React, { Component } from 'react';

import '../../styles/About.scss';
import '../../styles/TabView.scss';

class Updates extends Component {
	render() {
		const { tab } = this.props.match.params
		return (
			<div className="TabViewBody">
				This alpha version does not support automatic updates yet. Please follow updates on Twitter.
                <a target="_blank" href="https://twitter.com/luigifcruz">@luigifcruz</a>
			</div>
        );
	}
}

export default Updates