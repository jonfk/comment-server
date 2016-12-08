
import { SWITCH_VIEW, ADD_COMMENT, REPLY_COMMENT, viewTypes } from '../actions';

const initialState = {
    view: viewTypes.COMMENT_VIEW,
    comments: {
        roots: [],
        parentsToChildren: {}
    },
    authenticationToken: ""
};

export default function commentApp(state = initialState, action) {
    return {
        view: view(state.view, action),
        comments: comments(state.comments, action)
    };
}

function view(state = viewTypes.COMMENT_VIEW, action) {
    switch (action.type) {
    case SWITCH_VIEW:
        return action.view;
    default:
        return state;
    }
}

function comments(state = initialState.comments, action) {
    switch (action.type) {
    case ADD_COMMENT:
        return Object.assign({}, state, { roots: [...state.roots, {
            text: action.comment.text,
            author: action.comment.author
        }]});
    case REPLY_COMMENT:
        return Object.assign({}, state, { parentsToChildren: {...state.parentsToChildren,  [action.parentId]: {
            text: action.comment.text,
            author: action.comment.author
        }}});
    default:
        return state;
    }
}
