import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import '../styles/FilePicker.scss'

class FilePicker extends Component {

    render() {
        const { match: { params } } = this.props;

        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">Select the satellite frequency band...</h1>
                    <h2 className="Description">
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. 
                    </h2>
                </div>
                <div className="Body">
                    pick file name pls
                    <Link to="decoder">Continue...</Link>
                </div>
            </div>
        )
    }

}

export default FilePicker
