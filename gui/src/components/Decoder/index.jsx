import React, { Component } from 'react';
import Websocket from 'react-websocket';
import * as rxa from '../../redux/actions';
import { connect } from 'react-redux';
import request from 'superagent';

import Constellation from './Constellation';
import { decoder as headerText } from 'static/HeaderText';

import 'styles/Decoder';

function _base64ToArrayBuffer(base64) {
    var wordArray = window.atob(base64);
    var len = wordArray.length,
        u8_array = new Uint8Array(len),
        offset = 0, word, i
        ;
    for (i = 0; i < len; i++) {
        word = wordArray.charCodeAt(i);
        u8_array[offset++] = (word & 0xff);
    }
    return u8_array;
}

class Decoder extends Component {
    constructor(props) {
        super(props);
        this.state = {
            complex: [],
            stats: {
                TotalBytesRead: 0.0,
                TotalBytes: 0.0,
                AverageRSCorrections: [-1, -1, -1, -1],
                AverageVitCorrections: 0,
                SignalQuality: 0,
                ReceivedPacketsPerChannel: [],
                Finished: false,
                TaskName: "Starting decoder"
            },
            n: 0
        };

        this.handleConstellation = this.handleConstellation.bind(this);
        this.handleStatistics = this.handleStatistics.bind(this);
        this.handleAbort = this.handleAbort.bind(this);
        this.openDecodedFolder = this.openDecodedFolder.bind(this);
        this.openProcessor = this.openProcessor.bind(this);
    }

    handleConstellation(data) {
        this.setState({ complex: _base64ToArrayBuffer(data), n: this.state.n + 1 })
    }

    handleStatistics(payload) {
        const stats = JSON.parse(payload)
        if (this.state.stats.Finished != this.props.Finished && stats.Finished) {
            this.handleFinish()
        }
        this.setState({ stats })
    }

    componentDidMount() {
        const { datalink } = this.props.match.params
        const { processDescriptor } = this.props

        request
            .post(`http://localhost:3000/${datalink}/${processDescriptor}/start/decoder`)
            .field("inputFile", this.props.demodulatedFile)
            .then((res) => {
                let { Code, Description } = res.body;
                this.props.dispatch(rxa.updateProcessId(Code))
                this.props.dispatch(rxa.updateDecodedFile(Description))
                
            })
            .catch((err, res) => {
                console.log(err.response.body)
                alert(err.response.body.Code);
                this.props.history.goBack()
            })
    }

    handleFinish() {
        if (!document.hasFocus()) {
            new Notification('Decoder Finished', {
                body: 'WeatherDump finished decoding your file.'
            })
        }
        
        this.props.dispatch(rxa.updateProcessId(null))
    }

    handleAbort() {
        const { datalink } = this.props.match.params
        const { history, processId, processDescriptor } = this.props
        history.push(`/steps/${datalink}/decoder`)

        if (processId != null && processDescriptor != null) {
            request
            .post(`http://localhost:3000/${datalink}/${processDescriptor}/abort/decoder`)
            .field("id", processId)
            .then((res) => {
                this.handleFinish()
                console.log("Process aborted.")
            })
            .catch(err => console.log(err))
        }
    }

    openDecodedFolder() {
        let filePath = this.props.decodedFile.split('/')
        filePath.pop()
        window.open(filePath.join('/'), '_blank');
    }

    openProcessor() {
        const { datalink } = this.props.match.params
        this.props.history.push(`/processor/${datalink}`)
    }

    render() {
        const { datalink } = this.props.match.params
        const { stats } = this.state;

        let percentage = (stats.TotalBytesRead / stats.TotalBytes) * 100
        percentage = isNaN(percentage) ? 0 : percentage

        let droppedpackets = (stats.DroppedPackets / stats.TotalPackets) * 100
        droppedpackets = isNaN(droppedpackets) ? 0 : droppedpackets

        return (
            <div>
                {(this.props.processId != null) ? (
                    <div>
                        <Websocket 
                            reconnect={true}
                            debug={process.env.NODE_ENV == 'development'}
                            url={`ws://localhost:3000/${datalink}/${this.props.processId}/constellation`}
                            onMessage={this.handleConstellation}
                        />
                        <Websocket
                            reconnect={true}
                            debug={process.env.NODE_ENV == 'development'}
                            url={`ws://localhost:3000/${datalink}/${this.props.processId}/statistics`}
                            onMessage={this.handleStatistics}
                        />
                    </div>        
                ) :  null}
                <div className="main-header">
                    <h1 className="main-title">
                        <div onClick={this.handleAbort} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        {headerText.title}
                    </h1>
                    <h2 className="main-description">{headerText.description}</h2>
                </div>
                <div className="main-body Decoder">
                    <div className="LeftWindow">
                        <Constellation
                            percentage={percentage}
                            stats={this.state.stats}
                            complex={this.state.complex}
                            n={this.state.n}
                        />
                    </div>
                    <div className="CenterWindow">
                        <div className="ReedSolomon">
                            <div className="Indicator">
                                <div className="Block">
                                    <div className="Corrections">{stats.AverageRSCorrections[0]}</div>
                                    <div className="Label">B01</div>
                                </div>
                                <div className="Block">
                                    <div className="Corrections">{stats.AverageRSCorrections[1]}</div>
                                    <div className="Label">B02</div>
                                </div>
                                <div className="Block">
                                    <div className="Corrections">{stats.AverageRSCorrections[2]}</div>
                                    <div className="Label">B03</div>
                                </div>
                                <div className="Block">
                                    <div className="Corrections">{stats.AverageRSCorrections[3]}</div>
                                    <div className="Label">B04</div>
                                </div>
                            </div>
                            <div className="Name">Reed-Solomon Corrections</div>
                        </div>
                        <div className="SignalQuality">
                            <div className="Number">{stats.SignalQuality}%</div>
                            <div className="Name">Signal Quality</div>
                        </div>
                        <div className="DroppedPackets">
                            <div className="Number">{droppedpackets.toFixed(2)}%</div>
                            <div className="Name">Dropped Packets</div>
                        </div>
                        <div className="SignalQuality">
                            <div className="Number">{stats.FrameLock ? stats.VCID : 0}</div>
                            <div className="Name">VCID</div>
                        </div>
                        <div className="DroppedPackets">
                            <div className="Number">{stats.AverageVitCorrections}/{stats.FrameBits}</div>
                            <div className="Name">Viterbi Errors</div>
                        </div>
                        <div className="LockIndicator" style={{ background: stats.FrameLock ? "#00BA8C" : "#282A37" }}>
                            {stats.FrameLock ? "LOCKED" : "UNLOCKED"}
                        </div>
                    </div>
                    <div className="RightWindow">
                        <div className="ChannelList">
                            <div className="Label">Received Packets per Channel</div>
                            {
                                stats.ReceivedPacketsPerChannel.map((received, i) => {
                                    if (received > 0) {
                                        return (
                                            <div key={i} className="Channel">
                                                <div className="VCID">{i}</div>
                                                <div className="Count">{received}</div>
                                            </div>
                                        )
                                    }
                                })
                            }
                        </div>
                        <div className="controll-box">
                            {(this.props.processId != null && this.props.decodedFile != null) ? (
                                <div onClick={this.handleAbort} className="btn btn-orange btn-large">Abort Decoding</div>
                            ) : (
                                <div>
                                    <div onClick={this.openDecodedFolder} className="btn btn-blue btn-small btn-left">Open Folder</div>
                                    <div onClick={this.openProcessor} className="btn btn-green btn-small">Next Step</div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

Decoder.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Decoder)  