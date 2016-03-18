# GONG

simple constructor web site via widgets

Конструктор веб-сайта с синтаксисом markdown и структурой схожей с wiki.
Особенностью является возможность апперировать блоками с динамическим контентом - виджеты.

Вместа конструкции ```{{widget . "foobar"}}``` отображается значение виджета с ключем ```foobar```.
Отредактировать виджет можно по ссылке ```/@widgets:edit/foobar``` где _foobar_ название виджета. 

Страница это виджет. 
Виджет содержит множество виджетов. 
Виджет может быть обернут в layout. 
Виджет будет представлен в layout-виджете в параметре ```{{ V.Content }}```

Каждая страница дает возможность задать параметры (в [формате toml](https://github.com/toml-lang/toml#user-content-example)).

## Редактирование виджета

![Страница редактирования виджета](https://s3.amazonaws.com/idheap/ss/127.0.0.18080widgetseditdefault_2016-03-18_09-37-17.png)


## Настройки виджета

``` toml
title = "page title" # if used layout

[self]
render = "" # or markdown
# if the render=markdown in the widget data should not have dynamic parameters

# layout = "" # current widget will be in the .V.Content variable
link_edit = "/@widgets:edit/{{.V.Name}}" # example dynamic parameter
link_title = "edit"
```

## RUN

``` shell
go run main.go -db=gong.db -bind=:8080 -stderrthreshold=INFO
```

## Примеры

### Cоздание страницы
 
![Процесс редактирования](https://s3.amazonaws.com/idheap/ss/screencast_2016-03-18_09-26-21.gif)

* Переход по ссылке /
* Включив параметр editable=1 отображаем вспомогательные ссылки (для быстрого перехода на страницу редактирования виджета)
* Задаем значение виджета страницы /
* Определяем тип содрежимого виджета как markdown (prerender=md)
* Задаем layout в настройках страницы (layout=default)
* Указанный layout отображается в связанных виджетах. Переходим к редактированию виджета _default_
* Указали layout. Просматриваем полученную страницу
* Задаем значение виджета _brand_

### layout 

``` html
<!DOCTYPE html>
<html>
	<head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/uikit/2.25.0/css/uikit.almost-flat.min.css">
	</head>
    <body class="uk-position-relative">
        {{ widget . "top" }}
        <nav class="uk-navbar uk-navbar">
            <a href="/" class="uk-navbar-brand">{{widget . "brand" }}</a>
            <ul class="uk-navbar-nav">
                {{ widget . "navbar" }}
            </ul>
            <div class="uk-navbar-flip">
                <ul class="uk-navbar-nav">
                </ul>
            </div>
        </nav>
	    {{ .V.Content }}
        <footer>
        {{ widget . "footer" }}
        </footer>
	</body>
</html>
```
