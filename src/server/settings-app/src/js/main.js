var m = require("mithril");
var store = require('./store');

var Dashboard = require("./components/Dashboard");
var ItemsManager = require("./components/Items");

// TODO: Вынести в настройки
var Items = require("./actions/Items");

m.route.mode = "hash";

// TODO: выводить информацию о том в какой категории

var Header = {
    view: function(c, args) {
        var resource_name = m.route.param("resource_name");
        
        return <nav class="uk-navbar uk-navbar">
                    <a href="/" config={m.route} class="uk-navbar-brand">@</a>
                    <ul class="uk-navbar-nav uk-hidden-small">
                        <li class={resource_name=="buckets"?"uk-active":""}><a  href="/buckets/search" config={m.route}>Buckets</a></li>
                    </ul>

                    <div class="uk-navbar-flip">                        
                        <ul class="uk-navbar-nav uk-hidden-small">
                            <li><a href="/@settings/login" config={m.route}>Login</a></li>
                        </ul>
                    </div>
                </nav>
    }
}

var Layout = {
    controller: function(c) {
        var api = {
            content: m("div"),
            onunload: function() {
                console.log("unload", "Layout")
            }
        };
        
        api.content = m.component(c.content, c.args);
        
        return api;
    },
    view: function(c, args) {
        var header = m.component(Header)
        
        var contentWrap = m("div.uk-grid", m("div.uk-width-small-1-1 uk-width-medium-9-10 uk-container-center", c.content));
        
        return m("div", [header, contentWrap])
    }
}

m.route(document.body, "/", {
    "/": m.component(Layout, {content: Dashboard, args: {}}),
    "/buckets/search": m.component(Layout, {content: ItemsManager.Manager, args: {resource_name: Items.CLASSIFERS}}),
    "/buckets/new": m.component(Layout, {content: ItemsManager.NewItem, args: {resource_name: Items.CLASSIFERS}}),
    "/buckets/edit/:resource_id": m.component(Layout, {content: ItemsManager.NewItem, args: {resource_name: Items.CLASSIFERS}}),
    
    "/buckets/:bucket_id/files/search": m.component(Layout, {content: ItemsManager.Manager, args: {resource_name: Items.ITEMS}}),
    
    "/buckets/:bucket_id/files/edit/:resource_id": m.component(Layout, {content: ItemsManager.NewItem, args: {resource_name: Items.ITEMS}}),
    "/buckets/:bucket_id/files/new": m.component(Layout, {content: ItemsManager.NewItem, args: {resource_name: Items.ITEMS}}),
})