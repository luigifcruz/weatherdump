import {
    UPDATE_PROCESS_ID,
    UPDATE_DECODED_FILE,
    UPDATE_WORKING_PATH,
    UPDATE_PROCESS_DESCRIPTOR,
    UPDATE_MANIFEST,
    UPDATE_DEMOD_FILE,
    TOGGLE_PARSER,
    TOGGLE_COMPOSER,
    TOGGLE_ENHANCEMENT
} from "./actions"

export default function reducer(state, action) {
    switch (action.type) {
        case UPDATE_PROCESS_ID:
        return Object.assign({}, state, {
            processId: action.id
        })
        case UPDATE_DECODED_FILE:
        return Object.assign({}, state, {
            decodedFile: action.path
        })
        case UPDATE_WORKING_PATH:
        return Object.assign({}, state, {
            workingPath: action.path
        })
        case UPDATE_PROCESS_DESCRIPTOR:
        return Object.assign({}, state, {
            processDescriptor: action.descriptor
        })
        case UPDATE_MANIFEST:
        return Object.assign({}, state, {
            manifestParser: action.parser,
            manifestComposer: action.composer
        })
        case UPDATE_DEMOD_FILE:
        return Object.assign({}, state, {
            demodulatedFile: action.file
        })
        case TOGGLE_PARSER:
        return Object.assign({}, state, {
            manifestParser: {
                ...state.manifestParser,
                [action.apid]: {
                    ...state.manifestParser[action.apid],
                    Activated: !state.manifestParser[action.apid].Activated
                }
            }
        })
        case TOGGLE_COMPOSER:
        return Object.assign({}, state, {
            manifestComposer: {
                ...state.manifestComposer,
                [action.apid]: {
                    ...state.manifestComposer[action.apid],
                    Activated: !state.manifestComposer[action.apid].Activated
                }
            }
        })
        case TOGGLE_ENHANCEMENT:
        return Object.assign({}, state, {
            processorEnhancements: {
                ...state.processorEnhancements,
                [action.key]: {
                    ...state.processorEnhancements[action.key],
                    Activated: !state.processorEnhancements[action.key].Activated
                }
            }
        })
        default:
        return state;
    }
}