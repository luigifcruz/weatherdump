import 'styles/fonts';
import 'styles/app';

import * as rxa from 'redux/actions';

import { Link, Route, Switch } from 'react-router-dom';
import React, { Component } from 'react';

import Dashboard from 'components/Dashboard';
import Decoder from 'components/Decoder';
import Meta from 'components/Meta';
import Processor from 'components/Processor';
import Showroom from 'components/Showroom';
import StepPicker from 'components/StepPicker';
import { connect } from 'react-redux';
import { version } from '../../package.json';
import { withRouter } from 'react-router-dom';

class App extends Component {
    render() {
        return (
            <div className="main-app main-app-dark">
                <Switch>
                    <Route exact path="/" component={Dashboard}/>
                    <Route path="/meta/:tab" component={Meta}/>
                    <Route path="/steps/:datalink/:tab" component={StepPicker}/>
                    <Route path="/decoder/:datalink" component={Decoder}/>
                    <Route path="/processor/:datalink" component={Processor}/>
                    <Route path="/showroom/:datalink" component={Showroom}/>
                </Switch>
                <div className="main-footer">
                    <div className="main-footer-left">
                        Version {version}
                        <label> • </label> 
                        <Link to={{
                            pathname: '/meta/about',
                            state: {
                                previous: this.props.history.location.pathname
                            }
                        }}>About</Link>
                    </div>
                    <div className="main-footer-center">
                        WeatherDump
                    </div>
                    <div className="main-footer-right">
                        <a target="_blank" href="https://github.com/opensatelliteproject">Open Satellite Project</a>
                    </div>
                </div>
            </div>
        )
    }
}

App.propTypes = rxa.props
export default withRouter(connect(rxa.mapStateToProps)(App))
