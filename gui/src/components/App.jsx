import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import * as rxa from 'redux/actions';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { version } from '../../package.json';

import 'styles/fonts';
import 'styles/app';

class App extends Component {
    render() {
        return (
            <div className="main-app main-app-dark">
                {this.props.children}
                <div className="main-footer">
                    <div className="main-footer-left">
                        Version {version} • <Link to="/meta/about">About</Link>
                    </div>
                    <div className="main-footer-center">
                        WeatherDump
                    </div>
                    <div className="main-footer-right">
                        <a target="_blank" href="https://github.com/opensatelliteproject">Open Satellite Project</a>
                    </div>
                </div>
            </div>
        )
    }
}

App.propTypes = rxa.props
export default withRouter(connect(rxa.mapStateToProps)(App))
