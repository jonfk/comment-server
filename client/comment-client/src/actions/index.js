
export const SWITCH_VIEW = "SWITCH_VIEW";
export const ADD_COMMENT        = "ADD_COMMENT";
export const REPLY_COMMENT = "REPLY_COMMENT";
export const SIGN_IN            = "SIGN_IN";
export const CREATE_ACCOUNT     = "CREATE_ACCOUNT";

export const viewTypes = {
    ACCOUNT_VIEW: "ACCOUNT_VIEW",
    COMMENT_VIEW: "COMMENT_VIEW"
};

export function switchView(view) {
    return { type: SWITCH_VIEW, view: view };
}

export function addComment(comment) {
    return { type: ADD_COMMENT, comment: {
        id: comment.id,
        text: comment.text,
        author: comment.author
    }};
}

export function replyComment(parentId, comment) {
    return { type: REPLY_COMMENT, parentId: parentId, comment: {
        id: comment.id,
        text: comment.text,
        author: comment.author
    }};
}

export function signIn(token) {
    return { type: SIGN_IN, token: token };
}

export function createAccount(token) {
    return { type: CREATE_ACCOUNT, token: token};
}
