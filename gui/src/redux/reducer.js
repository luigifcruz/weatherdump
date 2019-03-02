import {
    UPDATE_PROCESS_ID,
    UPDATE_PROCESS_DATALINK,
    UPDATE_DECODED_FILE,
    UPDATE_PROCESSED_FILE
} from "./actions"

export default function reducer(state, action) {
    switch (action.type) {
        case UPDATE_PROCESS_ID:
        return Object.assign({}, state, {
            processId: action.id
        })
        case UPDATE_PROCESS_DATALINK:
        return Object.assign({}, state, {
            processDatalink: action.datalink
        })
        case UPDATE_DECODED_FILE:
        return Object.assign({}, state, {
            decodedFile: action.path
        })
        case UPDATE_PROCESSED_FILE:
        return Object.assign({}, state, {
            processedFile: action.path
        })
        default:
        return state;
    }
}