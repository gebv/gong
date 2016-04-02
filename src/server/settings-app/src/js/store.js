var m = require('mithril');
var redux = require("redux");
var app = require("./reducers/app");
var createLogger = require('redux-logger');
var thunk = require('redux-thunk');

var logger = createLogger({duration: true});

var store = redux.createStore(app, redux.applyMiddleware(logger));

store.subscribe(m.redraw.bind(m))

module.exports = store;