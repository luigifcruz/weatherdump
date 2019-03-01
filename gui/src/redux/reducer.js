import {
    UPDATE_PROCESS_ID,
    UPDATE_PROCESS_DATALINK
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
        default:
        return state;
    }
}