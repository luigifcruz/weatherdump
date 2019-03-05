import React, { Component } from 'react'
import Websocket from 'react-websocket'
import * as rxa from '../redux/actions'
import { connect } from 'react-redux'
import request from 'superagent'
import '../styles/Processor.scss'

class Processor extends Component {
    constructor(props) {
        super(props);
        this.state = {};
    }

    handleAbort() {
        this.props.history.goBack()
    }
    
    render() {
        const { match: { params } } = this.props;
        const { complex, n, stats } = this.state;

        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.handleAbort.bind(this)} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        Decoding the input file for NPOESS...
                    </h1>
                    <h2 className="Description">
                        In the decoding step, the data from the demodulator is synchronized and corrected using Error Correcting algorithms like Viterbi and Reed-Solomon. This step is computationally intensive and might take a while.
                    </h2>
                </div>
                <div className="Body Processor">
                    <div className="Options">
                        <div className="Option">
                            <div className="Name">Image Enhancement</div>
                            <div className="List">
                                <div className="Item">
                                    <div className="Label">Histogram Equalization</div>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                </div>
                                <div className="Item">
                                    <div className="Label">Invert Infrared Pixels</div>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                </div>
                                <div className="Item Active">
                                    <div className="Label">Flip Image</div>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                </div>
                            </div>
                        </div>
                        <div className="Option">
                            <div className="Name">Overlay Options</div>
                        </div>
                        <div className="Option">
                            <div className="Name">Export Format</div>
                            <div className="List">
                                <div className="Item">
                                    <div className="Label">Lossless PNG</div>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                </div>
                                <div className="Item">
                                    <div className="Label">Lossless JPEG</div>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                </div>
                            </div>
                        </div>
                        <div className="StartButton">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-play"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
                        </div>
                    </div>
                </div>
            </div>
        )
    }

}

Processor.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Processor)  