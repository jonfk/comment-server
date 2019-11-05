import { createStore } from 'redux';
import commentApp from './reducers';

let store = createStore(commentApp);
