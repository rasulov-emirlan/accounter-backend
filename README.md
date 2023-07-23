# esep-backend

Backend for the ESEP 2.0

## Idea behind this project

This app will provide an easy to use api with i18n. This api will help people with small businesses to keep track of their warehouse, profits and losses.
Our api will be able to analyze what kind of products are selling the most and when.

## Get UP and RUNNING

First you have to make sure you have docker and docker compose installed. After that you can run the following command in the root of the project.

```
docker-compose up
```

or

```
docker compose up
```

Now you must see a cool log with all the applications inside of our docker compose. And now you can go to the `localhost:8080` and start playing with out api. At the `localhost:8080/health/ready` you can check if all the systems we depend on work correctly.

If you wish to run our app without docker you can do so using makefile in the root. Its default command will compile and run our app in development mode, which adds colors to our logs. Also if you want you can run our app with --help flag, and it will show all available flags.

## Plans

[] - Logic for items. Create items which can have different sizes and can be placed in categories.

[] - Logic for sales. Owners must be able to sell their products. And the sales should be recorded and easy to filter through. All the sales should be reflected in the owners number of available products. For example if an owner sells a white t-shirt of size M, the number of t-shirts with exact same specs should be reduced.

[] - Oauth

[] - Errors from validation should be objects instead of a string

[] - i18n

[] - Kuber

[] - Rewrite in Rust
