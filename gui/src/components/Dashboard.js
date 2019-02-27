import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import '../styles/Dashboard.scss'

class Dashboard extends Component {

    render() {
        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">Select the satellite frequency band...</h1>
                    <h2 className="Description">
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. 
                    </h2>
                </div>
                <div className="Body Flex">
                    <div className="Satellite">
                        <h3>NPOESS</h3>
                        <h4>NOAA-20 & Suomi</h4>
                        <Link to="/steps/hrd" className="Band">X-Band HRD</Link>
                    </div>
                    <div className="Satellite">
                        <h3>Meteor</h3>
                        <h4>Meteor-MN2</h4>
                        <Link to="/steps/lrpt"  className="Band">VHF LRPT</Link>
                    </div>
                </div>
            </div>
        )
    }

}

export default Dashboard
