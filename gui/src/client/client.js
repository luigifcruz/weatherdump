import React, { Component } from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom'
import { Provider } from 'react-redux';
import { configureStore } from '../redux/store'
import App from '../components/App'
import Dashboard from '../components/Dashboard'
import Decoder from '../components/Decoder'
import About from '../components/About'
import Processor from '../components/Processor'
import StepPicker from '../components/StepPicker'

const store = configureStore();

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <App>
                        <Switch>
                            <Route exact path="/index.html" component={Dashboard}/>
                            <Route exact path="/about" component={About}/>
                            <Route exact path="/steps/:datalink" component={StepPicker}/>
                            <Route exact path="/decoder/:datalink" component={Decoder}/>
                            <Route exact path="/processor/:datalink" component={Processor}/>
                        </Switch>
                    </App>
                </BrowserRouter>
            </Provider>
        );
    }
}
