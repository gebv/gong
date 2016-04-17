var redux = require("redux");
var _ = require("lodash");
var actionTypes = require("../actions/types");

var appInitialState = {
    items: [],
    isWaitingRequest: false,
    hasNext: false,
    nextPage: 0,
    total: 0
}

var AppConfig = function(state, action) {
    
    state = state || appInitialState

    switch (action.type) {
        case actionTypes.GET_ITEM_REQUEST:
            state.items = [];
            break;
        case actionTypes.DELETE_TEMP_ITEM:
            state.items = _.filter(state.items, function(item){
                return item._TempId != action._TempId;
            });
            break;
        case actionTypes.NEW_ITEM:
            state.items = [_.extend({}, action.item)].concat(state.items);
            
            break;
        case actionTypes.CREATE_ITEM_SUCCESS:
            state.items = _.map(state.items, function(item){
                if (action._TempId == item._TempId) {
                    return _.extend({}, action.item);
                }
                
                return item
            })
            break;
        case actionTypes.UPDATE_ITEM_SUCCESS:
            state.items = _.map(state.items, function(item){
                
                if (action.id == item.ItemId) {
                    
                    return action.item;
                }
                
                return item
            })
            
            break;
        case actionTypes.DELETE_ITEM_SUCCESS:
            state.items = _.filter(state.items, function(item){
                return item.ItemId != action.item.ItemId;
            });
            break;
        case actionTypes.SEARCH_ITEM_SUCCESS:
            // Если первая страница, обновляем список
            // В противном случае дозагружаем
            
            state.hasNext = action.response.Data.HasNext
            state.nextPage = action.response.Data.NextPage
            state.total = action.response.Data.Total
            
            if (action.response && action.response.Code == "success") {
                if (action.response.Data.NextPage == 0) {
                    state.items = action.response.Data.Items;    
                } else {
                    _.each(action.response.Data.Items, function(item){
                        state.items.push(item)
                    })
                }
                    
            }
            
            break;  
        
        case actionTypes.CREATE_ITEM_REQUEST:
        case actionTypes.UPDATE_ITEM_REQUEST:
        case actionTypes.DELETE_ITEM_REQUEST:
        case actionTypes.SEARCH_ITEM_REQUEST:
            state.isWaitingRequest = true;
            
            break;
        case actionTypes.CREATE_ITEM_RESPONSE:
        case actionTypes.UPDATE_ITEM_RESPONSE:
        case actionTypes.DELETE_ITEM_RESPONSE:
        case actionTypes.SEARCH_ITEM_RESPONSE:
            state.isWaitingRequest = false;
            
            break;
    }
    
    return state;
};


var App = redux.combineReducers({
    app: AppConfig,
    lastAction: function(state, action) {
        return action;
    }
})

module.exports = App;