var m = require("mithril");
var actionTypes = require("./types");

module.exports = {
    Request: function(dispatch, config) {
        dispatch({ type: config.key + "_REQUEST", config: config });

        if (config.options) {
            m.request(config.options)
                .then(function(response) {
                    dispatch({ type: config.key + "_RESPONSE", ok: true, config: config, response: response });

                    return response;
                }, function(response) {
                    dispatch({ type: config.key + "_RESPONSE", ok: false, config: config, response: response });

                    return response;
                }
                ).then(config.handler)
        }
    }
}