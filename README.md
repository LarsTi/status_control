# Status Control

This repository is to manage one or any number of status across multiple PCs with a low overhead.
Imagine you are using a fileshare with no lock methods and you want to ensure, you are alone on the file.
You could message all other people, that you are going to start editing and leave a message, when you are done.

You could also call this webapi to check whether someone is already changing it.

## Customizing

The status are offered to the webapi by using environmental variables.
- DEFAULT_{{n}} = The default status, if not set its ""
- SECRET_{{n}} = The secret for updating this status
- STATUS_{{n}} = The id of the status

{{n}} is any number greater than 0.

The Programm starts to read with n = 1 and adds 1 to n as long as it finds no more SECRET_{{n}} or STATUS_{{n}}.

## Querying the API

You can simply call a GET-Request on the URL, replacing STATUS_{{n}} with its respective value defined above
/get?id=STATUS_{{n}}

## Updating the API
You can simply call a GET-Request or UPDATE-Request on the URL, replacing STATUS_{{n}}, SECRET_{{n}} with its respective value defined above.
NEW_STATUS can be choosed by you, its a free field.

/update?id=STATUS_{{n}}&pw=SECRET_{{n}}&status=NEW_STATUS

## DEMO

### Environment

- DEFAULT_1=free
- SECRET_1=SuperS3cret!
- STATUS_1=status

### Querying

/get?id=status

### Updating

/update?id=status&status=test&pw=SuperS3cret --> Will change it to "test"

/update?id=status&status=free&pw=SuperS3cret --> Will change it back to "free"
