# goTSV

# 🔨 Description 

Service for scanning directories with files in .tsv format and rewriting information into pdf files

# 🧩 Config

```json
{
  "host": "localhost:8080",
  "dir_from": "From",
  "dir_to": "To",
  "dsn": "postgres://bestuser:bestuser@localhost:5432/?sslmode=disable",
  "refresh_interval": 5
}
```
# 📞 Request and Response

Request

```http
POST http://localhost:8080/ HTTP/1.1
Content-Type: application/json
{
    "unitguid": "01749246-95f6-57db-b7c3-2ae0e8be671f",
    "page": 2
    "limit": 2
}
```

Response

```json
{
    [
        {
            "Number": "3",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "cold7_ComprSK_status",
            "MessageText": "Компрессор",
            "Context": "",
            "MessageClass": "working",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "cold7_status.ComprSK_status",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        },
        {
            "Number": "4",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "cold7_Temp_Al_HH",
            "MessageText": "Высокая температура",
            "Context": "",
            "MessageClass": "alarm",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "cold7_status.Temp_Al_HH",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        }
    ],
    [
        {
            "Number": "16",
            "MQTT": "",
            "InventoryID": "G-044322",
            "UnitGUID": "01749246-95f6-57db-b7c3-2ae0e8be671f",
            "MessageID": "test_alarm",
            "MessageText": "Тест Аларм",
            "Context": "",
            "MessageClass": "alarm",
            "Level": "100",
            "Area": "LOCAL",
            "Address": "TestingForMsg.Alarm",
            "Block": "",
            "Type": "",
            "Bit": "",
            "InvertBit": ""
        },
    ],
}
