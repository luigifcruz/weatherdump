import PropTypes from 'prop-types'

export const UPDATE_PROCESS_ID = "UPDATE_PROCESS_ID"
export const UPDATE_DECODED_FILE = "UPDATE_DECODED_FILE"
export const UPDATE_PROCESSED_FOLDER = "UPDATE_PROCESSED_FOLDER"
export const UPDATE_PROCESS_DATALINK = "UPDATE_PROCESS_DATALINK"

export const props = {
    'processId': PropTypes.string,
    'processDatalink': PropTypes.string,
    'decodedFile': PropTypes.string,
    'processedFolder': PropTypes.string,
}

export const mapStateToProps = (state) => ({
    'processId': state.processId,
    'processDatalink': state.processDatalink,
    'decodedFile': state.decodedFile,
    'processedFolder': state.processedFolder,
});

export function updateDecodedFile(path) {
    return { type: UPDATE_DECODED_FILE, path }
}

export function updateProcessedFolder(path) {
    return { type: UPDATE_PROCESSED_FOLDER, path }
}

export function updateProcessId(id) {
    return { type: UPDATE_PROCESS_ID, id }
}

export function updateProcessDatalink(datalink) {
    return { type: UPDATE_PROCESS_DATALINK, datalink }
}