CREATE TABLE Occurrence (
                     ID           SERIAL PRIMARY KEY,
                     Number       VARCHAR(255),
                     MQTT         VARCHAR(255),
                     InventoryID  VARCHAR(255),
                     UnitGUID     VARCHAR(255),
                     MessageID    VARCHAR(255),
                     MessageText  TEXT,
                     Context      VARCHAR(255),
                     MessageClass VARCHAR(255),
                     Level        VARCHAR(255),
                     Area         VARCHAR(255),
                     Address      VARCHAR(255),
                     Block        VARCHAR(255),
                     Type         VARCHAR(255),
                     Bit          VARCHAR(255),
                     InvertBit    VARCHAR(255)
);
CREATE TABLE checkedFiles (
                              name VARCHAR(255) PRIMARY KEY,
                              error VARCHAR(255)
);