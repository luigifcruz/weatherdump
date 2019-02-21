import {
    UPDATE_REGISTRY
} from "./actions"

export default function reducer(state, action) {
    switch (action.type) {
        case UPDATE_REGISTRY:
        return Object.assign({}, state, {
            appId: action.uuid
        })
        default:
        return state;
    }
}