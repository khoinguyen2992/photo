import path from 'path';
import fs from 'fs';

let configFile = process.env.CONFIG || path.join(__dirname, '../config-default.json');

let config = JSON.parse(fs.readFileSync(configFile).toString());
let dbconfig = {
	host: config.rethinkdb.RETHINKDB_ADDR,
	port: config.rethinkdb.RETHINKDB_PORT,
	db: config.rethinkdb.RETHINKDB_DBNAME
};

export default dbconfig;
