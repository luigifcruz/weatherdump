import { withRouter } from 'react-router-dom';
import React, { Component } from 'react';
import * as rxa from 'redux/actions';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { version } from '../../package.json';

import 'styles/fonts';
import 'styles/App';

class App extends Component {
    render() {
        return (
            <div className="App">
                {this.props.children}
                <div className="Footer">
                    <div className="Left">
                        Version {version} • <Link to="/meta/about">About</Link>
                    </div>
                    <div className="Center">
                        WeatherDump
                    </div>
                    <div className="Right">
                        <a target="_blank" href="https://github.com/opensatelliteproject">Open Satellite Project</a>
                    </div>
                </div>
            </div>
        )
    }
}

App.propTypes = rxa.props
export default withRouter(connect(rxa.mapStateToProps)(App))
