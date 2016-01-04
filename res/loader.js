var filename = process.argv[2];
if (!filename) {
	console.log('Usage:   node loader.js <filename.js>');
	return
}

var path = require('path');
var module = path.relative(__dirname, path.resolve(filename));

require('babel/register')({stage:1});
require('./' + module);
