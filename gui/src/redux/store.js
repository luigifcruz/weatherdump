import { createStore, applyMiddleware } from 'redux'
import Reducer from './reducer'
import { createLogger } from 'redux-logger'

let defaultState = {
    'processId': null,
    'processDatalink': null,
    'processDescriptor': null,
    'manifestParser': {},
    'manifestComposer': {},
    'decodedFile': null,
    'demodulatedFile': null,
    'workingFolder': null,
}

let middleware = new Array()

if (process.env.NODE_ENV == 'development') {
    middleware.push(createLogger())
}

function configureStore(initialState = defaultState) {
    const store = createStore(Reducer, initialState, applyMiddleware(...middleware));
    return store;
}

export { configureStore }