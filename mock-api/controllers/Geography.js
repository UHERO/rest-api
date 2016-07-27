'use strict';

var url = require('url');


var Geography = require('./GeographyService');


module.exports.geoGET = function geoGET (req, res, next) {
  Geography.geoGET(req.swagger.params, res, next);
};
