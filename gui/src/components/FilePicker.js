import React, { Component } from 'react'
import { Link } from 'react-router-dom'
import '../styles/FilePicker.scss'

class FilePicker extends Component {

    render() {
        const { match: { params } } = this.props;

        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">Choose the input format...</h1>
                    <h2 className="Description">
                    Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut. 
                    </h2>
                </div>
                <div className="Body">
                    <div className="Satellite">
                        <h3>Soft-Symbol File</h3>
                        <h4>Demodulator output with interleaved 8-bit soft-symbols.</h4>
                        <Link to="decoder" className="Band">Browse File</Link>
                    </div>
                    <div className="Satellite">
                        <h3>Decoded File</h3>
                        <h4>Pre-processed output file from WeatherDump decoder.</h4>
                        <Link to="decoders" className="Band">Browse File</Link>
                    </div>
                </div>
            </div>
        )
    }

}

export default FilePicker
