import React from 'react';
import ReactDOM from 'react-dom';
import AppC from './containers/App';
import './index.css';
import { viewTypes } from './actions';
import commentApp from './reducers';
import { Provider } from 'react-redux';
import { createStore } from 'redux'

let store = createStore(commentApp);

ReactDOM.render(
    <Provider store={store}>
        <AppC view={viewTypes.ACCOUNT_VIEW}/>
    </Provider>,
  document.getElementById('root')
);
