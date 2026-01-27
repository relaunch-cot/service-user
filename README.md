# User Service - Relaunch

Microserviço responsável pela gestão completa de usuários na plataforma Relaunch, implementando funcionalidades essenciais de autenticação, cadastro e gerenciamento de perfis.

## 🎯 Sobre o Projeto

Este serviço faz parte da arquitetura de microserviços do Relaunch, concentrando-se exclusivamente nas operações relacionadas aos usuários da plataforma.

## ⚙️ Funcionalidades

- **Cadastro de Usuários**: Criação de novas contas com validação de dados e criptografia de senhas usando bcrypt
- **Autenticação**: Sistema de login seguro com validação de credenciais
- **Gerenciamento de Perfil**: Atualização de informações pessoais, configurações e foto de perfil
- **Recuperação de Senha**: Fluxo completo para redefinição de senha
- **Busca de Usuários**: Pesquisa por nome com suporte a filtros
- **Exclusão de Conta**: Remoção segura de dados do usuário

## 🛠️ Tecnologias

- **Linguagem**: Go (Golang)
- **Banco de Dados**: MySQL
- **Comunicação**: gRPC com Protocol Buffers
- **Segurança**: bcrypt para hash de senhas
- **Gerenciamento de Contexto**: Context para controle de requisições

## 📦 Estrutura

O serviço implementa uma camada de repositório (`repositories/mysql`) que abstrai toda a lógica de persistência de dados, seguindo princípios de Clean Architecture e facilitando manutenção e testes.

## 🔐 Segurança

- Senhas armazenadas com hash bcrypt (cost factor 14)
- Validação de duplicidade de emails
- Autenticação segura antes de operações sensíveis
- Gestão adequada de erros com códigos gRPC status

## 🚀 Diferenciais

- Arquitetura baseada em microserviços
- Comunicação eficiente via gRPC
- Código modular e testável
- Suporte a configurações personalizadas por usuário (skills, biografia, etc.)

---

*Desenvolvido como parte do ecossistema Relaunch*
