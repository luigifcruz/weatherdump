import { createStore, applyMiddleware } from 'redux'
import Reducer from './reducer'
import createLogger from 'redux-logger'
import Ring from "ringjs"

let defaultState = {
    'history': new Ring(200),
    'mapdata': new Ring(100),
    'state': {
    	'gnss_fix_quality': 0,
    	'gnss_fix_type': 0,
    	'sat_count_gps': 0,
    	'sat_count_glonass': 0,
    	'longitude': 0,
    	'latitude': 0,
    	'true_north': 0,
    	'ground_speed': 0,
    	'altitude': 0,
        'gnss_sat_info': null
    },
    'settings': {
    	'led_status': 0,
    	'serial_tx_active': 0
    }
}

const middleware = [ ]

function configureStore(initialState = defaultState) {
    const store = createStore(Reducer, initialState, applyMiddleware(...middleware));
    return store;
}

export { configureStore }