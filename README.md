# 2024-spring-AB-Go-HW-1-template

Необходимо сделать общий чат.
Нужно реализовать и клиента и сервер.

Клиент - консольное приложение, которое подключается по websocket к серверу. При старте клиента, в консоле запрашивается nickname.
При подключении к серверу, клиент получает последние 10 сообщений чата и все будущие сообщения. Сообщения на клиенте должны отображаться в консоле в формате: "<имя пользователя>: <сообщение>".

Сервер, при получении сообщения, должен выводить их в консоль, сохранять в БД и рассылать их по клиентам. В качестве базы данных необходимо использовать postgres.

На клиенте и сервере необходимо реализовать Graceful Shutdown.
