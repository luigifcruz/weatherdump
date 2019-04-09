import 'styles/processor';
import 'styles/matrix';
import 'styles/btn';

import * as rxa from 'redux/actions';

import React, { Component } from 'react';

import { connect } from 'react-redux';
import { processor as headerText } from 'static/HeaderText';

class Processor extends Component {
    constructor(props) {
        super(props);

        this.goBack = this.goBack.bind(this);
        this.startProcessor = this.startProcessor.bind(this);
        this.datalink = this.props.match.params.datalink;
    }

    componentDidMount() {
        global.client.getManifest(this.datalink).then((manifest) => {
            this.props.dispatch(rxa.updateManifest(manifest.Parser, manifest.Composer))
        });
    }

    startProcessor() {
        global.client.startProcessor({
            datalink: this.datalink,
            inputPath: this.props.decodedFile,
            pipeline: JSON.stringify(this.props.processorEnhancements),
            manifest: JSON.stringify({
                Parser: this.props.manifestParser,
                Composer: this.props.manifestComposer
            })
        }).then((res) => {
            this.props.dispatch(rxa.updateProcessId(res.uuid))
            this.props.dispatch(rxa.updateDemodulatedFile(res.outputPath))
        }).then(() => {
            this.props.history.push(`/showroom/${this.datalink}`)
        });
    }

    goBack() {
        this.props.history.push("/index.html")
    }

    render() {
        return (
            <div>
                <div className="main-header">
                    <h1 className="main-title">
                        <div onClick={this.goBack} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-arrow-left"><line x1="19" y1="12" x2="5" y2="12"></line><polyline points="12 19 5 12 12 5"></polyline></svg>
                        </div>
                        {headerText.title}
                    </h1>
                    <h2 className="main-description">{headerText.description}</h2>
                </div>
                <div className="main-body mtx-container processor processor-dark">
                    <div className="mtx-block mtx-block-large">
                        <div className="mtx-block-name">Individual Bands</div>
                        <div className="mtx-container mtx-flex">
                            {Object.entries(this.props.manifestParser).map((p, i) => {
                                return (
                                    <div
                                        key={i}
                                        onClick={() => this.props.dispatch(rxa.toggleParserActivation(p[0]))}
                                        className={(p[1].Activated) ? "mtx-cell cell-active" : "mtx-cell"}
                                    >
                                        {p[1].Name}
                                    </div>
                                )
                            })}
                        </div>
                    </div>
                    <div className="mtx-block mtx-block-large">
                        <div className="mtx-block-name">Multispectral Composites</div>
                        <div className="mtx-container mtx-flex">
                            {Object.entries(this.props.manifestComposer).map((p, i) => {
                                return (
                                    <div
                                        key={i}
                                        onClick={() => this.props.dispatch(rxa.toggleComposerActivation(p[0]))}
                                        className={(p[1].Activated) ? "mtx-cell cell-active" : "mtx-cell"}
                                    >
                                        {p[1].Name}
                                    </div>
                                )
                            })}
                        </div>
                    </div>
                    <div className="mtx-block mtx-block-medium">
                        <div className="mtx-block-name">Image Enhancement</div>
                        <div className="mtx-container ch-list-dark">
                            {Object.entries(this.props.processorEnhancements).map((p, i) => {
                                if (!p[0].includes("Export")) {
                                    return (
                                        <div
                                            key={i}
                                            onClick={() => this.props.dispatch(rxa.toggleEnhancement(p[0]))}
                                            className={(p[1].Activated) ? "mtx-cell cell-active" : "mtx-cell"}
                                        >
                                            <div className="Label">{p[1].Name}</div>
                                            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                        </div>
                                    )
                                }
                            })}
                        </div>
                    </div>
                    <div className="mtx-block mtx-block-medium">
                        <div className="mtx-block-name">Overlay Options</div>
                    </div>
                    <div className="mtx-block mtx-block-medium">
                        <div className="mtx-block-name">Export Format</div>
                        <div className="mtx-container ch-list-dark">
                            {Object.entries(this.props.processorEnhancements).map((p, i) => {
                                if (p[0].includes("Export")) {
                                    return (
                                        <div
                                            key={i}
                                            onClick={() => this.props.dispatch(rxa.toggleEnhancement(p[0]))}
                                            className={(p[1].Activated) ? "mtx-cell cell-active" : "mtx-cell"}
                                        >
                                            <div className="Label">{p[1].Name}</div>
                                            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-check"><polyline points="20 6 9 17 4 12"></polyline></svg>
                                        </div>
                                    )
                                }
                            })}
                        </div>
                    </div>
                    <div
                        onClick={this.startProcessor}
                        className="mtx-block-small btn btn-flex btn-blue"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-play"><polygon points="5 3 19 12 5 21 5 3"></polygon></svg>
                    </div>
                </div>
            </div>
        );
    }
}

Processor.propTypes = rxa.props
export default connect(rxa.mapStateToProps)(Processor)  