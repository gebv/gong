var m = require("mithril");
var ItemActions = require("../actions/Items");
var ActionTypes = require("../actions/types");


var store = require("../store");
var _ = require("lodash");

var ace = require('brace');
require('brace/mode/toml');
require('brace/mode/html');

var SearchFilter = {
    controller: function(c) {
        var api = {
            query: m.prop(""),
            onSearch: function() {
                return function(e) {
                    e.preventDefault();

                    ItemActions.Search(store.dispatch, { 
                        query: this.query(), 
                        resource_name: c.resource_name, 
                        bucket_id: c.bucket_id, 
                        page: 0
                    })

                    return false;
                }.bind(this);
            },
            onLoadMore: function() {
                return function(e) {
                    e.preventDefault();

                    ItemActions.Search(store.dispatch, { 
                        query: this.query(), 
                        resource_name: c.resource_name, 
                        bucket_id: c.bucket_id, 
                        page: store.getState().app.nextPage 
                    })

                    return false;
                }.bind(this);
            }
        }

        ItemActions.Search(store.dispatch, { query: "", resource_name: c.resource_name, bucket_id: c.bucket_id, page: 0})

        return api;
    },
    view: function(c, args) {
        var searchResult = m.component(SearchResult, {resource_name: args.resource_name, bucket_id: args.bucket_id, onLoadMore: c.onLoadMore()});
        
        return <div class="uk-panel">
            <form class="uk-form uk-width-1-1" onsubmit={c.onSearch() }>
                <input placeholder="Search..." class="uk-form-large uk-width-1-1" oninput={m.withAttr("value", c.query) } value={c.query() } />
                <p class="uk-text-muted uk-float-right"></p>
            </form>
            {searchResult}
        </div>
    }
}

var SearchResult = {
    controller: function(c) {
        // resource_name

        var api = {
        }

        return api;
    },
    view: function(c, args) {
        var list = m.component(ListItems, args);
        
        var url = args.resource_name == ItemActions.ITEMS? 
            "/buckets/"+args.bucket_id+"/files/new":
            "/buckets/new";
        
        var createButton = m("a.uk-button uk-width-1-1 uk-margin-top", { href: url, config: m.route}, "+");

        return <div>
            {createButton}
            {list}
        </div>
    }
}

var ListItems = {
    controller: function(c) {
        var api = {
            onLoadMore: c.onLoadMore,
        }

        return api;
    },
    view: function(c, args) {
        var items = _.map(store.getState().app.items, function(item) {
            var itemList = m.component(ItemList, { data: item, resource_name: args.resource_name, bucket_id: args.bucket_id });

            return m("li", { key: "" + item._TempId + item.Id }, itemList)
        });
        
        if (store.getState().app.hasNext) {
            items.push(m("li", m("a.uk-button uk-width-1-1", {onclick: c.onLoadMore}, "more")))    
        }

        var waiting = store.getState().app.isWaitingRequest ? m("i.uk-icon-spinner uk-icon-spin") : false;

        return <ul class="uk-list uk-list-space">{waiting || items}</ul>
    }
}

var ItemView = {
    controller: function(c) {
        var api = {
            onEdit: function(mode_name) {
                return function(e) {
                    e.preventDefault();

                    m.route("/buckets/edit/" + args.data.Id)

                    return false;
                }
            }
        }

        api = _.extend({}, c, api);

        return api;
    },
    view: function(c, args) {
        var urlFiles = "/buckets/" + args.data.Id + "/files/search";
        var urlEdit = c.resource_name == ItemActions.ITEMS?
            "/buckets/" + m.route.param("bucket_id") + "/files/edit/"+args.data.Id:
            "/buckets/edit/"+args.data.Id;
        
        var isShowSubPositions = c.resource_name == ItemActions.CLASSIFERS;
        // var buttonPositions = <a class="uk-button" title="Файлы" href={url} config={m.route}><i class="uk-icon uk-icon-files-o"></i></a>

        var listRefWidgets = [];

        if (args.data.Props._BuildTraceWidgets) {
            _.each(args.data.Props._BuildTraceWidgets, function(trace){
                var url = "/buckets/"+trace.Bucket+"/files/edit/"+trace.Id;
                var label = trace.Description + ' "'+trace.ExtId+'" ('+trace.Collections.join(", ")+')';
                
                listRefWidgets.push(m("li", m("a", {href: url, config: m.route}, label)))
            })
        }
        
        /*
        <div class="uk-width-1-10">
                            {isShowSubPositions ? buttonPositions : ""}
                        </div>
        */
        
        return <div class="uk-panel uk-panel-box uk-panel-hover" id={args.data.ExtId + ":" + args.data.Id} >
                    <div class="uk-panel-badge uk-text-small uk-text-muted">{args.data.ExtId}</div>
                    <p class="uk-panel-title">{args.data.Description}</p>
                    <p class="uk-article-meta"><i class="uk-icon uk-icon-database"></i> {args.data.Collections} <i class="uk-icon uk-icon-tags"></i> {args.data.Tags}</p>
                    <p>related widgets:</p>
                    <ol class="">{listRefWidgets}</ol>
                    <div>
                        <a class="uk-button" href={urlEdit} config={m.route}><i class="uk-icon uk-icon-edit"></i> Edit</a>
                        {isShowSubPositions? <a href={urlFiles} config={m.route} class="uk-button uk-button-primary"><i class="uk-icon uk-icon-folder-open"></i> Files</a>: ""}
                    </div>
                </div>
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
        
        // TODO: визуализация временных записей
        
        // if (args.data._TempId) {
        //     c.mode("edit");
        // }
        
        // var special_id = m.route.param("special_id");
        
        // if (special_id && (special_id == args.data.Id || special_id == args.data.ExtId)) {
        //     m.route("/classifers/"+args.bucket_id+"/items")
        //     c.mode("edit");
        // }
        
        return m.component(ItemView, _.extend({}, args, { mode: c.mode, bucket_id: args.bucket_id, resource_name: c.resource_name }));

        // switch (c.mode()) {
        //     case "edit":
        //         return m.component(mapEditor[c.resource_name], _.extend({}, args, { mode: c.mode, bucket_id: args.bucket_id }));
        //     case "view":
        //         return m.component(ItemView, _.extend({}, args, { mode: c.mode, bucket_id: args.bucket_id }));
        //     default:
        //         return m("li", args.data.Id)
        // }
    }
}



