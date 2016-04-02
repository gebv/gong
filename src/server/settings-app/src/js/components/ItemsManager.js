var m = require("mithril");
var ItemActions = require("../actions/Items");

var store = require("../store");
var _ = require("lodash");

var ace = require('brace');
require('brace/mode/toml');
require('brace/mode/html');
// require('brace/theme/monokai');

var SearchFilter = {
    controller: function(c) {
        var api = {
            query: m.prop(""),
            onSearch: function() {
                return function(e) {
                    e.preventDefault();

                    ItemActions.Search(store.dispatch, { query: this.query(), resource_name: c.resource_name, classifer_id: c.classifer_id })

                    return false;
                }.bind(this);
            }
        }

        ItemActions.Search(store.dispatch, { query: "", resource_name: c.resource_name, classifer_id: c.classifer_id })

        return api;
    },
    view: function(c, args) {
        return <div class="uk-panel uk-panel-box">
            <form class="uk-form" onsubmit={c.onSearch() }>
                <input placeholder="Search..." class="uk-form-large uk-width-1-1" oninput={m.withAttr("value", c.query) } value={c.query() } />
                <p class="uk-text-muted uk-float-right"></p>
            </form>
        </div>
    }
}

var mapViews = {
    "classifers": {
        controller: function(c) {
            var api = {
                onEdit: function(mode_name) {
                    return function(e) {
                        e.preventDefault();

                        c.mode(mode_name);

                        return false;
                    }
                }
            }

            api = _.extend({}, api, c);

            return api;
        },
        view: function(c, args) {
            var url = "/@settings/classifers/" + args.data.ItemId + "/items";
            var isShowSubPositions = c.resource_name == ItemActions.CLASSIFERS;
            var buttonPositions = <a class="uk-button uk-button-primary" href={url} config={m.route}>Позиции</a>

            return <div class="uk-panel uk-panel-box">
                <div class="uk-panel-badge uk-text-small uk-text-muted">{args.data.ExtId}</div>
                <p class="uk-panel-title"><a onclick={c.onEdit("edit") } class="uk-text-success uk-icon-edit" href=""></a> {args.data.Title}</p>
                <p class="uk-article-meta">{args.data.Category} {args.data.Tags}</p>
                <div>
                    {isShowSubPositions ? buttonPositions : ""}
                </div>
            </div>
        }
    },
    "items": {
        controller: function(c) {
            var api = {
                onEdit: function(mode_name) {
                    return function(e) {
                        e.preventDefault();

                        c.mode(mode_name);

                        return false;
                    }
                }
            }

            api = _.extend({}, api, c);

            return api;
        },
        view: function(c, args) {
            var url = "/@settings/classifers/" + args.data.ItemId + "/items";
            var isShowSubPositions = c.resource_name == ItemActions.CLASSIFERS;
            var buttonPositions = <a class="uk-button uk-button-primary" href={url} config={m.route}>Позиции</a>

            var listRefWidgets = [];

            if (args.data.Props._BuildTraceWidgets) {
                for (var i = 0; i < args.data.Props._BuildTraceWidgets.length; i += 2) {
                    var url = "/@settings/classifers/"+args.data.Props._BuildTraceWidgets[i]+"/items/edit?special_id="+args.data.Props._BuildTraceWidgets[i+1];
                    var label = [args.data.Props._BuildTraceWidgets[i], args.data.Props._BuildTraceWidgets[i + 1]].join(" ");
                    listRefWidgets.push(m("li", m("a", {href: url, config: m.route}, label)))
                }
            }

            return <div class="uk-panel uk-panel-box">
                <div class="uk-panel-badge uk-text-small uk-text-muted">{args.data.ExtId}</div>
                <p class="uk-panel-title"><a onclick={c.onEdit("edit") } class="uk-text-success uk-icon-edit" href=""></a> {args.data.Title}</p>
                <p class="uk-article-meta">{args.data.Category} {args.data.Tags}</p>
                <div>
                    {isShowSubPositions ? buttonPositions : ""}
                </div>
                <p>related widgets:</p>
                <ul class="uk-list">{listRefWidgets}</ul>
            </div>
        }
    }
}

