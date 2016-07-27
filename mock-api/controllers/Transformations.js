'use strict';

var url = require('url');


var Transformations = require('./TransformationsService');


module.exports.transformsGET = function transformsGET (req, res, next) {
  Transformations.transformsGET(req.swagger.params, res, next);
};
