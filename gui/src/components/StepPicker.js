import React, { Component } from 'react'
import { connect } from 'react-redux'
import request from 'superagent'
import * as rxa from '../redux/actions'
import '../styles/StepPicker.scss'
import '../styles/TabView.scss'

const RECORDING = 0
const DEMODULATING = 1
const DECODING = 2
const PROCESSING = 3

const options = [
    [],
    [],
    [
        {
            title: "Soft-Symbol File",
            description: "Demodulator output with interleaved 8-bit soft-symbols."
        }
    ],
    [
        {
            title: "Transfer Frames File",
            description: "Decoder output with serialized CCSDS Transfer Frames."
        }
    ]
]

class StepPicker extends Component {
    constructor(props) {
        super(props);
        this.fileUpload = React.createRef();
        this.state = {
            currentTab: DECODING
        };
    }

    handleSelection(currentTab) {
        this.setState({ currentTab })
    }

    getUploadedFileName(e) {
        const inputFile = e.target.files[0].path
        if (inputFile == undefined) {
            console.log("Is this running inside a Electron application?")
            alert("Browser navigation isn't supported by this app.")
            return
        }

        switch (this.state.currentTab) {
            case DECODING:
            this.handleDecoding(inputFile)
            break;
            case PROCESSING:
            this.handleProcessing(inputFile)
            break;
        }
    }

    handleDecoding(inputFile) {
        const { match: { params } } = this.props;

        request
            .post(`http://localhost:3000/${params.datalink}/start/decoder`)
            .field("inputFile", inputFile)
            .then((res) => {
                this.props.dispatch(rxa.updateProcessId(res.body.Code))
                this.props.dispatch(rxa.updateProcessDatalink(params.datalink))
                this.props.dispatch(rxa.updateDecodedFile(res.body.Description))
                this.props.history.push(`/decoder/${params.datalink}`)
            })
            .catch(err => console.log(err))
    }

    handleProcessing(inputFile) {
        const { match: { params } } = this.props;

        request
            .post(`http://localhost:3000/${params.datalink}/start/processor`)
            .field("inputFile", inputFile)
            .then((res) => {
                this.props.dispatch(rxa.updateProcessId(res.body.Code))
                this.props.dispatch(rxa.updateProcessDatalink(params.datalink))
                this.props.dispatch(rxa.updateProcessedFolder(res.body.Description))
                this.props.history.push(`/processor/${params.datalink}`)
            })
            .catch(err => console.log(err))
    }

    selectInput() {
        this.fileUpload.current.click();
    }

    render() {
        const { match: { params } } = this.props;
        const { currentTab } = this.state

        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.props.history.goBack} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        Where are you at?
                    </h1>
                    <h2 className="Description">
                        From the recording thru processing, the WeatherDump supports a myriad of input options. To proceed, select below where are you at in the receiving process and what kind of input file you want to process.
                    </h2>
                </div>
                <div className="Body StepPicker">
                    <div className="TabViewHeader">
                        <div className={currentTab == RECORDING ? "Tabs Selected" : "Tabs"} onClick={this.handleSelection.bind(this, RECORDING)}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-radio"><circle cx="12" cy="12" r="2"></circle><path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"></path></svg>
                            <h3>Recording</h3>
                        </div>
                        <div className={currentTab == DEMODULATING ? "Tabs Selected" : "Tabs"} onClick={this.handleSelection.bind(this, DEMODULATING)}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-activity"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
                            <h3>Demodulating</h3>
                        </div>
                        <div className={currentTab == DECODING ? "Tabs Selected" : "Tabs"} onClick={this.handleSelection.bind(this, DECODING)}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-cpu"><rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect><rect x="9" y="9" width="6" height="6"></rect><line x1="9" y1="1" x2="9" y2="4"></line><line x1="15" y1="1" x2="15" y2="4"></line><line x1="9" y1="20" x2="9" y2="23"></line><line x1="15" y1="20" x2="15" y2="23"></line><line x1="20" y1="9" x2="23" y2="9"></line><line x1="20" y1="14" x2="23" y2="14"></line><line x1="1" y1="9" x2="4" y2="9"></line><line x1="1" y1="14" x2="4" y2="14"></line></svg>
                            <h3>Decoding</h3>
                        </div>
                        <div className={currentTab == PROCESSING ? "Tabs Selected" : "Tabs"} onClick={this.handleSelection.bind(this, PROCESSING)}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-aperture"><circle cx="12" cy="12" r="10"></circle><line x1="14.31" y1="8" x2="20.05" y2="17.94"></line><line x1="9.69" y1="8" x2="21.17" y2="8"></line><line x1="7.38" y1="12" x2="13.12" y2="2.06"></line><line x1="9.69" y1="16" x2="3.95" y2="6.06"></line><line x1="14.31" y1="16" x2="2.83" y2="16"></line><line x1="16.62" y1="12" x2="10.88" y2="21.94"></line></svg>
                            <h3>Processing</h3>
                        </div>
                    </div>
                    <div className="TabViewBody">
                        {(options[currentTab].length == 0) ? (
                            <div className="Option Deactivated">
                                <h3>No Options Yet</h3>
                                <h4>We're working hard to bring new features. They're coming soon!</h4>
                            </div>
                        ) : (
                        options[currentTab].map((o, i) =>
                            <div key={i} className="Option">
                                <h3 onClick={this.selectInput.bind(this)}>{o.title}</h3>
                                <h4>{o.description}</h4>
                            </div>
                        ))}
                    </div>
                </div>
                <input type="file" ref={this.fileUpload} onChange={this.getUploadedFileName.bind(this)} />
            </div>
        )
    }

}

StepPicker.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(StepPicker)  
