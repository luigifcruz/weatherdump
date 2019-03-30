import React, { Component } from 'react';
import * as rxa from '../redux/actions';
import { connect } from 'react-redux';

import request from 'superagent';
import '../styles/Processor.scss';

class Processor extends Component {
    constructor(props) {
        super(props);

        this.goBack = this.goBack.bind(this);
        this.startProcessor = this.startProcessor.bind(this);
    }

    componentDidMount() {
        const { datalink } = this.props.match.params
        const { processDescriptor } = this.props
        request
            .get(`http://localhost:3000/${datalink}/${processDescriptor}/manifest/processor`)
            .then((res) => {
                let { Code, Description } = res.body;
                if (Code == "MANIFEST") {
                    let { Parser, Composer } = JSON.parse(Description)
                    this.props.dispatch(rxa.updateManifest(Parser, Composer))
                }
            })
            .catch((err, res) => {
                console.log(err.response.body)
                alert(err.response.body.Code);
            })
    }

    startProcessor() {
        const { datalink } = this.props.match.params
        request
            .post(`http://localhost:3000/${datalink}/${this.props.processDescriptor}/start/processor`)
            .field("inputFile", this.props.decodedFile)
            .field("pipeline", JSON.stringify(this.props.processorEnhancements))
            .field("manifest", JSON.stringify({
                Parser: this.props.manifestParser,
                Composer: this.props.manifestComposer
            }))
            .then((res) => {
                let { Code, Description } = res.body;
                this.props.dispatch(rxa.updateProcessId(Code))
                this.props.dispatch(rxa.updateWorkingFolder(Description))
                this.props.history.push(`/showroom/${datalink}`)
            })
            .catch((err, res) => {
                console.log(err.response.body)
                alert(err.response.body.Code);
            })
    }

    goBack() {
        this.props.history.push(`/index.html`)
    }
    
    render() {
        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.goBack} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        Customize processor output...
                    </h1>
                    <h2 className="Description">
                        In the decoding step, the data from the demodulator is synchronized and corrected using Error Correcting algorithms like Viterbi and Reed-Solomon. This step is computationally intensive and might take a while.
                    </h2>
                </div>
                <div className="Body Processor">
                    <div className="Channels">
                        <div className="Channel">
                            <div className="Name">Individual Bands</div>
                            <div className="List">
                            {Object.entries(this.props.manifestParser).map((p, i) => {
                                return (
                                    <div
                                        key={i}
                                        onClick={() => this.props.dispatch(rxa.toggleParserActivation(p[0]))}
                                        className={(p[1].Activated) ? "Item Active" : "Item"}>
                                        {p[1].Name}
                                    </div>
                                )
                            })}
                            </div>
                        </div>
                        <div className="Channel Last">
                            <div className="Name">Multispectral Composites</div>
                            <div className="List">
                            {Object.entries(this.props.manifestComposer).map((p, i) => {
                                return (
                                    <div
                                        key={i}
                                        onClick={() => this.props.dispatch(rxa.toggleComposerActivation(p[0]))}
                                        className={(p[1].Activated) ? "Item Active" : "Item"}>
                                        {p[1].Name}
                                    </div>
                                )
                            })}
                            </div>
                        </div>
                    </div>
                    <div className="Options">
                        <div className="Option">
                            <div className="Name">Image Enhancement</div>
                            <div className="List">
                                {Object.entries(this.props.processorEnhancements).map((p, i) => {
                                    if (!p[0].includes("Export")) {
                                        return (
                                            <div
                                                key={i}
                                                onClick={() => this.props.dispatch(rxa.toggleEnhancement(p[0]))}
                                                className={(p[1].Activated) ? "Item Active" : "Item"}>
                                                <div className="Label">{p[1].Name}</div>
                                                <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                            </div>
                                        )
                                    }
                                })}
                            </div>
                        </div>
                        <div className="Option">
                            <div className="Name">Overlay Options</div>
                        </div>
                        <div className="Option">
                            <div className="Name">Export Format</div>
                            <div className="List">
                                {Object.entries(this.props.processorEnhancements).map((p, i) => {
                                    if (p[0].includes("Export")) {
                                        return (
                                            <div
                                                key={i}
                                                onClick={() => this.props.dispatch(rxa.toggleEnhancement(p[0]))}
                                                className={(p[1].Activated) ? "Item Active" : "Item"}>
                                                <div className="Label">{p[1].Name}</div>
                                                <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                            </div>
                                        )
                                    }
                                })}
                            </div>
                        </div>
                        <div onClick={this.startProcessor} className="StartButton">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-play"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

Processor.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Processor)  