# Nome do Projeto: ReLaunch

## Integrantes:
- Matheus Oliveira Mangualde - 22301194
- Henrique de Freitas Issa - 22300732
- João Pedro Bastos Neves - 22301330
- Eduardo Mapa Avelar Damasceno - 22301674
- Eike Levy Albano Neves - 22402772
- Vinícius Theodoro Giovani - 22300821

**Turma 3B2**

# Como rodar
- Baixar o golang
- Setar no terminal 'go env -w GOPRIVATE=*' para conseguir acessar os repositorios privados do github
- Rodar 'go mod tidy' no terminal para instalar as dependencias
- Setar as variaveis de ambiente:
  - PORT: (porta em que o microserviço vai rodar)
  - ## Variáveis de ambiente referente ao banco de dados:
    **Lembrar de rodar o MySql localmente com uma instancia para o banco de dados que contenha uma tabela 'users' para alimentar as requisições**
    - MYSQL_HOST: (host do banco de dados MySql)
    - MYSQL_PORT: (porta em que o banco de daods MySql está rodando)
    - MYSQL_USER: (usuario do MySql)
    - MYSQL_PASS: (senha do MySql)
    - MYSQL_DBNAME: (nome do banco de dados que receberá as requisições)
    - SENDGRID_API_KEY (chave de api do sendgrid)
    - EMAIL (email do usuário que receberá a mensagem para recuperação de senha)
    - NAME (nome do usuário que aparecerá no email de recuperação de senha)
    - JWT_SECRET (secret key utilizada na criação do token jwt do usuário ao fazer login)
- Rodar 'go build main.go' no terminal
- Rodar 'go run main.go' no terminal

## Funcionalidades implementadas
- [x]  Permitir login do usuário
- [x]  Permitir cadastro do usuário
- [x]  Usuário redefinir  a senha
- [x]  Permitir deletar usuário
- [x]  O usuário deve poder personalizar as configurações do perfil
- [x]  Buscar informações de perfil do usuario
- [x]  Deve ser possível exportar relatórios em PDF.
- [x]  Enviar Email de recuperação de senha
- [x]  Usuário deletar sua conta
- [x]  Usuário fazer logout da plataforma
- [x]  Criar um novo chat entre usuarios
- [x]  Enviar mensagens no chat entre usuários
- [x]  Buscar todas as mensagens de um chat específico
- [x]  Buscar todos os chats de um usuário
- [x]  Criar um novo projeto (usuários que sejam clientes)
- [x]  Buscar um projeto específico
- [x]  Buscar todos os projetos de um usuário
- [x]  Adicionar freelancer a um projeto
- [x]  Remover freelancer de um projeto
- [x]  Listar todos os projetos que estejam sem um freelancer desenvolvendo o mesmo, ou seja, disponíveis para desenvolvimento
- [x]  Enviar norificações para o usuário (seja de uma mensagem nova, seja de solicitação para participar de um projeto...)
- [x]  Buscar informações de uma notificação específica
- [x]  Buscar todas as notificações de um usuário

## Padrões requisitados
- padrão singleton aplicado
### Padrões GoF aplicados além do singleton:
- Adapter
- Facade
- Strategy
- Factory
- Iterator
### além disso o projeto também aplica padrões de arquitetura (Repository, Dependency Injection) que não são parte dos GoF clássicos, mas complementam a estrutura.
