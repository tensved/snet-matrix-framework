---
description: >-
  Information about requests for payments that allow users to pay for SNET
  service calls.
---

# Payments

## Get payment info

<mark style="color:green;">`GET`</mark> `/payment`

Allows you to get the payment information needed to make a payment on the blockchain.

**Query params**

| Key  | Type   | Description     |
| ---- | ------ | --------------- |
| `id` | string | UUID of payment |

**Response**

{% tabs %}
{% tab title="200" %}
```json
{
  "data": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "url": "ethereum:{tokenAddr}/transfer?address={toAddr}&uint256={amount}",
    "key": "roomId userId serviceName methodName",
    "txHash": "0xabc123...",
    "tokenAddress": "0xdef456...",
    "toAddress": "0xghi789...",
    "amount": 1,
    "status": "pending",
    "createdAt": "2024-07-19T12:34:56Z",
    "updatedAt": "2024-07-19T13:34:56Z",
    "expiresAt": "2024-07-20T12:34:56Z"
  }
}
```
{% endtab %}

{% tab title="400" %}
```
Failed to parse uuid
```
{% endtab %}

{% tab title="400" %}
```
Failed to get payment state
```
{% endtab %}
{% endtabs %}



## Update tx hash of payment

<mark style="color:green;">`PUT`</mark> `/payment`

Allows you to update tx hash of payment.

**Headers**

| Name         | Value              |
| ------------ | ------------------ |
| Content-Type | `application/json` |

**Body**

| Name     | Type   | Description        |
| -------- | ------ | ------------------ |
| `id`     | string | UUID of payment    |
| `txHash` | string | Tx hash of payment |

**Response**

{% tabs %}
{% tab title="200" %}
```
OK
```
{% endtab %}

{% tab title="400" %}
```json
Invalid input
```
{% endtab %}

{% tab title="500" %}
```
Failed to update payment state
```
{% endtab %}
{% endtabs %}
