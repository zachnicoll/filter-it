import * as feed from "./feedEndpoint";
import * as upload from "./uploadEndpoint";
import * as queue from "./queueEndpoint";
import * as s3 from "./s3Endpoint";

const API = {
  feed,
  upload,
  queue,
  s3,
};

export default API;
