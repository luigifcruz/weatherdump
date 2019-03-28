import React, { Component } from 'react'
import { connect } from 'react-redux'
import { Link } from 'react-router-dom'
import * as rxa from '../redux/actions'
import '../styles/StepPicker.scss'
import '../styles/TabView.scss'

const options = {
    hrd: {
        recorder: [

        ],
        demodulator: [

        ],
        decoder: [{
            descriptor: "soft",
            title: "Soft-Symbol File",
            description: "Demodulator output with interleaved 8-bit soft-symbols."
        },{
            descriptor: "cadu",
            title: "CADU Frames File",
            description: "Randomized and unsynchronized sequential CADU frames with ASM."
        },{
            descriptor: "asm",
            title: "CADU Frames File",
            description: "Unrandomized and synchronized sequential CADU frames with ASM."
        }],
        processor: [{
            title: "Transfer Frames File",
            description: "Decoder output with serialized CCSDS Transfer Frames."
        }]
    },
    lrpt: {
        recorder: [

        ],
        demodulator: [

        ],
        decoder: [{
            descriptor: "soft",
            title: "Soft-Symbol File",
            description: "Demodulator output with interleaved 8-bit soft-symbols."
        }],
        processor: [{
            title: "Transfer Frames File",
            description: "Decoder output with serialized CCSDS Transfer Frames."
        }]
    },
}

class StepPicker extends Component {
    constructor(props) {
        super(props);
        this.fileUpload = React.createRef();
        this.state = {};
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

        const { tab, datalink } = this.props.match.params
        switch (tab) {
            case 'decoder':
            this.props.dispatch(rxa.updateDemodulatedFile(inputFile))
            break;
            case 'processor':
            this.props.dispatch(rxa.updateDecodedFile(inputFile))
            break;
        }

        e.target.value = null;
        this.props.history.push(`/${tab}/${datalink}`)
    }

    selectInput(descriptor) {
        this.props.dispatch(rxa.updateProcessDescriptor(descriptor))
        this.fileUpload.current.click();
    }

    handleTab(datalink, to) {
        this.props.history.push(`/steps/${datalink}/${to}`)
    }

    goBack() {
        this.props.history.push(`/index.html`)
    }

    render() {
        const { tab, datalink } = this.props.match.params
        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.goBack.bind(this)} className="icon">
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
                        <Link to={`/steps/${datalink}/recorder`} className={tab == "recorder" ? "Tabs Selected" : "Tabs"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-radio"><circle cx="12" cy="12" r="2"></circle><path d="M16.24 7.76a6 6 0 0 1 0 8.49m-8.48-.01a6 6 0 0 1 0-8.49m11.31-2.82a10 10 0 0 1 0 14.14m-14.14 0a10 10 0 0 1 0-14.14"></path></svg>
                            <h3>Recorder</h3>
                        </Link>
                        <Link to={`/steps/${datalink}/demodulator`} className={tab == "demodulator" ? "Tabs Selected" : "Tabs"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-activity"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
                            <h3>Demodulator</h3>
                        </Link>
                        <Link to={`/steps/${datalink}/decoder`} className={tab == "decoder" ? "Tabs Selected" : "Tabs"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-cpu"><rect x="4" y="4" width="16" height="16" rx="2" ry="2"></rect><rect x="9" y="9" width="6" height="6"></rect><line x1="9" y1="1" x2="9" y2="4"></line><line x1="15" y1="1" x2="15" y2="4"></line><line x1="9" y1="20" x2="9" y2="23"></line><line x1="15" y1="20" x2="15" y2="23"></line><line x1="20" y1="9" x2="23" y2="9"></line><line x1="20" y1="14" x2="23" y2="14"></line><line x1="1" y1="9" x2="4" y2="9"></line><line x1="1" y1="14" x2="4" y2="14"></line></svg>
                            <h3>Decoder</h3>
                        </Link>
                        <Link to={`/steps/${datalink}/processor`} className={tab == "processor" ? "Tabs Selected" : "Tabs"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-aperture"><circle cx="12" cy="12" r="10"></circle><line x1="14.31" y1="8" x2="20.05" y2="17.94"></line><line x1="9.69" y1="8" x2="21.17" y2="8"></line><line x1="7.38" y1="12" x2="13.12" y2="2.06"></line><line x1="9.69" y1="16" x2="3.95" y2="6.06"></line><line x1="14.31" y1="16" x2="2.83" y2="16"></line><line x1="16.62" y1="12" x2="10.88" y2="21.94"></line></svg>
                            <h3>Processor</h3>
                        </Link>
                    </div>
                    <div className="TabViewBody">
                        {(Object.entries(options[datalink][tab]).length == 0) ? (
                            <div className="Option Deactivated">
                                <h3>No Options Yet</h3>
                                <h4>We're working hard to bring new features. They're coming soon!</h4>
                            </div>
                        ) : (
                        Object.entries(options[datalink][tab]).map((o, i) =>
                            <div key={i} className="Option">
                                <h3 onClick={this.selectInput.bind(this, o[1].descriptor)}>{o[1].title}</h3>
                                <h4>{o[1].description}</h4>
                            </div>
                        ))}
                    </div>
                </div>
                <input type="file" ref={this.fileUpload} onInput={this.getUploadedFileName.bind(this)} />
            </div>
        )
    }

}

StepPicker.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(StepPicker)  
