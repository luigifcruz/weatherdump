import 'styles/Decoder';

import * as rxa from '../../redux/actions';

import React, { Component } from 'react';

import Constellation from './Constellation';
import { RingLoader } from 'react-spinners';
import Websocket from 'react-websocket';
import { connect } from 'react-redux';
import { decoder as headerText } from 'static/HeaderText';

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
            socketOpen: false,
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
            }
        };

        this.handleAbort = this.handleAbort.bind(this);
        this.openDecodedFolder = this.openDecodedFolder.bind(this);
        this.openProcessor = this.openProcessor.bind(this);
        this.handleSocketEvent = this.handleSocketEvent.bind(this);
        this.handleSocketMessage = this.handleSocketMessage.bind(this);
        this.datalink = this.props.match.params.datalink;
    }

    handleSocketMessage(payload) {
        const data = JSON.parse(payload)
        if (this.state.stats.Finished != this.props.Finished && data.Finished) {
            this.handleFinish()
        }
        this.setState({
            stats: data,
            complex: _base64ToArrayBuffer(data.Constellation)
        })
    }

    handleSocketEvent() {
        this.setState({ socketOpen: !this.state.socketOpen });
    }

    componentDidMount() {
        if (this.props.processId == null) {
            global.client.startDecoder({
                datalink: this.datalink,
                inputFile: this.props.demodulatedFile,
                decoder: this.props.processDescriptor
            }).then((res) => {
                this.props.dispatch(rxa.updateProcessId(res.uuid))
                this.props.dispatch(rxa.updateDecodedFile(res.outputPath))
            });
        }
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
        const { history, processId } = this.props
        history.push(`/steps/${this.datalink}/decoder`)

        if (processId != null) {
            global.client.abortTask(processId).then(() => {
                this.handleFinish()
            });
        }
    }

    openDecodedFolder() {
        let filePath = this.props.decodedFile.split('/')
        filePath.pop()
        window.open(filePath.join('/'), '_blank');
    }

    openProcessor() {
        this.props.history.push(`/processor/${this.datalink}`)
    }

    render() {
        const { stats } = this.state;

        let percentage = (stats.TotalBytesRead / stats.TotalBytes) * 100
        percentage = isNaN(percentage) ? 0 : percentage

        let droppedpackets = (stats.DroppedPackets / stats.TotalPackets) * 100
        droppedpackets = isNaN(droppedpackets) ? 0 : droppedpackets

        return (
            <div>
                <div className="main-header">
                    <h1 className="main-title">
                        <div onClick={this.handleAbort} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        {headerText.title}
                    </h1>
                    <h2 className="main-description">{headerText.description}</h2>
                </div>
                {(this.props.processId != null) ? (
                    <div>
                        <Websocket 
                            reconnect={true}
                            debug={process.env.NODE_ENV == 'development'}
                            url={`ws://${global.client.enginePath}/socket/${this.datalink}/${this.props.processId}`}
                            onMessage={this.handleSocketMessage}
                            onOpen={this.handleSocketEvent}
                            onClose={this.handleSocketEvent}
                        />
                    </div>
                ) : null}
                {(this.state.socketOpen || this.state.stats.Finished || false) ? (
                    <div className="main-body Decoder">
                        <div className="LeftWindow">
                            <Constellation
                                percentage={percentage}
                                stats={this.state.stats}
                                complex={this.state.complex}
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
                ) : (
                    <div className='sockets-loader'>
                        <RingLoader
                            sizeUnit={"px"}
                            size={100}
                            color={'#63667B'}
                        />
                    </div>
                )}
            </div>
        );
    }
}

Decoder.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Decoder)  