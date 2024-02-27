# Rodando o projeto

Esse projeto utiliza o RabbitMQ como Message Broker e Postgres como banco de dados relacional.
Para facilitar a configuração, basta iniciá-los utilizando o docker compose da seguinte forma:

```bash
docker compose up db rabbitmq
```

Existem dois entrypoints na aplicação, um deles sendo uma API e o outro sendo um worker.
Para automatizar a parte de configuração das variáveis de ambiente, compilação e execução, estou utilizando um Makefile.
Para executá-los basta usar os seguintes commandos na raíz do projeto:

```bash
make api
make worker
```

# Rotas da API

## Criação de cliente
URL base local: localhost:8765

POST {base_url}/clients

Body:
```json
{
    "name": "breno",
    "surname": "almeida",
    "contact": "99999999999",
    "address": "rua abc",
    "birthdate": "09-21-1993",
    "cpf": "99999999999"
}
```

Essa rota é responsável por criar um evento de criação de usuário e enviá-lo para uma fila do Rabbitmq para ser processado posteriormente pelo worker.
O evento criado possui todas essas informações do cliente, mas antes de ser enviado para a fila do Rabbitmq ele é criptografado utilizado criptografia AES para garantir a criptografia dos dados em repouso. Essa é uma abordagem padrão para lidar com dados sensíveis e mais indicada do que apenas criptografar os dados durante o tráfego.

O evento possui o seguinte formato:

```go
type CreateClientEvent struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Contact   string `json:"contact"`
	Address   string `json:"address"`
	Birthdate string `json:"birthdate"`
	Cpf       string `json:"cpf"`
}
```

## Atualização de um cliente

PATCH {base_url}/clients/:id

Body:
```json
{
    "name": "nome atualizado",
    "surname": "sobrenome atualizado",
    "contact": "99999999999",
}
```

Essa rota atualiza os dados do cliente no banco. O que for enviado para atualização será atualizado e o que for omitido do payload não será atualiazado, mantendo o seu valor atual do banco.
Caso não exista usuário com o ID fornecido, é retornada uma resposta de registro não encontrado.

## Deleção de um cliente

DELETE {base_url}/clients/:id

Essa rota deleta um cliente do banco que possui um determinado id. Caso o usuário não exista, é retornada uma resposta de registro não encontrado.

## Listagem de Clientes

GET {base_url}/clients?page=1&size=10

Essa rota retorna os clientes de forma paginada. Os parâmetros page e size são utilizados para indicar a página e o tamanho da página respectivamente.

A resposta possui o seguinte formato:

```json
{
    "data": [
        {
            "ID": 1,
            "Name": "updated",
            "Surname": "updated2",
            "Contact": "99999999999",
            "Address": "rua abc",
            "Birthdate": "1993-09-21T00:00:00Z",
            "Cpf": "99999999999",
            "CreatedAt": "2024-02-27T14:31:00.193216Z",
            "UpdatedAt": "2024-02-27T14:31:00.193216Z"
        }
    ],
    "total": 1,
    "page": 1,
    "page_size": 5
}
```