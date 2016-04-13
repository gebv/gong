# GONG


[![Build Status](https://travis-ci.org/gebv/gong.svg?branch=master)](https://travis-ci.org/gebv/gong)

!Пример ниже не соответствует обновленному интерфейсу взаимодействия с приложением

constructor web site via widgets

Конструктор веб-сайта.
Особенностью является возможность апперировать блоками с динамическим контентом - виджеты.

Вместа конструкции ```{{widget . "widgets" "foobar"}}``` отображается значение виджета с ключем ```foobar```.
Отредактировать виджет можно по ссылке ```/@settings/classifers/widgets/items/edit?special_id=foobar``` где _foobar_ название виджета в подмножестве _widgets_. 

Страница это виджет. 
Виджет содержит множество виджетов. 
Виджет может быть обернут в layout. 
Виджет будет представлен в layout-виджете в параметре ```{{ V.Content }}```

Каждая страница дает возможность задать параметры (в [формате toml](https://github.com/toml-lang/toml#user-content-example)).

## Редактирование виджета

![Страница редактирования коллекции виджетов](https://s3.amazonaws.com/idheap/ss/192.168.1.368081settingsclassifers_2016-04-02_22-00-22.png)

![Страница радктирования виджета](https://s3.amazonaws.com/idheap/ss/192.168.1.368081settingsclassifersbc3f1a64-f8c1-11e5-9cf5-10ddb19b9d24itemseditspecial_id2Fwiki_2016-04-02_22-02-25.png)


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

### Применение

* WIKI
* Протипирование сайтов. Альтернатива axure, pencel, .... Cоздать прототип максимально приближенный к реальности
* Ведение блога (после реализации коллекций и системы управления с авторизацией для модератора)

# TODO

* закрыть от всех настройки (например http-авторизация)
* изменить роутинг в настроке виджетов
* возможность экспортировать\импортировать настройки
* реализовать интерфейс для написания доп.компонентов\модулей
* для модели реализовать функции для быстрой валидации данных [(Required, Length, Equals, Greater Than, Greater Than or Equal, Less Than, Less Than or Equal, Contains, etc)](https://godoc.org/gopkg.in/bluesuncorp/validator.v8)

# Идеи

* реализовать коллекции, интерфейс работы с ними
* предоставить возможность создавать контроллеры через виджеты, например для обработки формы (csrf)
