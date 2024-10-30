# Example of service calls

To call an snet service via bot, follow these steps:

1. In the room with the bot, send a command using the following format:

```
proto_name snet_id service_name method_name
```

* `descriptor_name`: protocol name of the snet service.
* `snet_id`: id for the snet service.
* `service_name`: the specific snet service you want to use.
* `method_name`: the method you want to invoke in the service.

Example:

```
paraphrase-generation paraphrase paraphrase paraphrase
```

2. After sending the command, the bot will reply with a link to proceed with payment. Click on the link and follow the instructions to complete the payment.
3. Once the payment is successfully processed, the bot will notify you in the chat and proceed to the next step.
4. The bot will then prompt you to provide the required input parameters one by one.
5. After receiving all required inputs, the bot will call the snet service on your behalf. The result will be returned directly in the chat.
