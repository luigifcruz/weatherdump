import React, { Component } from 'react';
import { Link } from 'react-router-dom';

import { index as headerText } from 'static/HeaderText';

import 'styles/dashboard';

class Dashboard extends Component {
    render() {
        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">{headerText.title}</h1>
                    <h2 className="Description">{headerText.description}</h2>
                </div>
                <div className="Body Flex dashboard">
                    <div className="sat-option">
                        <h3>NPOESS</h3>
                        <label>BETA</label>
                        <h4>NOAA-20 & Suomi</h4>
                        <Link to="/steps/hrd/decoder" className="btn btn-block btn-blue">X-Band</Link>
                    </div>
                    <div className="sat-option">
                        <h3>Meteor</h3>
                        <label>ALPHA</label>
                        <h4>Meteor-MN2</h4>
                        <Link to="/steps/lrpt/decoder" className="btn btn-block btn-blue">VHF</Link>
                    </div>
                </div>
            </div>
        )
    }
}

export default Dashboard
