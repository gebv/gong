var actionTypes = require("./types");
var m = require("mithril");
var api = require("./Api");
var appConfig = require("../config");

var CLASSIFERS = "classifers"
var ITEMS = "items"

var ApiSearch = function(config) {
    var params = {
        q: config.query
    }
    
    if (config.resource_name == ITEMS) {
        params["classifer_id"] = config.classifer_id;
    }
    
    return appConfig.ApiPrefix + "/resources/"+config.resource_name+"/search?"+m.route.buildQueryString(params);
}

var ApiUpdateOrDelete = function(config) {
    var params = {};
    
    if (config.resource_name == ITEMS) {
        params["classifer_id"] = config.classifer_id;
    }
    
    return appConfig.ApiPrefix + "/resources/"+config.resource_name+"/"+config.id+"?"+m.route.buildQueryString(params);
}

var ApiCreate = function(config) {
    var params = {};
    
    if (config.resource_name == ITEMS) {
        params["classifer_id"] = config.classifer_id;
    }
    
    return appConfig.ApiPrefix + "/resources/"+config.resource_name+"?"+m.route.buildQueryString(params);
}

var Search = function(dispatch, config) {
    var handler = function(response) {
        
        dispatch({
            type: actionTypes.SEARCH_ITEM_SUCCESS,
            items: response,
            resource_name: config.resource_name,
        })
        
        return response;
    }
    
    api.Request(dispatch, {
        key: actionTypes.SEARCH_ITEM, 
        options: { 
            background: true, 
            url: ApiSearch({resource_name: config.resource_name, query: config.query, classifer_id: config.classifer_id}), 
            method: 'GET' 
        }, 
        handler: handler})
}

var New = function(dispatch, config) {
    dispatch({
        type: actionTypes.NEW_ITEM,
        item: config.data
    })
}

var Create = function(dispatch, config) {
    var handler = function(response) {
        
        dispatch({
            type: actionTypes.CREATE_ITEM_SUCCESS,
            item: response,
            resource_name: config.resource_name,
            _TempId: config._TempId,
            classifer_id: config.classifer_id,
        })
        
        return response;
    }
    
    var data = {
        ExtId: config.data.ExtId,
        Articul: config.data.Articul,
        Title: config.data.Title,
        Categories: config.data.Categories,
        Attributes: config.data.Attributes,
        Props: config.data.Props,
        Tags: config.data.Tags,
    }
    
    api.Request(dispatch, {
        key: actionTypes.CREATE_ITEM, 
        options: { 
            background: true, 
            url: ApiCreate({resource_name: config.resource_name, classifer_id: config.classifer_id}),
            data: data, 
            method: 'POST' 
        }, 
        handler: handler})
}

var Update = function(dispatch, config) {
    var handler = function(response) {
        
        dispatch({
            type: actionTypes.UPDATE_ITEM_SUCCESS,
            item: response,
            resource_name: config.resource_name,
            id: config.id
        })
        
        return response;
    }
    
    var data = {
        ExtId: config.data.ExtId,
        Articul: config.data.Articul,
        Title: config.data.Title,
        Categories: config.data.Categories,
        Attributes: config.data.Attributes,
        Props: config.data.Props,
        Tags: config.data.Tags,
    }
    
    api.Request(dispatch, {
        key: actionTypes.UPDATE_ITEM, 
        options: { 
            background: true, 
            url: ApiUpdateOrDelete({id: config.id, resource_name: config.resource_name}),
            data: data, 
            method: 'PUT' 
        }, 
        handler: handler})
}

var Delete = function(dispatch, config) {
    var handler = function(response) {
        
        dispatch({
            type: actionTypes.DELETE_ITEM_SUCCESS,
            item: response,
            resource_name: config.resource_name,
        })
        
        return response;
    }
    
    api.Request(dispatch, {
        key: actionTypes.DELETE_ITEM, 
        options: { 
            background: true, 
            url: ApiUpdateOrDelete({id: config.id, resource_name: config.resource_name}), 
            method: 'DELETE' 
        }, 
        handler: handler})
}

var DeleteTemp = function(dispatch, config) {
    dispatch({
        type: actionTypes.DELETE_TEMP_ITEM,
        _TempId: config._TempId,
        id: config.id,
        resource_name: config.resource_name,
    })
}

module.exports = {
    ITEMS: ITEMS,
    CLASSIFERS: CLASSIFERS,
    
    ApiCreate: ApiCreate,
    ApiUpdate: ApiUpdateOrDelete,
    ApiSearch: ApiSearch,
    
    Search: Search,
    Create: Create,
    Update: Update,
    Delete: Delete,
    DeleteTemp: DeleteTemp, 
    New: New,
}
