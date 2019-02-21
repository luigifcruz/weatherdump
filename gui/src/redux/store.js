import { createStore, applyMiddleware } from 'redux'
import Reducer from './reducer'
import createLogger from 'redux-logger'

let defaultState = {
    "appId": null
}

const middleware = [ createLogger ]

function configureStore(initialState = defaultState) {
    const store = createStore(Reducer, initialState, applyMiddleware(...middleware));
    return store;
}

export { configureStore }