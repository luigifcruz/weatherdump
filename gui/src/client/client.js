import React, { Component } from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom'
import { Provider } from 'react-redux'
import { configureStore } from '../redux/store'

import App from '../components/App'
import Dashboard from '../components/Dashboard'
import Decoder from '../components/Decoder'

const store = configureStore();

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <App>
                        <Switch>
                            <Route exact path="/" component={Dashboard}/>
                            <Route exact path="/decoder" component={Decoder}/>
                        </Switch>
                    </App>
                </BrowserRouter>
            </Provider>
        );
    }
}
