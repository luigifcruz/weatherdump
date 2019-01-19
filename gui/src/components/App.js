import { withRouter } from 'react-router-dom'
import React, { Component } from 'react'
import Websocket from 'react-websocket'
import * as rxa from '../redux/actions'
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import { TimeEvent } from 'pondjs'
import request from 'superagent'
import prefix from 'superagent-prefix'

import '../styles/App.scss'

class App extends Component {

    render() {
        return (
            <div className="App">
                <Link to="/" className="Khronos">WeatherDump</Link>
                {this.props.children}
            </div>
        )
    }

}

App.propTypes = rxa.props
export default withRouter(connect(rxa.mapStateToProps)(App))
