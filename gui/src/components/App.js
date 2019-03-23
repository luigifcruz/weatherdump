import { withRouter } from 'react-router-dom'
import React, { Component } from 'react'
import * as rxa from '../redux/actions'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'

import '../styles/App.scss'

class App extends Component {

    render() {
        return (
            <div className="App">
                {this.props.children}
                <div className="Footer">
                    <div className="Left">
                        Alpha Version 2 • <Link to="/about">About</Link>
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
