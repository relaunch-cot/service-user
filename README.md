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
- Rodar 'go build main.go' no terminal
- Rodar 'go run main.go' no terminal

## Funcionalidade de Relatórios em PDF via gRPC

Este microserviço expõe métodos gRPC para geração de relatórios em PDF que devem ser chamados pelo BFF.

### Métodos gRPC Disponíveis

#### 1. GenerateReportFromJSON  
- **Input**: `GenerateReportRequest` (JSON string)
- **Output**: `GenerateReportResponse` (bytes do PDF)
- **Descrição**: Gera PDF personalizado baseado em dados JSON

### Estruturas para o arquivo .proto

```proto
message GenerateReportRequest {
  string json_data = 1;
}

message GenerateReportResponse {
  bytes pdf_data = 1;
}

service UserService {
  // Métodos existentes...
  rpc ExportUserReportPDF(google.protobuf.Empty) returns (GenerateReportResponse);
  rpc GenerateReportFromJSON(GenerateReportRequest) returns (GenerateReportResponse);
}
```

### Formato JSON para Relatórios Personalizados

```json
{
  "title": "Título do Relatório",
  "subtitle": "Subtítulo (opcional)",
  "headers": ["Coluna 1", "Coluna 2", "Coluna 3"],
  "rows": [
    ["Dados 1", "Dados 2", "Dados 3"],
    ["Dados 1", "Dados 2", "Dados 3"]
  ],
  "footer": "Rodapé (opcional)"
}
```

### Implementação no BFF

Veja o arquivo `docs/grpc_pdf_integration.md` para exemplos completos de como implementar no BFF.

## Funcionalidades
### Backend/Frontend
- [x]  Permitir login do usuário
- [x]  Permitir cadastro do usuário
- [x]  Usuário redefinir  a senha
- [x]  Permitir deletar usuário
- [x]  Permitir atualizar email e nome
- [ ]  O usuário deve poder personalizar as configurações do perfil
- [x]  Deve ser possível exportar relatórios em PDF.
- [ ]  O freelancer define o tempo para o desenvolvimento da aplicação

### Frontend
- [x]  Usuário fazer logout da plataforma
- [ ]  O usuário deve conseguir selecionar o tema da plataforma
