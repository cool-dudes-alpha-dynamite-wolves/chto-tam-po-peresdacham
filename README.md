# [WIP] chto-tam-po-peresdacham

### Функциональные требования

- Excel таблицы будут скачиваться вручную, напрямую с сайта университета. Они будут храниться статически, напрямую на сервере 
(т.к. на первом этапе изменение данных возможно раз в семестр, когда и публикуют список пересдач).
- *Должна существовать реализация рестарта сервиса, на случай, если упадет. Самый простой способ реализовать это - написать Bash скрипт, проверяющий процесс. (пока можно забить)
- Процесс парсинга проходит не на "голых" файлах с сайта университета, а предварительно обработанных вручную (в будущем планируется упразднить данный шаг, если на сайте университета будут выкладывать одинаковые таблицы с пересдачами).

### Не функциональные требования

- Парсим только файлы в формате .XLSX
- На первом этапе позволительно парсить не все виды таблиц, т.к. многие из них отличаются названием полей/столбцов, наполнением и т.д.
Достаточно поддержать парсинг наиболее распространенного шаблона.
- Никаких оптимизаций по взаимодействию с telegram API и с парсингом .xsls файлов на данном этапе не требуется.