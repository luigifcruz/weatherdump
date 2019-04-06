import React, { Component } from 'react';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { Provider } from 'react-redux';

import { configureStore } from 'redux/store';
import App from 'components/App';
import Dashboard from 'components/Dashboard';
import Decoder from 'components/Decoder';
import Meta from 'components/Meta';
import Processor from 'components/Processor';
import StepPicker from 'components/StepPicker';
import Showroom from 'components/Showroom';

const store = configureStore();

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <App>
                        <Switch>
                            <Route exact path="/" component={Dashboard}/>
                            <Route exact path="/meta/:tab" component={Meta}/>
                            <Route exact path="/steps/:datalink/:tab" component={StepPicker}/>
                            <Route exact path="/decoder/:datalink" component={Decoder}/>
                            <Route exact path="/processor/:datalink" component={Processor}/>
                            <Route exact path="/showroom/:datalink" component={Showroom}/>
                        </Switch>
                    </App>
                </BrowserRouter>
            </Provider>
        );
    }
}
