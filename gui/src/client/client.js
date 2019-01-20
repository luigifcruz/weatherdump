import React, { Component } from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom'
import { Provider } from 'react-redux'
import { configureStore } from '../redux/store'

import App from '../components/App'
import Dashboard from '../components/Dashboard'

const store = configureStore();

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <App>
                        <Switch>
                            <Route exact path="/" component={Dashboard}/>
                        </Switch>
                    </App>
                </BrowserRouter>
            </Provider>
        );
    }
}
