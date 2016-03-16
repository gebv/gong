# GONG

simple constructor web site via widgets

Конструктор веб-сайта с синтаксисом markdown и структурой схожей с wiki.
Особенностью является возможность оперировать блоками с динамическим контентом - виджеты.

Страница имеет уникальный URL. На странице размещены виджеты
Вместа конструкции ```{{widget . "/foobar"}}``` отображается значение виджета с ключем "/foobar"

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
        {{ widget . "/top" }}
        <nav class="uk-navbar uk-navbar">
            <a href="/" class="uk-navbar-brand">{{widget . "/brand" }}</a>
            <ul class="uk-navbar-nav">
                {{ widget . "/navbar" }}
            </ul>
            <div class="uk-navbar-flip">
                <ul class="uk-navbar-nav">
                </ul>
            </div>
        </nav>
	    {{ .V.Content }}
        <footer>
        {{ widget . "/footer" }}
        </footer>
	</body>
</html>
```

![Страница и log запрашиваемых виджетов](https://s3.amazonaws.com/idheap/ss/localhost8080page_2016-03-15_16-55-59.png)

Ниже пример редактирования виджетов 
![Процесс редактирования](full_example.gif)

[Оригинальная gif-демо](https://s3.amazonaws.com/idheap/ss/screencast_2016-03-15_16-59-01.gif)

* задается контент страницы /page (спец URL /@pages:edit/page
* задается значение виджета /brand (спец URL /@widgets:edit/brand)
* редактируется значение виджета /navbar (спец URL /@widgets:edit/navbar) - добавляется кнопка редактирования страницы
