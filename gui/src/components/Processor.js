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

    componentDidMount() {
        const { match: { params } } = this.props;
        request
            .get(`http://localhost:3000/${params.datalink}/${this.props.processDescriptor}/manifest/processor`)
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

    handleAbort() {
        this.props.history.goBack()
        if (this.props.processId != null) {
            request
            .post(`http://localhost:3000/${this.props.processDatalink}/${this.props.processDescriptor}/abort/decoder`)
            .field("id", this.props.processId)
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
        this.props.dispatch(rxa.updateProcessDatalink(null))
    }

    start() {
        const { match: { params } } = this.props;
        request
            .post(`http://localhost:3000/${params.datalink}/${this.props.processDescriptor}/start/processor`)
            .field("inputFile", inputFile)
            .then((res) => {
                this.props.dispatch(rxa.updateProcessId(res.body.Code))
                this.props.dispatch(rxa.updateProcessDatalink(params.datalink))
                this.props.dispatch(rxa.updateWorkingFolder(res.body.Description))
                this.props.history.push(`/processor/${params.datalink}`)
            })
            .catch((err, res) => {
                console.log(err.response.body)
                alert(err.response.body.Code);
            })
    }
    
    render() {
        return (
            <div className="View">
                <div className="Header">
                    <h1 className="Title">
                        <div onClick={this.handleAbort.bind(this)} className="icon">
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
                            {Object.entries(this.props.manifestParser).map((parser, i) => {
                                console.log(parser, i)
                                return (<div key={i} className="Item">{parser[1].Name}</div>)
                            })}
                            </div>
                        </div>
                        <div className="Channel Last">
                            <div className="Name">Multispectral Composites</div>
                            <div className="List">
                            {Object.entries(this.props.manifestComposer).map((parser, i) => {
                                return (<div key={i} className="Item">{parser[1].Name}</div>)
                            })}
                            </div>
                        </div>
                    </div>
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
                                    <div className="Label">Horizontally Flip Image</div>
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