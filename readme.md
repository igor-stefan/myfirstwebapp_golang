### Projeto de aplicação web utilizando a linguagem GO
Repositótio utilizado para armazenar/controlar os arquivos. Vou escrevendo os updates conforme aprendo.

* Inicio em 10/05/2022
* Até o momento possui:
    - sessions; 
    - middlewares; 
    - proteção csrf
    - templates html
    - handlers

# cmd/web
Constitui do comando da aplicação. As ferramentas essenciais para o funcionamento, a camada mais superficial.

## main.go
Arquivo possui duas funções, a main() e a run().
A função run() retorna uma conexão com o database da aplicação através de um objeto to tipo 'myDriver', que utiliza o driver *pgx* para estabelecer a conexão com o database SQL. Nesta função também são estabelecidos as configurações iniciais da *session*, o *channel* para o envio de emails e as configurações gerais da aplicação que são utilizadas por outros packages (render, handlers, helpers, etc.). Além disso são carregados os templates *html* da aplicação.
A função main() executa a função run(), estabelece *defers* para fechamento do *channel* e a conexão com o db. Por último, inicializa o server e o coloca para escutar as solicitações.

## middleware.go

## routes.go

## send-mail.go
