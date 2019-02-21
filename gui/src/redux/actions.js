import PropTypes from 'prop-types'

export const UPDATE_REGISTRY = "UPDATE_REGISTRY"

export const props = {
    'appId': PropTypes.string
}

export const mapStateToProps = (state) => ({
    'appId': state.appId
});

export function updateRegistry(uuid) {
    return { type: UPDATE_REGISTRY, uuid }
}