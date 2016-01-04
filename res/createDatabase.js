import r from 'rethinkdb';
import dbconfig from './_dbconfig.js';

async function dropDatabase() {
  try {
    let conn = await r.connect(dbconfig);
    await r.dbDrop(dbconfig.db).run(conn);
    console.log('Db dropped', dbconfig.db);
  } catch (e) {
    console.error('Error dropping: ', dbconfig.db, e);
    return;
  }
}

async function createDatabase() {
  try {
    let conn = await r.connect(dbconfig);
    await r.dbCreate(dbconfig.db).run(conn);
    console.log('Db created', dbconfig.db);
  } catch (e) {
    console.error('Error creating: ', dbconfig.db, e);
    return;
  }
}

async function createTable(tableName, data, indexes) {
  try {
    let conn = await r.connect(dbconfig)
    await r.db(dbconfig.db).tableCreate(tableName).run(conn);
    for (var i = 0; i < indexes.length; i++) {
      await r.db(dbconfig.db).table(tableName).indexCreate(indexes[i], {multi: true}).run(conn)
    }

    console.log('Table created', dbconfig.db, tableName);
    for (let item of data) {
      await r.db(dbconfig.db).table(tableName).insert(item).run(conn);
    }
  } catch(e) {
    console.error('Error creating table', dbconfig.db, tableName, e);
    return;
  }
}


async function main() {
  await dropDatabase();
  await createDatabase();
  await createTable('comment', [], ['photo_id', 'account_id', 'notification_id', 'is_known', 'tags']);
  await createTable('photo', [], ['account_id', 'is_private']);
  await createTable('follower', [], ['account_id', 'follower_id']);
  await createTable('account', [], ['username']);

}

main().then(() => {
  console.log('Done!');
  process.exit(0);
}, err => {
  console.log('Error', err.stack)
})
