docker pull microsoft/mssql-server-linux
docker run -d --name ms-sql-server -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=sek1gunzC8' -p 1433:1433 microsoft/mssql-server-linux