import React, { Component } from 'react'
import Websocket from 'react-websocket'
import * as rxa from '../redux/actions'
import { connect } from 'react-redux'
import request from 'superagent'
import '../styles/Processor.scss'

class Processor extends Component {
    constructor(props) {
        super(props);
        this.state = {
            
        };
    }

    handleStatistics(payload) {
        const stats = JSON.parse(payload)
        if (this.state.stats.Finished != this.props.Finished && stats.Finished) {
            this.handleFinish()
        }
        this.setState({ stats })
    }

    handleEvent(data) {
        console.log("[STREAM] Connected to decoder via WebSocket.");
    }

    handleOpenDecodedFolder() {
        let filePath = this.props.decodedFile.split('/')
        filePath.pop()
        window.open(filePath.join('/'), '_blank');
    }

    render() {
        const { match: { params } } = this.props;
        const { complex, n, stats } = this.state;

        let percentage = (stats.TotalBytesRead / stats.TotalBytes) * 100
        percentage = isNaN(percentage) ? 0 : percentage

        let droppedpackets = (stats.DroppedPackets / stats.TotalPackets) * 100
        droppedpackets = isNaN(droppedpackets) ? 0 : droppedpackets

        return (
            <div className="View">
                {(this.props.processId != null) ? (
                    <div>
                        <Websocket url={`ws://localhost:3000/${params.datalink}/${this.props.processId}/constellation`}
                            onOpen={this.handleEvent.bind(this)} onMessage={this.handleConstellation.bind(this)} />
                    </div>        
                ) :  null}
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
                <div className="Body Flex Processor"></div>
            </div>
        )
    }

}

Decoder.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Decoder)  