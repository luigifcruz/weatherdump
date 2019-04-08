import { AppContainer } from 'react-hot-loader';
import Client from './client';
import React from 'react';
import ReactDOM from 'react-dom';

const render = (Component) => {
    ReactDOM.render(
        <AppContainer>
            <Component />
        </AppContainer>,
        document.getElementById('root'),
    );
};
 
render(Client);

if (module && module.hot) {
    module.hot.accept('./client', () => {
        const HotApp = require('./client').default;
        render(HotApp);
    });
}