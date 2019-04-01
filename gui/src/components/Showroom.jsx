import React, { Component } from 'react';
import { connect } from 'react-redux';
import * as rxa from 'redux/actions';
import Websocket from 'react-websocket';

import { showroom as headerText } from 'static/HeaderText';

import 'styles/showroom';
import 'styles/progressbar';
import 'styles/btn';
import 'styles/grid';
import 'styles/scrollbar';

class Showroom extends Component {
    constructor(props) {
        super(props);

        this.handleStatistics = this.handleStatistics.bind(this);
        this.handleAbort = this.handleAbort.bind(this);
    }

    handleAbort() {
        const { datalink } = this.props.match.params
        const { history, processId, processDescriptor } = this.props
        history.push(`/steps/${datalink}/processor`)

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

    handleFinish() {
        if (!document.hasFocus()) {
            new Notification('Processing Finished', {
                body: 'WeatherDump finished processing your file.'
            })
        }
        
        this.props.dispatch(rxa.updateProcessId(null))
    }

    handleStatistics() {

    }

    render() {
        const { tab, datalink } = this.props.match.params
        return (
            <div className="View">
                {(this.props.processId != null) ? (
                    <div>
                        <Websocket reconnect={true} debug={true} url={`ws://localhost:3000/${datalink}/${this.props.processId}/statistics`}
                            onMessage={this.handleStatistics} />
                    </div>        
                ) :  null}
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.handleAbort} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        {headerText.title}
                    </h1>
                    <h2 className="Description">{headerText.description}</h2>
               </div> 
                <div className="Body showroom">
                    <div className="products grid-container-four-two scroll-bar">
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                        <div className="product product-dark">
                            <div className="img"><img src=""></img></div>
                            <div className="title">Channel 69</div>
                            <div className="description">2330x512 • 44 MB</div>
                        </div>
                    </div>
                    <div className="controller">
                        <div className="progress-bar progress-bar-green-dark">
                            <div className="bar">
                                <div style={{ background: "#059C75", width: "50%" }} className="filler"></div>
                            </div>
                            <div className="text">
                                <div className="description">Processing CCSDS packets</div>
                                <div className="percentage">{2}%</div>
                            </div>
                        </div>
                        <div className="btn btn-orange">Open Folder</div>
                    </div>
                </div>
            </div>
        );
    }
}

Showroom.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Showroom)  
