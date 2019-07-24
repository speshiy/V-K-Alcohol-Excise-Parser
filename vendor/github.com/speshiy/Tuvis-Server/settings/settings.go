package settings

//ServerConfigType тип конфигурации сервера
var ServerConfigType = "local"

//IsRelease release flag
var IsRelease = false

//DBHostDefault default host
var DBHostDefault = "127.0.0.1"

//DBHostMainDefault default host
var DBHostMainDefault = "127.0.0.1"

//DBHostClientDefault default host
var DBHostClientDefault = "127.0.0.1"

//DBHostAnalytics default host
var DBHostAnalytics = "127.0.0.1"

//URL of server
var URL string

//URLFrontend of client
var URLFrontend string

//URLFrontendClient of client for client
var URLFrontendClient string

//EmailFrom of service
var EmailFrom string

//EmailFromPassword application
var EmailFromPassword string

//Locale application RU/EN
var Locale string

//PortHTTP application
var PortHTTP string

//PortHTTPS application
var PortHTTPS string

//PortService application
var PortService string

//DBRP application root
var DBRP string

//DBRTUP application root
var DBRTUP string

//DBMP application main
var DBMP string

//DBCP application client
var DBCP string

//DBDP application demo
var DBDP string

//DBAP application analytics
var DBAP string

//ResourcesPath path of pictures
var ResourcesPath string

//FrontendPath path of web frontend
var FrontendPath string

//IsSSL release flag
var IsSSL = false

//CertPath path of web frontend
var CertPath string

//KeyPath path of web frontend
var KeyPath string
