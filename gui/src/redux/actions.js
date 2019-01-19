import PropTypes from 'prop-types'
import React from 'react'

export const UPDATE_STATE = "UPDATE_STATE"
export const UPDATE_SETTINGS = "UPDATE_SETTINGS"
export const UPDATE_HISTORY = "UPDATE_HISTORY"
export const UPDATE_MAPDATA = "UPDATE_MAPDATA"

export const props = {
    'state': PropTypes.object.isRequired,
    'settings': PropTypes.object.isRequired,
    'history': PropTypes.object.isRequired,
    'mapdata': PropTypes.object.isRequired
}

export const mapStateToProps = (state) => ({
    'state': state.state,
    'settings': state.settings,
    'history': state.history,
    'mapdata': state.mapdata
});

export function updateState(key, value) {
    return { type: UPDATE_STATE, key, value }
}

export function updateSettings(key, value) {
    return { type: UPDATE_SETTINGS, key, value }
}

export function updateHistory(value) {
    return { type: UPDATE_HISTORY, value }
}

export function updateMapdata(value) {
    return { type: UPDATE_MAPDATA, value }
}