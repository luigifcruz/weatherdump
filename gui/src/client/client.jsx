import React, { Component } from 'react';

import App from 'components/App';
import { BrowserRouter } from 'react-router-dom';
import { Provider } from 'react-redux';
import { configureStore } from 'redux/store';
import { WeatherServer, WeatherClient } from 'weather-dump';
import getPort from 'get-port';

const store = configureStore();

window.addEventListener('load', (event) => {
    if (process.env.NODE_ENV === "development") {
        global.client = new WeatherClient("localhost", "3000");
        return;
    }

    (async () => {
        let enginePort = await getPort({
            host: "127.0.0.1",
            port: getPort.makeRange(3050, 3150)
        });
    
        if (global.server === undefined) {
            global.server = new WeatherServer(enginePort);
            global.client = new WeatherClient("localhost", enginePort);
            global.server.startEngine();
        }
    })();
});

window.addEventListener('unload', (event) => {
    global.server.stopEngine();
});

export default class Client extends Component {
    render() {
        return (
            <Provider store={store}>
                <BrowserRouter> 
                    <App/>
                </BrowserRouter>
            </Provider>
        );
    }
}
