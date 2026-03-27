# actajus

## ActaJus é um projeto de desenvolvimento de software que visa simplificar e facilitar a vida das pessoas ao utilizar uma plataforma de gestão de processos de negócio.

## Pré-requisitos

#

# - Docker

# - Docker Compose

# - Git

## Instalação

#

# 1. Clone o repositório (se necessário):

```bash
git clone https://github.com/paladignus/actajus.git
cd actajus
```

2. Instale as dependências:

```bash
docker compose build
docker compose up -d
```

3. Acesse a plataforma em http://localhost:8080

4. Utilize o usuario padrao: admin e senha: admin123

## Documentação

#

# A documentação completa pode ser encontrada em [https://github.com/paladignus/actajus](https://github.com/paladignus/actajus).

- Regras arquiteturais atuais: [ARCHITECTURE.md](/home/marcelo/Workspace/go/projects/actajus/ARCHITECTURE.md)

## Contribua para o projeto

Se vocês gostaria de contribuir para o projeto, por favor, visite [https://github.com/paladignus/actajus](https://github.com/paladignus/actajus) para mais informaçoes.

## Contato

#

# Para mais informaçoes, visite [https://github.com/paladignus/actajus](https://github.com/paladignus/actajus).

## Novos arquivos para publish and dispatch

```internal/
  module/
    identity/
      application/
        dto/
          request_password_reset_command.go
          request_password_reset_read_model.go
        event/
          subjects.go
          email_send_requested.go
        mapper/
          auth_mapper.go
        repository/
          outbox_repository.go
          password_reset_repository.go
          user_repository.go
        service/
          clock.go
          refresh_token_service.go
        usecase/
          request_password_reset.go
      domain/
        model/
          password_reset_token_create.go

    notification/
      application/
        dto/
          email_send_requested.go
        repository/
          sent_email_repository.go
        service/
          email_sender.go
        usecase/
          process_email_send_requested.go
      infrastructure/
        consumer/
          email_send_requested_consumer.go

  shared/
    application/
      messaging/
        outbox_message.go
        publisher.go
      service/
        id_generator.go
        message_serializer.go
      uow/
        unit_of_work.go
    infrastructure/
      messaging/
        nats/
          jetstream_publisher.go
          bootstrap.go
        outbox/
          dispatcher.go
      persistence/
        postgres/
          unit_of_work.go
          outbox_repository.go
      serialization/
        json_serializer.go
      service/
        uuid_generator.go

internal/shared/application/uow/unit_of_work.go
internal/shared/application/service/message_serializer.go
internal/shared/application/service/id_generator.go
internal/shared/application/messaging/outbox_message.go
internal/shared/application/messaging/publisher.go
internal/shared/infrastructure/serialization/json_serializer.go
internal/shared/infrastructure/service/uuid_generator.go
internal/shared/infrastructure/persistence/postgres/unit_of_work.go
internal/shared/infrastructure/persistence/postgres/outbox_repository.go
internal/shared/infrastructure/messaging/nats/jetstream_publisher.go
internal/shared/infrastructure/messaging/outbox/dispatcher.go
internal/shared/infrastructure/messaging/nats/bootstrap.go

internal/module/identity/application/repository/user_repository.go, já existe e mais completo
internal/module/identity/application/repository/password_reset_repository.go já existe e mais completo
internal/module/identity/application/repository/outbox_repository.go
internal/module/identity/application/service/clock.go, já existe
internal/module/identity/application/service/refresh_token_service.go, já existe e mais completo
internal/module/identity/application/event/subjects.go
internal/module/identity/application/event/email_send_requested.go
internal/module/identity/domain/model/password_reset_token_create.go, já existe
internal/module/identity/application/dto/request_password_reset_command.go, já existe
internal/module/identity/application/dto/request_password_reset_read_model.go, já existe e MODIFICADO
internal/module/identity/application/usecase/request_password_reset.go, já existe e MODIFICADO

internal/module/notification/application/dto/email_send_requested.go
internal/module/notification/application/service/email_sender.go
internal/module/notification/application/repository/sent_email_repository.go
internal/module/notification/application/usecase/process_email_send_requested.go
internal/module/notification/infrastructure/consumer/email_send_requested_consumer.go
internal/module/notification/infrastructure/persistence/postgres/sent_email_repository.go
```
