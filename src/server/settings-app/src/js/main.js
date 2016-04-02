var m = require("mithril");
var store = require('./store');

var Dashboard = require("./components/Dashboard");
var ItemsManager = require("./components/ItemsManager");

// TODO: Вынести в настройки
var Items = require("./actions/Items");

m.route.mode = "pathname";

// TODO: выводить информацию о том в какой категории

var Header = {
    view: function(c, args) {
        var resource_name = m.route.param("resource_name");
        
        return <nav class="uk-navbar uk-navbar">
                    <a href="/" config={m.route} class="uk-navbar-brand">@</a>
                    <ul class="uk-navbar-nav uk-hidden-small">
                        <li class={resource_name=="classifers"?"uk-active":""}><a  href="/@settings/classifers" config={m.route}>Classifers</a></li>
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
    view: function(c, args) {
        var header = m.component(Header)
        var content = m.component(args.content, args.config);
        
        var contentWrap = m("div.uk-grid", m("div.uk-width-small-1-1 uk-width-medium-9-10 uk-container-center", content));
        
        return m("div", [header, contentWrap])
    }
}

m.route(document.body, "/@settings/", {
    "/@settings/": m.component(Layout, {content: Dashboard, config: {}}),
    "/@settings/:resource_name": m.component(Layout, {content: ItemsManager.Manager, config: {resource_name: Items.CLASSIFERS}}),
    "/@settings/classifers/:classifer_id/:resource_name": m.component(Layout, {content: ItemsManager.Manager, config: {resource_name: Items.ITEMS}}),
})