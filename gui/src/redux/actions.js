import PropTypes from 'prop-types'

export const UPDATE_PROCESS_ID = "UPDATE_PROCESS_ID"
export const UPDATE_DECODED_FILE = "UPDATE_DECODED_FILE"
export const UPDATE_PROCESSED_FILE = "UPDATE_PROCESSED_FILE"

export const props = {
    'processId': PropTypes.string,
    'decodedFile': PropTypes.string,
    'processedFile': PropTypes.string,
}

export const mapStateToProps = (state) => ({
    'processId': state.processId,
    'decodedFile': state.decodedFile,
    'processedFile': state.processedFile,
});

export function updateDecodedFile(path) {
    return { type: UPDATE_DECODED_FILE, path }
}

export function updateProcessedFile(path) {
    return { type: UPDATE_PROCESSED_FILE, path }
}

export function updateProcessId(id) {
    return { type: UPDATE_PROCESS_ID, id }
}