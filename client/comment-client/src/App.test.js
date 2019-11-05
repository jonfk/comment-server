import React from 'react';
import ReactDOM from 'react-dom';
import App from './containers/App';

import { createStore } from 'redux';
import commentApp from './reducers';
import { viewTypes, switchView, addComment, replyComment } from './actions';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<App />, div);
});

it('store dispatch testing', () => {
    let store = createStore(commentApp);

    console.log(JSON.stringify(store.getState()));

    store.dispatch(addComment({id: '1', text: 'first', author: 'Amy'}));
    store.dispatch(addComment({id: '2', text: 'secundo', author: 'Jfk'}));
    store.dispatch(addComment({id: '3', text: 'secundo2', author: 'Jfk'}));
    store.dispatch(addComment({id: '4', text: 'secundo3', author: 'Jfk'}));
    store.dispatch(replyComment(2, {id: '5', text: 'fivos', author: 'Amy'}));
    store.dispatch(replyComment(5, {id: '6', text: 'sexos', author: 'Jfk'}));
    store.dispatch(switchView(viewTypes.ACCOUNT_VIEW));

    console.log(JSON.stringify(store.getState()));
});
