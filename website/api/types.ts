import { FileWithPath } from "react-dropzone";

export enum Filter {
  GRAYSCALE,
  SEPIA,
  INVERT,
}

export enum Progress {
  READY,
  PROCESSING,
  DONE,
  FAILED,
}

export interface ImageDocument {
  id: string; // UUID to identify document
  title: string; // Title of art piece
  author: string; // Author of art piece
  tag: Filter; // Filter that has been applied
  progress: Progress; // Current progress of image processing in SQS
  image: string; // S3 filename of image related to this document
}

export interface UploadResponse {
  image: string; // S3 filename of image
  url: string; // S3 pre-sign URL
}

export interface PreviewFileWithPath {
  preview: string;
  file: FileWithPath;
}

export interface QueueRequestBody {
  title: string;
  author: string;
  image: string;
  tag: Filter;
}

export interface QueueResponse {
  id: string;
}

export interface ProgressResponse {
  id: string;
  progress: Progress;
  imageurl: string;
}
