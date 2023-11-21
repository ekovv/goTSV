CREATE TABLE tsv (
                     ID           SERIAL PRIMARY KEY,
                     Number       int,
                     MQTT         VARCHAR(255),
                     InventoryID  VARCHAR(255),
                     UnitGUID     VARCHAR(255),
                     MessageID    VARCHAR(255),
                     MessageText  TEXT,
                     Context      VARCHAR(255),
                     MessageClass VARCHAR(255),
                     Level        INTEGER,
                     Area         VARCHAR(255),
                     Address      VARCHAR(255),
                     Block        BOOLEAN,
                     Type         VARCHAR(255),
                     Bit          INTEGER,
                     InvertBit    INTEGER
);
CREATE TABLE checkedFiles (
                              name VARCHAR(255) PRIMARY KEY,
                              error VARCHAR(255)
);