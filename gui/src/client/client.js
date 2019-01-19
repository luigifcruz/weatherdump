import React, { Component } from 'react'
import { BrowserRouter, Route, Switch, Redirect } from 'react-router-dom'
import { Provider } from 'react-redux'
import { configureStore } from '../redux/store'

import App from '../components/App'
import Dashboard from '../components/Dashboard'
import TimeDate from '../components/TimeDate'
import Settings from '../components/Settings'

const store = configureStore();

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter>
                    <App>
                        <Switch>
                            <Route exact path="/" component={Dashboard}/>
                            <Route path="/clock" component={TimeDate}/>
                            <Route path="/settings" component={Settings}/>
                        </Switch>
                    </App>
                </BrowserRouter>
            </Provider>
        );
    }
}
