import 'styles/meta';
import 'styles/tabview';

import { AnimatedSwitch, spring } from 'react-router-transition';
import { Link, Route } from 'react-router-dom';
import React, { Component } from 'react';

import About from './About';
import Advanced from './Advanced';
import Feedback from './Feedback';
import Licenses from './Licenses';
import Updates from './Updates';

function mapStyles(styles) {
    return {
        opacity: styles.opacity,
        transform: `scale(${styles.scale})`,
    };
}
  
function bounce(val) {
    return spring(val, {
        stiffness: 330,
        damping: 22,
    });
}

const bounceTransition = {
    atEnter: {
        opacity: 0,
        scale: 1.05,
    },
    atLeave: {
        opacity: bounce(0),
        scale: bounce(0.95),
    },
    atActive: {
        opacity: bounce(1),
        scale: bounce(1),
    },
};

class Meta extends Component {
    constructor(props) {
        super(props);
        this.state = {
            previous: props.location.state.previous
        }
        this.handleClose = this.handleClose.bind(this);
    }

    handleClose() {
        const { history } = this.props;
        history.push(this.state.previous);
    }

    render() {
        const { tab } = this.props.match.params
        return (
            <div>
                <div className="main-header main-header-small">
                    <h1 className="main-title">
                        <div onClick={this.handleClose} className="icon">
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                        </div>
                        WeatherDump
                    </h1>
                </div>
                <div className="meta main-body main-body-large">
                    <div className="tab-view-header">
                        <Link to={`/meta/about`} className={tab == "about" ? "tab-view-tab tab-view-tab-selected" : "tab-view-tab"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-star"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
                            <h3>About</h3>
                        </Link>
                        <Link to={`/meta/feedback`} className={tab == "feedback" ? "tab-view-tab tab-view-tab-selected" : "tab-view-tab"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-heart"><path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"></path></svg>
                            <h3>Feedback</h3>
                        </Link>
                        <Link to={`/meta/updates`} className={tab == "updates" ? "tab-view-tab tab-view-tab-selected" : "tab-view-tab"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-download-cloud"><polyline points="8 17 12 21 16 17"></polyline><line x1="12" y1="12" x2="12" y2="21"></line><path d="M20.88 18.09A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.29"></path></svg>
                            <h3>Updates</h3>
                        </Link>
                        <Link to={`/meta/licenses`} className={tab == "licenses" ? "tab-view-tab tab-view-tab-selected" : "tab-view-tab"}>
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-pen-tool"><path d="M12 19l7-7 3 3-7 7-3-3z"></path><path d="M18 13l-1.5-7.5L2 2l3.5 14.5L13 18l5-5z"></path><path d="M2 2l7.586 7.586"></path><circle cx="11" cy="11" r="2"></circle></svg>
                            <h3>Licenses</h3>
                        </Link>
                        <Link to={`/meta/advanced`} className={tab == "advanced" ? "tab-view-tab tab-view-tab-selected" : "tab-view-tab"}>               
                            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="feather feather-terminal"><polyline points="4 17 10 11 4 5"></polyline><line x1="12" y1="19" x2="20" y2="19"></line></svg>
                            <h3>Advanced</h3>
                        </Link>
                    </div>
                    <AnimatedSwitch
                        atEnter={bounceTransition.atEnter}
                        atLeave={bounceTransition.atLeave}
                        atActive={bounceTransition.atActive}
                        mapStyles={mapStyles}
                        className="tab-view-body"
                    >
                        <Route exact path="/meta/about" component={About}/>
                        <Route exact path="/meta/feedback" component={Feedback}/>
                        <Route exact path="/meta/updates" component={Updates}/>
                        <Route exact path="/meta/licenses" component={Licenses}/>
                        <Route exact path="/meta/advanced" component={Advanced}/>
                    </AnimatedSwitch>
                </div>
            </div>
        );
    }
}

export default Meta
