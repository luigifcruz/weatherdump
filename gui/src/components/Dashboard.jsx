import 'styles/dashboard';
import 'styles/grid';

import React, { Component } from 'react';

import { Link } from 'react-router-dom';
import { index as headerText } from 'static/HeaderText';

class Dashboard extends Component {
    render() {
        return (
            <div className="dashboard">
                <div className="main-header">
                    <h1 className="main-title">{headerText.title}</h1>
                    <h2 className="main-description">{headerText.description}</h2>
                </div>
                <div className="main-body grid-container">
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