var mapEditor = {
    "classifers": {
        controller: function(c) {

            // TODO: Что бы знаения обновились, переключать режим

            var api = {
                resource_name: c.resource_name,
                classifer_id: c.classifer_id,

                _TempId: m.prop(c.data._TempId || ""),
                ItemId: m.prop(c.data.ItemId || ""),
                ExtId: m.prop(c.data.ExtId || ""),
                Title: m.prop(c.data.Title || ""),

                onDelete: function() {
                    return function(e) {
                        e.preventDefault();

                        if (!confirm("Вы уверены в удалении?")) {
                            return false;
                        }

                        if (this._TempId()) {
                            ItemActions.DeleteTemp(store.dispatch, { _TempId: this._TempId(), id: this.ItemId(), resource_name: this.resource_name });
                        } else {
                            ItemActions.Delete(store.dispatch, { id: this.ItemId(), resource_name: this.resource_name });
                        }

                        return false;
                    }.bind(this);
                },
                onCancel: function() {
                    return function(e) {
                        c.mode("view");
                        return false
                    }
                },
                onCreate: function() {
                    return function(e) {
                        e.preventDefault();

                        var data = {
                            ExtId: this.ExtId(),
                            Title: this.Title(),
                        };

                        c.mode("view");

                        if (this._TempId()) {
                            ItemActions.Create(store.dispatch, { _TempId: this._TempId(), data: data, resource_name: this.resource_name, classifer_id: this.classifer_id });
                        } else {
                            ItemActions.Update(store.dispatch, { id: this.ItemId(), data: data, resource_name: this.resource_name });
                        }


                        return false;
                    }.bind(this);
                }
            }

            api = _.extend({}, api, c);

            return api;
        },
        view: function(c, args) {

            return <div class="uk-panel uk-panel-box">
                <form class="uk-form uk-form-horizontal" onsubmit={c.onCreate() }>
                    <div class="uk-form-row">
                        <label class="uk-form-label">Title</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" oninput={m.withAttr("value", c.Title) } value={c.Title() }/>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <label class="uk-form-label">Name</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" oninput={m.withAttr("value", c.ExtId) } value={c.ExtId() }/>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <button class="uk-button" onclick={c.onCancel() }>Отменить</button>
                        <button class="uk-button uk-button-success">Сохранить</button>

                        <a onclick={c.onDelete() } class="uk-button uk-float-right uk-button-danger"  href="">Удалить</a>
                    </div>
                </form>
            </div>
        }
    },
    "items": {
        controller: function(c) {

            // TODO: Что бы знаения обновились, переключать режим

            var api = {
                InitConfigEditor: function(value) {
                    return function(e) {
                        var editor = ace.edit([c.data.ItemId, c.data._TempId, "config"].join(""));
                        editor.getSession().setMode('ace/mode/toml');
                        // editor.setTheme('ace/theme/monokai');
                        editor.setValue(value());
                        editor.getSession().on('change', function() {
                            value(editor.getSession().getValue());
                        });
                    }
                },

                InitContentEditor: function(value) {
                    return function(e) {
                        var editor = ace.edit([c.data.ItemId, c.data._TempId, "content"].join(""));
                        editor.getSession().setMode('ace/mode/html');
                        // editor.setTheme('ace/theme/monokai');
                        editor.setValue(value());
                        editor.getSession().on('change', function() {
                            value(editor.getSession().getValue());
                        });
                    }
                },
                resource_name: c.resource_name,
                classifer_id: c.classifer_id,

                _TempId: m.prop(c.data._TempId || ""),
                ItemId: m.prop(c.data.ItemId || ""),
                ExtId: m.prop(c.data.ExtId || ""),
                Title: m.prop(c.data.Title || ""),
                Props: m.prop(c.data.Props || {}),

                Config: m.prop(c.data.Props && c.data.Props.Config ? c.data.Props.Config : ""),
                Content: m.prop(c.data.Props && c.data.Props.Content ? c.data.Props.Content : ""),

                onDelete: function() {
                    return function(e) {
                        e.preventDefault();

                        if (!confirm("Вы уверены в удалении?")) {
                            return false;
                        }

                        if (this._TempId()) {
                            ItemActions.DeleteTemp(store.dispatch, { _TempId: this._TempId(), id: this.ItemId(), resource_name: this.resource_name });
                        } else {
                            ItemActions.Delete(store.dispatch, { id: this.ItemId(), resource_name: this.resource_name });
                        }

                        return false;
                    }.bind(this);
                },
                onCancel: function() {
                    return function(e) {
                        c.mode("view");
                        return false
                    }
                },
                onCreate: function() {
                    return function(e) {
                        e.preventDefault();

                        var data = {
                            ExtId: this.ExtId(),
                            Title: this.Title(),
                            Props: { Config: this.Config(), Content: this.Content() },
                        };

                        c.mode("view");

                        if (this._TempId()) {
                            ItemActions.Create(store.dispatch, { _TempId: this._TempId(), data: data, resource_name: this.resource_name, classifer_id: this.classifer_id });
                        } else {
                            ItemActions.Update(store.dispatch, { id: this.ItemId(), data: data, resource_name: this.resource_name });
                        }


                        return false;
                    }.bind(this);
                }
            }

            api = _.extend({}, api, c);

            return api;
        },
        view: function(c, args) {

            return <div class="uk-panel uk-panel-box">
                <form class="uk-form uk-form-horizontal" onsubmit={c.onCreate() }>
                    <div class="uk-form-row">
                        <label class="uk-form-label">Title</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" oninput={m.withAttr("value", c.Title) } value={c.Title() }/>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <label class="uk-form-label">Name</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" oninput={m.withAttr("value", c.ExtId) } value={c.ExtId() }/>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <label class="uk-form-label">Config</label>

                        <div class="uk-form-controls">
                            <div id={[c.data.ItemId, c.data._TempId, "config"].join("") } config={c.InitConfigEditor(c.Config) }  style="height:200px;"></div>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <label class="uk-form-label">Content</label>

                        <div class="uk-form-controls">
                            <div id={[c.data.ItemId, c.data._TempId, "content"].join("") } config={c.InitContentEditor(c.Content) } style="height:200px;"></div>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <button class="uk-button" onclick={c.onCancel() }>Отменить</button>
                        <button class="uk-button uk-button-success">Сохранить</button>

                        <a onclick={c.onDelete() } class="uk-button uk-float-right uk-button-danger"  href="">Удалить</a>
                    </div>
                </form>
            </div>
        }
    }
}

