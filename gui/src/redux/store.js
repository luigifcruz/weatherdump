import { createStore, applyMiddleware } from 'redux'
import Reducer from './reducer'
import createLogger from 'redux-logger'

let defaultState = {
    'processId': null,
    'processDatalink': null,
    'decodedFile': null,
    'processedFile': null,
}

const middleware = [ createLogger ]

function configureStore(initialState = defaultState) {
    const store = createStore(Reducer, initialState, applyMiddleware(...middleware));
    return store;
}

export { configureStore }