package mysql

import (
	dbe "github.com/kucjac/go-rest-sdk/errors/dberrors"
)

var mysqlErrMap = map[interface{}]dbe.DBError{
	"08": dbe.ErrConnExc,
	"22": dbe.ErrDataException,
	"25": dbe.ErrInvalidTransState,
	"42": dbe.ErrInvalidSyntax,
	"XA": dbe.ErrInvalidTransState,

	"08S01": dbe.ErrConnExc,
	"08001": dbe.ErrConnExc,
	"21000": dbe.ErrCardinalityViolation,

	"22000": dbe.ErrDataException,
	"22007": dbe.ErrDataException,
	"23000": dbe.ErrIntegrConstViolation,
	"25000": dbe.ErrInvalidTransState,
	"28000": dbe.ErrInvalidAuthorization,

	"42000": dbe.ErrInvalidSyntax,
	"44000": dbe.ErrCheckViolation,
	"HY000": dbe.ErrInternalError,

	uint16(1022): dbe.ErrUniqueViolation,
	uint16(1046): dbe.ErrInvalidCatalogName,
	uint16(1048): dbe.ErrNotNullViolation,
	uint16(1050): dbe.ErrInvalidSyntax,
	uint16(1062): dbe.ErrUniqueViolation,
	uint16(1114): dbe.ErrProgramLimitExceeded,
	uint16(1118): dbe.ErrProgramLimitExceeded,
	uint16(1130): dbe.ErrInvalidAuthorization,
	uint16(1131): dbe.ErrInvalidAuthorization,
	uint16(1132): dbe.ErrInvalidPassword,
	uint16(1133): dbe.ErrInvalidPassword,
	uint16(1169): dbe.ErrUniqueViolation,
	uint16(1182): dbe.ErrTransRollback,
	uint16(1216): dbe.ErrForeignKeyViolation,
	uint16(1217): dbe.ErrForeignKeyViolation,
	uint16(1227): dbe.ErrInsufficientPrivilege,
	uint16(1251): dbe.ErrInvalidAuthorization,
	uint16(1400): dbe.ErrInvalidTransState,
	uint16(1401): dbe.ErrInternalError,
	uint16(1451): dbe.ErrForeignKeyViolation,
	uint16(1452): dbe.ErrForeignKeyViolation,
	uint16(1557): dbe.ErrUniqueViolation,
	uint16(1568): dbe.ErrUniqueViolation,
	uint16(1698): dbe.ErrInvalidPassword,
}

