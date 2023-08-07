
# Postman

[Postman](https://www.postman.com/product/what-is-postman/) allows calling
HTTP endpoints with some flexibility to execute scripts, replace variables,
setup an environment and share these with other developers.

This project provides a [MGC postman_collection](./MGC.postman_collection.json)
to explore some APIs we cover in our CLI and SDK.

## Setup

To properly setup Postman follow the steps bellow:

- Install [Postman](https://www.postman.com/downloads/)
- Import the collection file in Postman: `File -> Import...`
- Click on the MGC collection and go to the `Variables` tab
- Add the ClientID, Client Secret and Client UUID
- Go to the `Authorization` tab
- Click on `Create New Token`
  - Authorize all scopes if requested
- On the `Current token` section, select `mgc_access_token` for the Token property
- Enable the `Auto-refresh token` option

Now using the VPN you should be able to make a request to the MGC servers.

## FAQ

**My request is not using the latest token**

Make sure under the MGC Collection `Authorization` tab there is only one token
on the `Token` property in the `Current Token` section. If there is more than
one click on `Manage Tokens` and remove all of them.
