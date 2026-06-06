-- +goose Up

CREATE TABLE IF NOT EXISTS dbmapping (
    id int unsigned NOT NULL AUTO_INCREMENT,
    -- numbers
    anint int, -- in Go: int32
    abigint bigint, -- in Go: int64
    anintunsighed int unsigned, -- in Go: uint32
    abigintunsigned bigint unsigned, -- in Go: uint64
    adecimal decimal(11, 2), -- in Go: float64
    afloat float(11, 2), -- in Go: float32
    adouble double(11, 2), -- in Go: float64
    -- booleans
    aboolean boolean, -- in Go: bool
    -- strings
    avarchar varchar(20), -- Max length : 65535 bytes; in Go: string
    adatetime datetime, -- Format: YYYY-MM-DD hh:mm:ss; in Go: time.Time
    ablob blob(20), -- Max length : 65535 bytes; in Go: []byte
    PRIMARY KEY (id)
);

INSERT INTO dbmapping (anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob) 
VALUES
    (-42, -123456789, 100, 200, 123.45, 123.456789, 123.456789, TRUE, 'Bake a cake', STR_TO_DATE('2025-02-18 15:44:04', '%Y-%m-%d %H:%i:%s'), NULL),
    (0, 0, 0, 0, 0.0, 0.0, 0.0, FALSE, '', 0, NULL),
    (2147483647, 9223372036854775807, 4294967295, 18446744073709551615, 123456789.01, 123456789.01, 123456789.01, FALSE, 'Take out the trash90', STR_TO_DATE('9999-12-31 23:59:59', '%Y-%m-%d %H:%i:%s'), NULL),
    (NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL)
;

-- +goose Down

DROP TABLE IF EXISTS dbmapping;