var codeSQLState = map[uint16]string{
	1022: "23000",
	1037: "HY001",
	1038: "HY001",
	1040: "08004",
	1042: "08S01",
	1043: "08S01",
	1044: "42000",
	1045: "28000",
	1046: "3D000",
	1047: "08S01",
	1048: "23000",
	1049: "42000",
	1050: "42S01",
	1051: "42S02",
	1052: "23000",
	1053: "08S01",
	1054: "42S22",
	1055: "42000",
	1056: "42000",
	1057: "42000",
	1058: "21S01",
	1059: "42000",
	1060: "42S21",
	1061: "42000",
	1062: "23000",
	1063: "42000",
	1064: "42000",
	1065: "42000",
	1066: "42000",
	1067: "42000",
	1068: "42000",
	1069: "42000",
	1070: "42000",
	1071: "42000",
	1072: "42000",
	1073: "42000",
	1074: "42000",
	1075: "42000",
	1080: "08S01",
	1081: "08S01",
	1082: "42S12",
	1083: "42000",
	1084: "42000",
	1090: "42000",
	1091: "42000",
	1101: "42000",
	1102: "42000",
	1103: "42000",
	1104: "42000",
	1106: "42000",
	1107: "42000",
	1109: "42S02",
	1110: "42000",
	1112: "42000",
	1113: "42000",
	1115: "42000",
	1118: "42000",
	1120: "42000",
	1121: "42000",
	1129: "HY000",
	1130: "HY000",
	1131: "42000",
	1132: "42000",
	1133: "42000",
	1136: "21S01",
	1138: "42000",
	1139: "42000",
	1140: "42000",
	1141: "42000",
	1142: "42000",
	1143: "42000",
	1144: "42000",
	1145: "42000",
	1146: "42S02",
	1147: "42000",
	1148: "42000",
	1149: "42000",
	1152: "08S01",
	1153: "08S01",
	1154: "08S01",
	1155: "08S01",
	1156: "08S01",
	1157: "08S01",
	1158: "08S01",
	1159: "08S01",
	1160: "08S01",
	1161: "08S01",
	1162: "42000",
	1163: "42000",
	1164: "42000",
	1166: "42000",
	1167: "42000",
	1169: "23000",
	1170: "42000",
	1171: "42000",
	1172: "42000",
	1173: "42000",
	1176: "42000",
	1177: "42000",
	1178: "42000",
	1179: "25000",
	1184: "08S01",
	1189: "08S01",
	1190: "08S01",
	1203: "42000",
	1205: "40001",
	1207: "25000",
	1211: "42000",
	1213: "40001",
	1216: "23000",
	1217: "23000",
	1218: "08S01",
	1222: "21000",
	1226: "42000",
	1227: "42000",
	1230: "42000",
	1231: "42000",
	1232: "42000",
	1234: "42000",
	1235: "42000",
	1239: "42000",
	1241: "21000",
	1242: "21000",
	1247: "42S22",
	1248: "42000",
	1249: "01000",
	1250: "42000",
	1251: "08004",
	1252: "42000",
	1253: "42000",
	1261: "01000",
	1262: "01000",
	1263: "01000",
	1264: "01000",
	1265: "01000",
	1280: "42000",
	1281: "42000",
	1286: "42000",
	1292: "22007",
	1303: "2F003",
	1304: "42000",
	1305: "42000",
	1308: "42000",
	1309: "42000",
	1310: "42000",
	1311: "01000",
	1312: "0A000",
	1313: "42000",
	1314: "0A000",
	1315: "42000",
	1316: "42000",
	1317: "70100",
	1318: "42000",
	1319: "42000",
	1320: "42000",
	1321: "2F005",
	1322: "42000",
	1323: "42000",
	1324: "42000",
	1325: "24000",
	1326: "24000",
	1327: "42000",
	1329: "02000",
	1330: "42000",
	1331: "42000",
	1332: "42000",
	1333: "42000",
	1335: "0A000",
	1336: "0A000",
	1337: "42000",
	1338: "42000",
	1339: "20000",
	1365: "22012",
	1367: "22007",
	1370: "42000",
	1397: "XAE04",
	1398: "XAE05",
	1399: "XAE07",
	1400: "XAE09",
	1401: "XAE03",
	1402: "XA100",
	1403: "42000",
	1406: "22001",
	1407: "42000",
	1410: "42000",
	1413: "42000",
	1414: "42000",
	1415: "0A000",
	1416: "22003",
	1425: "42000",
	1426: "42000",
	1427: "42000",
	1437: "42000",
	1439: "42000",
	1440: "XAE08",
	1441: "22008",
	1451: "23000",
	1452: "23000",
	1453: "42000",
	1458: "42000",
	1460: "42000",
	1461: "42000",
	1463: "42000",
	1557: "23000",
	1568: "25001",
	1582: "42000",
	1583: "42000",
	1584: "42000",
	1586: "23000",
	1613: "XA106",
	1614: "XA102",
	1630: "42000",
	1641: "42000",
	1642: "01000",
	1643: "02000",
	1645: "0K000",
	1687: "42000",
	1690: "22003",
	1698: "28000",
	1701: "42000",
	1758: "35000",
	1761: "23000",
	1762: "23000",
	1792: "25006",
	1845: "0A000",
	1846: "0A000",
	1859: "23000",
	1873: "28000",
	1887: "0Z002",
	1903: "2201E",
}