var SearchResult = {
    controller: function(c) {
        // resource_name

        var api = {
            onCreateTempItem: function() {
                return function(e) {
                    e.preventDefault();

                    ItemActions.New(store.dispatch, { data: { _TempId: _.now() } });

                    return false;
                }
            }
        }

        return api;
    },
    view: function(c, args) {
        var list = m.component(ListItems, args);
        var createButton = m("a.uk-button uk-width-1-1 uk-margin-top", { onclick: c.onCreateTempItem() }, "+");

        return <div>
            {createButton}
            {list}
        </div>
    }
}

var ListItems = {
    controller: function(c) {
        var api = {

        }

        return api;
    },
    view: function(c, args) {
        var items = _.map(store.getState().app.items, function(item) {
            var itemList = m.component(ItemList, { data: item, resource_name: args.resource_name, classifer_id: args.classifer_id });

            return m("li", { key: "" + item._TempId + item.ItemId }, itemList)
        });

        var waiting = store.getState().app.isWaitingRequest ? m("i.uk-icon-spinner uk-icon-spin") : false;

        return <ul class="uk-list uk-list-space">{waiting || items}</ul>
    }
}


var ItemList = {
    controller: function(c) {

        var api = {
            resource_name: c.resource_name,
            mode: m.prop("view"),
        }

        return api;
    },
    view: function(c, args) {

        if (args.data._TempId) {
            c.mode("edit");
        }
        
        var special_id = m.route.param("special_id");
        
        if (special_id && (special_id == args.data.ItemId || special_id == args.data.ExtId)) {
            c.mode("edit");
        }

        switch (c.mode()) {
            case "edit":
                return m.component(mapEditor[c.resource_name], _.extend({}, args, { mode: c.mode, classifer_id: args.classifer_id }));
            case "view":
                return m.component(mapViews[c.resource_name], _.extend({}, args, { mode: c.mode, classifer_id: args.classifer_id }));
            default:
                return m("li", args.data.ItemId)
        }
    }
}

var Manager = {
    controller: function(c) {
        var api = {

        }

        return api;
    },
    view: function(c, args) {
        var classifer_id = m.route.param("classifer_id");

        var filter = m.component(SearchFilter, { resource_name: args.resource_name, classifer_id: classifer_id })
        var result = m.component(SearchResult, { resource_name: args.resource_name, classifer_id: classifer_id })

        return m("div.uk-margin-top", [filter, result]);
    }
}

module.exports = {
    Manager: Manager
}