# Photo

This is a server of the Golang test app. It makes sure appropriate application logic and provides necessary API for client side development.

# How to run

*  recommend linux or vagrant
*  install go1.4 - [gvm](https://github.com/moovweb/gvm)
*  install node - [nvm](https://github.com/creationix/nvm)
*  install rethinkdb - [rethinkdb.com](http://rethinkdb.com/)
*  git clone git@github.com:khoinguyen2992/photo.git
*  cd photo
* export GOPATH=$PWD
* go install photo/cmd/photo-server

#References
* Architecture - [grox](https://github.com/golang-vietnam/grox)
*  Gorethink - [gorethink](https://github.com/dancannon/gorethink) - for communicating with rethinkdb
*  Imaging - [imaging](https://github.com/disintegration/imaging) - for croping images
*  Gomail - [gomail](https://github.com/go-gomail/gomail/) - for sending email
*  Govalidator - [govalidator](https://github.com/asaskevich/govalidator) - for validation

# API Docs

Coming soon
