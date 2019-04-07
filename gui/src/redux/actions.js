import PropTypes from 'prop-types'

export const UPDATE_PROCESS_ID = "UPDATE_PROCESS_ID"
export const UPDATE_DECODED_FILE = "UPDATE_DECODED_FILE"
export const UPDATE_WORKING_PATH = "UPDATE_WORKING_PATH"
export const UPDATE_PROCESS_DESCRIPTOR = "UPDATE_PROCESS_DESCRIPTOR"
export const UPDATE_MANIFEST = "UPDATE_MANIFEST"
export const UPDATE_DEMOD_FILE = "UPDATE_DEMOD_FILE"
export const TOGGLE_PARSER = "TOGGLE_PARSER"
export const TOGGLE_COMPOSER = "TOGGLE_COMPOSER"
export const TOGGLE_ENHANCEMENT = "TOGGLE_ENHANCEMENT"

export const props = {
    'processId': PropTypes.string,
    'processDescriptor': PropTypes.string,
    'processorEnhancements': PropTypes.object,
    'manifestParser': PropTypes.object,
    'manifestComposer': PropTypes.object,
    'decodedFile': PropTypes.string,
    'demodulatedFile': PropTypes.string,
    'workingPath': PropTypes.string,
}

export const mapStateToProps = (state) => ({
    'processId': state.processId,
    'processDescriptor': state.processDescriptor,
    'processorEnhancements': state.processorEnhancements,
    'manifestParser': state.manifestParser,
    'manifestComposer': state.manifestComposer,
    'decodedFile': state.decodedFile,
    'demodulatedFile': state.demodulatedFile,
    'workingPath': state.processedFolder,
});

export function updateDecodedFile(path) {
    return { type: UPDATE_DECODED_FILE, path }
}

export function updateWorkingPath(path) {
    return { type: UPDATE_WORKING_PATH, path }
}

export function updateProcessId(id) {
    return { type: UPDATE_PROCESS_ID, id }
}

export function updateProcessDescriptor(descriptor) {
    return { type: UPDATE_PROCESS_DESCRIPTOR, descriptor }
}

export function updateManifest(parser, composer) {
    return { type: UPDATE_MANIFEST, parser, composer }
}

export function toggleParserActivation(apid) {
    return { type: TOGGLE_PARSER, apid }
}

export function toggleComposerActivation(apid) {
    return { type: TOGGLE_COMPOSER, apid }
}

export function toggleEnhancement(key) {
    return { type: TOGGLE_ENHANCEMENT, key }
}

export function updateDemodulatedFile(file) {
    return { type: UPDATE_DEMOD_FILE, file }
}