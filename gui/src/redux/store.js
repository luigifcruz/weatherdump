import { applyMiddleware, createStore } from 'redux'

import Reducer from './reducer'
import { createLogger } from 'redux-logger'

let defaultState = {
    'processId': null,
    'processDescriptor': null,
    'manifestParser': {},
    'manifestComposer': {},
    'decodedFile': null,
    'processorEnhancements': {
        "Invert": {
            "Name": "Invert Infrared Pixels",
            "Activated": true
        },
        "Flop": {
            "Name": "Horizontally Flip Image",
            "Activated": false
        },
        "Equalize": {
            "Name": "Histogram Equalization",
            "Activated": true
        },
        "ExportPNG": {
            "Name": "Lossless PNG",
            "Activated": false
        },
        "ExportJPEG": {
            "Name": "Lossless JPEG",
            "Activated": true
        }
    },
    'demodulatedFile': null,
    'workingPath': null
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