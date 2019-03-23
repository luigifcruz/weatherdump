import PropTypes from 'prop-types'

export const UPDATE_PROCESS_ID = "UPDATE_PROCESS_ID"
export const UPDATE_DECODED_FILE = "UPDATE_DECODED_FILE"
export const UPDATE_WORKING_FOLDER = "UPDATE_WORKING_FOLDER"
export const UPDATE_PROCESS_DATALINK = "UPDATE_PROCESS_DATALINK"
export const UPDATE_PROCESS_DESCRIPTOR = "UPDATE_PROCESS_DESCRIPTOR"
export const UPDATE_MANIFEST = "UPDATE_MANIFEST"
export const UPDATE_DEMOD_FILE = "UPDATE_DEMOD_FILE"

export const props = {
    'processId': PropTypes.string,
    'processDatalink': PropTypes.string,
    'processDescriptor': PropTypes.string,
    'manifestParser': PropTypes.object,
    'manifestComposer': PropTypes.object,
    'decodedFile': PropTypes.string,
    'demodulatedFile': PropTypes.string,
    'workingFolder': PropTypes.string,
}

export const mapStateToProps = (state) => ({
    'processId': state.processId,
    'processDatalink': state.processDatalink,
    'processDescriptor': state.processDescriptor,
    'manifestParser': state.manifestParser,
    'manifestComposer': state.manifestComposer,
    'decodedFile': state.decodedFile,
    'demodulatedFile': state.demodulatedFile,
    'workingFolder': state.processedFolder,
});

export function updateDecodedFile(path) {
    return { type: UPDATE_DECODED_FILE, path }
}

export function updateWorkingFolder(path) {
    return { type: UPDATE_WORKING_FOLDER, path }
}

export function updateProcessId(id) {
    return { type: UPDATE_PROCESS_ID, id }
}

export function updateProcessDatalink(datalink) {
    return { type: UPDATE_PROCESS_DATALINK, datalink }
}

export function updateProcessDescriptor(descriptor) {
    return { type: UPDATE_PROCESS_DESCRIPTOR, descriptor }
}

export function updateManifest(parser, composer) {
    return { type: UPDATE_MANIFEST, parser, composer }
}

export function updateDemodulatedFile(file) {
    return { type: UPDATE_DEMOD_FILE, file }
}