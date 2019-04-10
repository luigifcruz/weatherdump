import 'styles/showroom';
import 'styles/progressbar';
import 'styles/btn';
import 'styles/grid';
import 'styles/scrollbar';

import * as rxa from 'redux/actions';

import React, { Component } from 'react';

import Websocket from 'react-websocket';
import { connect } from 'react-redux';
import { showroom as headerText } from 'static/HeaderText';
import open from 'open';

class Showroom extends Component {
    constructor(props) {
        super(props);

        this.state = {
            socketOpen: false
        }

        this.openDecodedFolder = this.openDecodedFolder.bind(this);
        this.openDecodedFile = this.openDecodedFile.bind(this);
        this.decodedFilePath = this.decodedFilePath.bind(this);
        this.handleSocketMessage = this.handleSocketMessage.bind(this);
        this.handleSocketEvent = this.handleSocketEvent.bind(this);
        this.handleAbort = this.handleAbort.bind(this);
    }

    decodedFilePath(path) {
        let ext = ".jpeg";
        if (!this.props.processorEnhancements.ExportJPEG.Activated) {
            ext = ".png";
        }
        return path + ext
    }

    openDecodedFile(path) {
        (async () => {
            await open(this.decodedFilePath(path), {wait: true});
        })();
    }

    openDecodedFolder() {
        (async () => {
            await open(this.props.demodulatedFile, {wait: true});
        })();
    }

    handleAbort() {
        this.handleFinish()
        this.props.history.push("/index.html")

        // To-do: ADD DECODER ABORT WHEN AVAILABLE
    }

    handleFinish() {
        if (!document.hasFocus()) {
            new Notification('Processing Finished', {
                body: 'WeatherDump finished processing your file.'
            })
        }
        
        this.props.dispatch(rxa.updateProcessId(null))
    }

    handleSocketMessage(payload) {
        const manifest = JSON.parse(payload)
        this.props.dispatch(rxa.updateManifest(manifest.Parser, manifest.Composer))
    }

    handleSocketEvent() {
        this.setState({ socketOpen: !this.state.socketOpen });
    }

    render() {
        const { datalink } = this.props.match.params
        const { manifestComposer, manifestParser } = this.props;
        const manifestMerged = Object.assign(manifestComposer, manifestParser);
        
        let count = 0, finished = 0;
        Object.entries(manifestMerged).map((p, i) => {
            if (p[1].Finished) {
                finished++;
            }
            count++;
        });

        const ratio = (finished/count)*100;

        return (
            <div>
                {(this.props.processId != null) ? (
                    <div>
                        <Websocket
                            reconnect={true}
                            debug={true}
                            url={`ws://${global.client.enginePath}/socket/${datalink}/${this.props.processId}`}
                            onMessage={this.handleSocketMessage}
                            onOpen={this.handleSocketEvent}
                            onClose={this.handleSocketEvent}
                        />
                    </div>        
                ) :  null}
                <div className="main-header">
                    <h1 className="main-title">
                        <div onClick={this.handleAbort} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                        </div>
                        {headerText.title}
                    </h1>
                    <h2 className="main-description">{headerText.description}</h2>
               </div> 
                <div className="main-body showroom">
                    <div className="products grid-container scroll-bar">
                        {Object.entries(manifestMerged).map((p, i) => {
                            const { Filename, Finished, Name, Description } = p[1];
                            const filePath = this.decodedFilePath(Filename);

                            if (Finished && Filename != "") {
                                return (
                                    <div 
                                        key={i}
                                        onClick={() => this.openDecodedFile(Filename)}
                                        className="product product-dark"
                                    >
                                        <div className="img">
                                            <img src={`http://${global.client.enginePath}/get/thumbnail?filepath=${filePath}`}/>
                                        </div>
                                        <div className="title">{Name}</div>
                                        <div className="description">{Description}</div>
                                    </div>
                                )
                            }
                        })}
                    </div>
                    <div className="controller">
                        <div className="progress-bar progress-bar-green-dark">
                            <div className="bar">
                                <div style={{ background: "#059C75", width: ratio + "%" }} className="filler"></div>
                            </div>
                            <div className="text">
                                <div className="description">Processing packets</div>
                                <div className="percentage">{finished}/{count} {ratio.toFixed(0)}%</div>
                            </div>
                        </div>
                        <div onClick={this.openDecodedFolder} className="btn btn-orange">Open Folder</div>
                    </div>
                </div>
            </div>
        );
    }
}

Showroom.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Showroom)  
