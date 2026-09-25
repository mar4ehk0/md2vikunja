# md2vikunja

Читает задачи из markdown-файла и создаёт их в Vikunja. Ссылка на созданную задачу
записывается обратно в файл, поэтому повторный запуск дубли не создаёт.

## Конфиг

`config.yaml` (в `.gitignore`):

```yaml
url: https://vikunja.example.com
token: tk_xxx
project_id: 1
```

`token` - API-токен из настроек Vikunja. `project_id` - число из адреса проекта
`https://<vikunja>/projects/<id>`.

## Запуск

```
go build -o md2vikunja .
./md2vikunja -dry-run          # только показать, что будет создано
./md2vikunja                   # создать задачи
./md2vikunja -file other.md -config other.yaml
```

## Формат задачи

Задача - секция `## T<N>: <название>` с полями сразу под заголовком.
Всё после строки `vikunja:` до следующего `## ` уходит в описание как есть (markdown).

```markdown
## T4: Lorem ipsum dolor sit amet

priority: P1
status: todo
cases: 7.4a
vikunja:

### Где

Consectetur adipiscing elit, sed do eiusmod tempor.

### Что сделать

- Ut enim ad minim veniam.

### Приёмка

- [ ] quis nostrud exercitation ullamco laboris
```

В Vikunja получится задача `T4: Lorem ipsum dolor sit amet` с приоритетом high
и описанием:

```
cases: 7.4a

### Где
...
```

После создания строка `vikunja:` станет `vikunja: https://vikunja.example.com/tasks/123`.

## Что пропускается

- `status: done`
- `vikunja:` уже заполнен
- нет строки `vikunja:`
- секции `## ...` без префикса `T<N>:` (группы, бэклог, заметки)

## Приоритет

| md | Vikunja |
| --- | --- |
| P1 | 3 (high) |
| P2 | 2 (medium) |
| P3 | 1 (low) |

Остальное - 0 (не задан).