var Manager = {
    controller: function(c) {
        var bucket_id = m.route.param("bucket_id");

        var api = {
            bucket_id: bucket_id,
        }

        return api;
    },
    view: function(c, args) {
        var filter = m.component(SearchFilter, { resource_name: args.resource_name, bucket_id: c.bucket_id})
        // var result = m.component(SearchResult, { resource_name: args.resource_name, bucket_id: c.bucket_id, onSearch: c.onSearch })

        return m("div.uk-margin-top", [filter]);
    }
}

var NewFile = {
    controller: function(c) {
        var resource_id = m.route.param("resource_id");
        var resource_name = c.resource_name;
        var bucket_id = m.route.param("bucket_id");
        
        if (!c.data) {
            c.data = {}
        }
        
        var isInit = m.prop(false);
        
        var api = {
            resource_name: resource_name,
            bucket_id: bucket_id,
            config: function(element, initialized, context) {
                isInit(initialized);  
            },
            ExtId: m.prop(""),
            Id: m.prop(0),
            Articul: m.prop(""),
            
            Description: m.prop(""),
            Tags: m.prop([]),
            Props: {
                Config: m.prop(""),
                Content: m.prop(""),
            },
            CustomProps: [], // [{key, value}, ...]
            onunload: function() {
                if (typeof this.unlistener == "function" && isInit()) {
                    this.unlistener();
                }                
            },
            unlistener: null,
            listener: function() {
                this.unlistener = store.subscribe(function() {
                    switch (store.getState().lastAction.type) {
                        case ActionTypes.UPDATE_ITEM_SUCCESS:
                            alert("updated ok")
                            break;
                        case ActionTypes.CREATE_ITEM_SUCCESS:
                            alert("created ok")
                            if (store.getState().lastAction.item.Code == "success") {
                                m.route("/buckets/edit/"+store.getState().lastAction.item.Data.Id)    
                            } else {
                                m.route("/buckets/search")
                            }
                                
                            break;
                        case ActionTypes.GET_ITEM_SUCCESS:
                            var _propSeq = 1;
                            if (store.getState().lastAction.response.Code == "success") {
                                var file = store.getState().lastAction.response.Data; 
                                
                                this.Id(file.Id)
                                this.ExtId(file.ExtId)
                                this.Articul(file.Articul)
                                this.Description(file.Description)
                                this.Tags(file.Tags)
                                this.Props.Config(file.Props.Config)
                                this.Props.Content(file.Props.Content)
                                _.each(file.Props, function(value, key){
                                    if (key == "Config") {
                                        this.Props.Config(value)
                                        return
                                    }
                                    
                                    if (key == "Content") {
                                        this.Props.Content(value)
                                        return
                                    }
                                    
                                    this.CustomProps.push({key: m.prop(key), value: m.prop(value), id: _propSeq++});
                                }.bind(this))
                            } else {
                                m.route("/buckets/new")    
                            }
                            break;
                        case ActionTypes.GET_ITEM_FAIL:
                            m.route("/buckets/new")
                            break;
                    }
                    
                    // m.endComputation();
                    
                }.bind(this));
            },
            onCreate: function() {
                return function(e) {
                    e.preventDefault();
                    
                    var data = {
                        Id: this.Id,
                        ExtId: this.ExtId(),
                        Articul: this.Articul(),
                        
                        Description: this.Description(),
                        Tags: this.Tags(),
                        Props: {},
                    }
                    
                    _.each(this.CustomProps, function(item){
                        data.Props[item.key()] = item.value(); 
                    })
                    
                    data.Props["Config"] = this.Props.Config();
                    data.Props["Content"] = this.Props.Content();
                    
                    if (!this.Id()) {
                        ItemActions.Create(store.dispatch, { data: data, resource_name: this.resource_name, bucket_id: this.bucket_id });
                    } else {
                        ItemActions.Update(store.dispatch, { id: this.Id(), data: data, resource_name: this.resource_name, bucket_id: this.bucket_id });
                    }
                    
                    return false;
                }.bind(this);
            },
            onDeleteFile: function() {
                return function(e) {
                    e.preventDefault();
                    
                    if (!confirm("Вы уверены в удалении?")) {
                        return false;
                    }
                    
                    if (this.Id) {
                        ItemActions.Delete(store.dispatch, { id: this.Id(), resource_name: this.resource_name, bucket_id: this.bucket_id});   
                    }
                    
                    return false;
                }.bind(this);
            },
            onAddProperty: function() {
                return function(e) {
                    e.preventDefault();
                    
                    this.CustomProps.push({id: _.now(), key: m.prop(""), value: m.prop("")});
                    
                    return false;
                }.bind(this);
            },
            onRemoveProperty: function(id) {
                 return function(e) {
                    e.preventDefault();
                    
                    this.CustomProps = _.filter(this.CustomProps, function(item){ return item.id != id;});
                    
                    return false;
                 }.bind(this);
            },
            initEditor: function(config) {
                    return function(e) {
                        var editor = ace.edit(config.id);
                        editor.getSession().setMode(config.mode); // 'ace/mode/html'
                        // editor.setTheme('ace/theme/monokai');
                        editor.setValue(config.value());
                        editor.getSession().on('change', function() {
                            config.value(editor.getSession().getValue());
                        });
                        editor.$blockScrolling = Infinity;
                    }
                    
                },
        };
        
        api = _.merge({}, api, c);
        api.listener();
        
        if (resource_id) {
            ItemActions.Get(store.dispatch, {id: resource_id, bucket_id: bucket_id, resource_name: resource_name})
        }
        
        return api;
    },
    view: function(c, args) {
        var configId = [c.data.Id, c.data._TempId, "config"].join(":");
        var contentId = [c.data.Id, c.data._TempId, "content"].join(":");
        
        var tabSequence = 1;
        var customProps = _.map(c.CustomProps, function(item){
            
            return <div class="uk-form-controls">
                        <a href="" tabindex="-1" class="uk-button uk-width-1-10" onclick={c.onRemoveProperty(item.id)}><i class="uk-icon uk-icon-remove"></i></a>
                        <input tabindex={tabSequence++} type="text" placeholder="key" class="uk-width-3-10" onchange={m.withAttr("value", item.key) } value={item.key() }/>
                        <input tabindex={tabSequence++} type="text" placeholder="value" class="uk-width-6-10" onchange={m.withAttr("value", item.value) } value={item.value() }/>
                    </div>
        });
        
        return <div class="uk-panel uk-panel-box" id={c.ExtId()+":"+c.Id()} config={c.config}>
                <form class="uk-form uk-form-horizontal" onsubmit={c.onCreate() }>
                    
                    <div class="uk-form-row">
                        <label class="uk-form-label">ExtId</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" onchange={m.withAttr("value", c.ExtId) } value={c.ExtId() }/>
                        </div>
                    </div>
                    
                    <div class="uk-form-row">
                        <label class="uk-form-label">Articul</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" onchange={m.withAttr("value", c.Articul) } value={c.Articul() }/>
                        </div>
                    </div>
                    
                    <div class="uk-form-row">
                        <label class="uk-form-label">Description</label>
                        <div class="uk-form-controls">
                            <input type="text" placeholder="" class="uk-width-1-1" onchange={m.withAttr("value", c.Description) } value={c.Description() }/>
                        </div>
                    </div>
                    
                    <div class="uk-form-row">
                        <label class="uk-form-label">Config</label>

                        <div class="uk-form-controls">
                            <div id={configId} config={c.initEditor({value: c.Props.Config, id: configId, mode: "ace/mode/toml"}) }  style="height:200px;"></div>
                        </div>
                    </div>

                    <div class="uk-form-row">
                        <label class="uk-form-label">Content</label>

                        <div class="uk-form-controls">
                            <div id={contentId} config={c.initEditor({value: c.Props.Content, id: contentId, mode: "ace/mode/html"}) }  style="height:200px;"></div>
                        </div>
                    </div>
                    
                    <div class="uk-form-row">
                        <label class="uk-form-label">Properties</label>
                        {customProps}
                        <div class="uk-form-controls">
                            <a href="#" onclick={c.onAddProperty()} class="uk-button uk-button-small uk-width-1-1"><i class="uk-icon uk-icon-plus"></i></a>
                        </div>
                    </div>
                    
                    <div class="uk-form-row">
                        <a href="" class="uk-button uk-button-danger uk-float-right" onclick={c.onDeleteFile()}>Delete</a>
                        
                        <button class="uk-button uk-button-primary uk-width-1-4">Save</button>
                    </div>
                </form>
            </div>
    }
}


module.exports = {
    Manager: Manager,
    NewItem: NewFile
}